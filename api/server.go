package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/mateoradman/tempus/internal/db/sqlc"
)

// Server stores information about a server.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("gender", validGender)
	}

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
	router.GET("/users", server.listUser)

	server.router = router
	return server
}

// Start runs a server on a given address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
