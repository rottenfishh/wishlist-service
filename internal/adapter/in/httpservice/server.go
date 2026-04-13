package httpservice

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	httpServer *http.Server
	cfg        AuthConfig
}

type AuthConfig struct {
	JWTSecret string
}

type Handlers struct {
	User     *UserHandler
	Gift     *GiftHandler
	Wishlist *WishlistHandler
}

func NewServer(url string, cfg AuthConfig, handlers Handlers) *Server {
	router := gin.Default()

	httpServer := &http.Server{
		Addr:    url,
		Handler: router,
	}

	server := &Server{httpServer: httpServer, cfg: cfg}

	server.RegisterRoutes(router, handlers)

	return server
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) RegisterRoutes(router *gin.Engine, h Handlers) {
	authMiddleWare := AuthMiddleware(s.cfg)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")

	//groups
	auth := api.Group("/auth")
	wishlist := api.Group("/wishlist")
	gift := wishlist.Group("/:wishlistId/gifts")

	public := api.Group("/public")

	h.User.RegisterRoutes(auth)
	h.Wishlist.RegisterRoutes(wishlist, authMiddleWare)
	h.Gift.RegisterRoutes(gift, authMiddleWare)

	//public endpoints
	public.POST("/wishlist/:token/gifts/:id", h.Gift.Book)
	public.GET("/wishlist/token/:token", h.Wishlist.GetByToken)
}
