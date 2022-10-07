package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	db "github.com/mateoradman/tempus/internal/db/sqlc"
	"github.com/mateoradman/tempus/util"
)

type UserRequest struct {
	Username  string     `json:"username" binding:"required,alphanum"`
	Password  string     `json:"password" binding:"required,min=6"`
	Name      string     `json:"name" binding:"required,min=1"`
	Surname   string     `json:"surname" binding:"required,min=1"`
	Email     string     `json:"email" binding:"required,email"`
	CompanyID *int64     `json:"company_id"`
	Gender    string     `json:"gender" binding:"required,gender"`
	BirthDate *time.Time `json:"birth_date" binding:"lt"`
	Language  string     `json:"language" binding:"required,len=2,ascii"`
	Country   string     `json:"country" binding:"required,len=2,ascii"`
	Timezone  *string    `json:"timezone"`
	ManagerID *int64     `json:"manager_id"`
	TeamID    *int64     `json:"team_id"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req UserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:  req.Username,
		Password:  hashedPassword,
		Email:     req.Email,
		Name:      req.Name,
		Surname:   req.Surname,
		CompanyID: req.CompanyID,
		Gender:    req.Gender,
		BirthDate: *req.BirthDate,
		Language:  req.Language,
		Country:   req.Country,
		Timezone:  req.Timezone,
		ManagerID: req.ManagerID,
		TeamID:    req.TeamID,
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (server *Server) getUser(ctx *gin.Context) {
	var req RequestWithID
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (server *Server) deleteUser(ctx *gin.Context) {
	var req RequestWithID
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.DeleteUser(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (server *Server) updateUser(ctx *gin.Context) {
	var reqID RequestWithID
	var req UserRequest
	if err := ctx.ShouldBindUri(&reqID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams{
		ID:        reqID.ID,
		Username:  req.Username,
		Email:     req.Email,
		Name:      req.Name,
		Surname:   req.Surname,
		CompanyID: req.CompanyID,
		Gender:    req.Gender,
		BirthDate: *req.BirthDate,
		Language:  req.Language,
		Country:   req.Country,
		Timezone:  req.Timezone,
		ManagerID: req.ManagerID,
		TeamID:    req.TeamID,
	}
	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type listUserRequest struct {
	Offset int32 `form:"offset" binding:"min=0"`
	Limit  int32 `form:"limit,default=10" binding:"min=1,max=100"`
}

func (server *Server) listUser(ctx *gin.Context) {
	var req listUserRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	}
	users, err := server.store.ListUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}
