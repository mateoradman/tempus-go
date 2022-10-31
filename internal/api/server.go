package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/mateoradman/tempus/internal/config"
	db "github.com/mateoradman/tempus/internal/db/sqlc"
	"github.com/mateoradman/tempus/internal/rbac"
	"github.com/mateoradman/tempus/internal/token"
)

// Server stores information about a server.
type Server struct {
	config      config.Config
	store       db.Store
	tokenMaker  token.Maker
	rbacService *rbac.RBACService
	router      *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(config config.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	// role based access control service
	rbacService := rbac.NewRBACService(store)

	server := &Server{
		config:      config,
		store:       store,
		tokenMaker:  tokenMaker,
		rbacService: rbacService,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("gender", validGender)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// routes not protected by auth middleware
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/refresh", server.refreshToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/companies", server.createCompany)
	authRoutes.GET("/companies/:id", server.getCompany)
	authRoutes.DELETE("/companies/:id", server.deleteCompany)
	authRoutes.PUT("/companies/:id", server.updateCompany)
	authRoutes.GET("/companies", server.listCompany)
	authRoutes.GET("/companies/:id/employees", server.listCompanyEmployees)

	authRoutes.GET("/users/:id", server.getUser)
	authRoutes.DELETE("/users/:id", server.deleteUser)
	authRoutes.PATCH("/users/:id", server.updateUser)
	authRoutes.GET("/users", server.listUsers)

	server.router = router
}

// Start runs a server on a given address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
