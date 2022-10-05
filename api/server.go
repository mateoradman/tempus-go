package api

import (
	"github.com/gin-gonic/gin"
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

	router.POST("/companies", server.createCompany)
	router.GET("/companies/:id", server.getCompany)
	router.GET("/companies", server.listCompany)

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
