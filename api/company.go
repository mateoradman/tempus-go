package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	db "github.com/mateoradman/tempus/internal/db/sqlc"
)

type createCompanyRequest struct {
	Name string `json:"name" binding:"required,min=1"`
}

func (server *Server) createCompany(ctx *gin.Context) {
	var req createCompanyRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	company, err := server.store.CreateCompany(ctx, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, company)
}

func (server *Server) getCompany(ctx *gin.Context) {
	var req RequestWithID
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	company, err := server.store.GetCompany(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, company)
}

func (server *Server) deleteCompany(ctx *gin.Context) {
	var req RequestWithID
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	company, err := server.store.DeleteCompany(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, company)
}

func (server *Server) updateCompany(ctx *gin.Context) {
	var reqID RequestWithID
	var reqName createCompanyRequest
	if err := ctx.ShouldBindUri(&reqID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&reqName); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateCompanyParams{
		ID:   reqID.ID,
		Name: reqName.Name,
	}
	company, err := server.store.UpdateCompany(ctx, arg)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, company)
}

func (server *Server) listCompany(ctx *gin.Context) {
	var req PaginationRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListCompaniesParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	}
	companies, err := server.store.ListCompanies(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, companies)
}

func (server *Server) listCompanyEmployees(ctx *gin.Context) {
	var idReq RequestWithID
	var queryReq PaginationRequest
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if _, err := server.store.GetCompany(ctx, idReq.ID); err != nil {
		// check if a company with a given ID exists.
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindQuery(&queryReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListCompanyEmployeesParams{
		ID:     idReq.ID,
		Offset: queryReq.Offset,
		Limit:  queryReq.Limit,
	}

	employees, err := server.store.ListCompanyEmployees(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, employees)
}
