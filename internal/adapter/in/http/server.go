package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	Router *gin.Engine
	url    string
	cfg    AuthConfig
	Handlers
}

type Handlers struct {
	User     *UserHandler
	Gift     *GiftHandler
	Wishlist *WishlistHandler
}

func NewServer(url string, cfg AuthConfig, handlers Handlers) *Server {
	router := gin.Default()
	server := &Server{Router: router, url: url, cfg: cfg, Handlers: handlers}

	server.RegisterRoutes()
	return server
}

func (s *Server) Run() error {
	return s.Router.Run(s.url)
}

func (s *Server) RegisterRoutes() {
	authMiddleWare := AuthMiddleware(s.cfg)

	s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := s.Router.Group("/api")

	//groups
	auth := api.Group("/auth")
	wishlist := api.Group("/wishlist")
	gift := wishlist.Group("/:wishlistId/gifts")

	public := api.Group("/public")

	s.User.RegisterRoutes(auth)
	s.Wishlist.RegisterRoutes(wishlist, authMiddleWare)
	s.Gift.RegisterRoutes(gift, authMiddleWare)

	//public endpoints
	public.POST("wishlist/:token/gifts/:id", s.Gift.Book)
	public.GET("wishlist/token/:token", s.Wishlist.GetByToken)
}
