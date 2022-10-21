package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	db "github.com/mateoradman/tempus/internal/db/sqlc"
)

type EntryRequest struct {
	UserID    int64      `json:"user_id" binding:"required,min=1"`
	StartTime time.Time  `json:"start_time" binding:"required,timestamp"`
	EndTime   *time.Time `json:"end_time" binding:"omitempty,timestamp"`
	Date      time.Time  `json:"date" binding:"required,timestamp"`
}

func (server *Server) createEntry(ctx *gin.Context) {
	var req EntryRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateEntryParams{
		UserID:    req.UserID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Date:      req.Date,
	}
	entry, err := server.store.CreateEntry(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, entry)
}

func (server *Server) getEntry(ctx *gin.Context) {
	var req RequestWithID
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	entry, err := server.store.GetEntry(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entry)
}

func (server *Server) deleteEntry(ctx *gin.Context) {
	var req RequestWithID
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	entry, err := server.store.DeleteEntry(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entry)
}

func (server *Server) updateEntry(ctx *gin.Context) {
	var reqID RequestWithID
	var req EntryRequest
	if err := ctx.ShouldBindUri(&reqID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateEntryParams{
		ID:   reqID.ID,
		UserID: req.UserID,
		StartTime: req.StartTime,
		EndTime: req.EndTime,
	}
	entry, err := server.store.UpdateEntry(ctx, arg)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entry)
}

func (server *Server) listEntries(ctx *gin.Context) {
	var req PaginationRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListEntriesParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	}
	companies, err := server.store.ListEntries(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, companies)
}
