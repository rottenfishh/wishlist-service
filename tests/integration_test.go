//go:build integration

package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"wishlist-service/internal/adapter/in/dto"
	"wishlist-service/internal/adapter/in/httpservice"
	"wishlist-service/internal/app"
	"wishlist-service/internal/model"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type integrationEnv struct {
	server *httptest.Server
}

// register -> login -> create wishlist
func TestRegisterLogin_CreateWishlist(t *testing.T) {
	t.Parallel()

	env := newIntegrationEnv(t)
	email := uniqueEmail(t)

	registerResp := postJSON[dto.RegisterRequest, model.User](t, env.server.Client(), env.server.URL+"/api/auth/register", nil, dto.RegisterRequest{
		Email:    email,
		Password: "secret123",
	})

	require.Equal(t, http.StatusCreated, registerResp.StatusCode)
	require.NotEqual(t, "", registerResp.Body.ID.String())

	loginResp := postJSON[dto.LoginRequest, dto.Token](t, env.server.Client(), env.server.URL+"/api/auth/login", nil, dto.LoginRequest{
		Email:    email,
		Password: "secret123",
	})

	require.Equal(t, http.StatusOK, loginResp.StatusCode)
	require.NotEmpty(t, loginResp.Body.Token)

	createWishlistResp := postJSON[dto.CreateWishlistRequest, dto.WishlistResponse](t, env.server.Client(), env.server.URL+"/api/wishlists", map[string]string{
		"Authorization": "Bearer " + loginResp.Body.Token,
	}, dto.CreateWishlistRequest{
		Title:       "Birthday",
		Description: "Wishlist for birthday",
		Date:        mustParseTime(t, "2030-01-02T15:04:05Z"),
	})

	require.Equal(t, http.StatusCreated, createWishlistResp.StatusCode)
	require.Equal(t, "Birthday", createWishlistResp.Body.Title)
	require.NotEqual(t, "", createWishlistResp.Body.Token.String())
}

// get wishlist by public token -> book gift -> second book returns conflict
func TestPublicBook_GiftConflict(t *testing.T) {
	t.Parallel()

	env := newIntegrationEnv(t)

	email := uniqueEmail(t)
	postJSON[dto.RegisterRequest, model.User](t, env.server.Client(), env.server.URL+"/api/auth/register", nil, dto.RegisterRequest{
		Email:    email,
		Password: "secret123",
	})

	loginResp := postJSON[dto.LoginRequest, dto.Token](t, env.server.Client(), env.server.URL+"/api/auth/login", nil, dto.LoginRequest{
		Email:    email,
		Password: "secret123",
	})
	require.Equal(t, http.StatusOK, loginResp.StatusCode)

	authHeaders := map[string]string{"Authorization": "Bearer " + loginResp.Body.Token}

	wishlistResp := postJSON[dto.CreateWishlistRequest, dto.WishlistResponse](t, env.server.Client(), env.server.URL+"/api/wishlists", authHeaders, dto.CreateWishlistRequest{
		Title:       "New Year",
		Description: "Public wishlist",
		Date:        mustParseTime(t, "2031-01-01T00:00:00Z"),
	})
	require.Equal(t, http.StatusCreated, wishlistResp.StatusCode)

	wishlistID := wishlistResp.Body.ID
	publicToken := wishlistResp.Body.Token.String()

	giftResp := postJSON[dto.CreateGiftRequest, dto.GiftResponse](t, env.server.Client(),
		fmt.Sprintf("%s/api/wishlists/%d/gifts", env.server.URL, wishlistID),
		authHeaders,
		dto.CreateGiftRequest{
			Name:        "LEGO",
			Description: "Big set",
			Link:        "https://example.com/gift",
			Priority:    5,
		},
	)
	require.Equal(t, http.StatusCreated, giftResp.StatusCode)

	giftID := giftResp.Body.ID

	bookURL := fmt.Sprintf("%s/api/public/wishlists/%s/gifts/%d", env.server.URL, publicToken, giftID)

	firstBookResp := postJSON[struct{}, dto.GiftResponse](t, env.server.Client(), bookURL, nil, struct{}{})
	require.Equal(t, http.StatusOK, firstBookResp.StatusCode)
	require.True(t, firstBookResp.Body.Booked)

	secondBookResp := postJSON[struct{}, dto.ErrorResponse](t, env.server.Client(), bookURL, nil, struct{}{})
	require.Equal(t, http.StatusConflict, secondBookResp.StatusCode)
	require.Equal(t, "already_booked", secondBookResp.Body.Error)
}

type jsonResponse[T any] struct {
	StatusCode int
	Body       T
}

func postJSON[Req any, Resp any](t *testing.T, client *http.Client, url string, headers map[string]string, payload Req) jsonResponse[Resp] {
	t.Helper()

	rawPayload, err := json.Marshal(payload)
	require.NoError(t, err)

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(rawPayload))
	require.NoError(t, err)

	request.Header.Set("Content-Type", "application/json")
	applyHeaders(request, headers)

	return executeJSON[Resp](t, client, request)
}

func applyHeaders(request *http.Request, headers map[string]string) {
	for key, value := range headers {
		request.Header.Set(key, value)
	}
}

func executeJSON[Resp any](t *testing.T, client *http.Client, request *http.Request) jsonResponse[Resp] {
	t.Helper()

	response, err := client.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()

	rawResponse, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	var parsed Resp
	if len(rawResponse) > 0 {
		err = json.Unmarshal(rawResponse, &parsed)
		require.NoError(t, err, "response body: %s", string(rawResponse))
	}

	return jsonResponse[Resp]{
		StatusCode: response.StatusCode,
		Body:       parsed,
	}
}

func newIntegrationEnv(t *testing.T) *integrationEnv {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	t.Cleanup(cancel)

	req := testcontainers.ContainerRequest{
		Image:        "postgres:17",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "wishlist",
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = container.Terminate(context.Background())
	})

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)

	server := httptest.NewServer(buildHandler(t, host, port.Port()))
	t.Cleanup(server.Close)

	return &integrationEnv{
		server: server,
	}
}

func buildHandler(t *testing.T, host, port string) http.Handler {
	t.Helper()

	cfg := &app.Config{
		DatabaseConfig: app.DatabaseConfig{
			DatabaseName:  "wishlist",
			Username:      "postgres",
			Password:      "postgres",
			Host:          host,
			Port:          port,
			MigrationsDir: "../database",
		},
		AuthConfig: httpservice.AuthConfig{
			JWTSecret:  "integration-secret",
			JwtExpires: 3600,
		},
		ServerPort: "0",
	}

	application, err := app.NewApp(context.Background(), cfg)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = application.Shutdown()
	})

	return application.Handler()
}

func uniqueEmail(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("user-%d@example.com", time.Now().UnixNano())
}

func mustParseTime(t *testing.T, value string) time.Time {
	t.Helper()

	parsed, err := time.Parse(time.RFC3339, value)
	require.NoError(t, err)
	return parsed
}
