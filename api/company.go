package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type createCompanyRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) createCompany(ctx *gin.Context) {
	var req createCompanyRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	company, err := server.store.CreateCompany(ctx, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, company)
}
