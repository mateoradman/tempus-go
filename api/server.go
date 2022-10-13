package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/mateoradman/tempus/internal/db/sqlc"
	"github.com/mateoradman/tempus/internal/token"
	"github.com/mateoradman/tempus/util"
)

// Server stores information about a server.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("gender", validGender)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/companies", server.createCompany)
	router.GET("/companies/:id", server.getCompany)
	router.DELETE("/companies/:id", server.deleteCompany)
	router.PUT("/companies/:id", server.updateCompany)
	router.GET("/companies", server.listCompany)
	router.GET("/companies/:id/employees", server.listCompanyEmployees)

	router.POST("/users", server.createUser)
	router.GET("/users/:id", server.getUser)
	router.DELETE("/users/:id", server.deleteUser)
	router.PUT("/users/:id", server.updateUser)
	router.GET("/users", server.listUsers)
	router.POST("/users/login", server.loginUser)

	server.router = router
}

// Start runs a server on a given address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
