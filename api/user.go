package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	db "github.com/mateoradman/tempus/internal/db/sqlc"
	"github.com/mateoradman/tempus/util"
)

type UserRequest struct {
	Username  string     `json:"username" binding:"required,alphanum"`
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

type createUserRequest struct {
	Password string `json:"password" binding:"required,min=6"`
	UserRequest
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
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

type updateUserRequest struct {
	Name      *string    `json:"name,omitempty" binding:"omitempty,min=1"`
	Surname   *string    `json:"surname,omitempty" binding:"omitempty,min=1"`
	Gender    *string    `json:"gender,omitempty" binding:"omitempty,gender"`
	BirthDate *time.Time `json:"birth_date,omitempty" binding:"omitempty,lt"`
	Language  *string    `json:"language,omitempty" binding:"omitempty,len=2,ascii"`
	Country   *string    `json:"country,omitempty" binding:"omitempty,len=2,ascii"`
}

func (server *Server) updateUser(ctx *gin.Context) {
	var reqID RequestWithID
	var req updateUserRequest
	if err := ctx.ShouldBindUri(&reqID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, reqID.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams{
		ID:        user.ID,
		Name:      req.Name,
		Surname:   req.Surname,
		Gender:    req.Gender,
		BirthDate: req.BirthDate,
		Language:  req.Language,
		Country:   req.Country,
	}
	user, err = server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type listUserRequest struct {
	Offset int32 `form:"offset" binding:"min=0"`
	Limit  int32 `form:"limit,default=10" binding:"min=1,max=100"`
}

func (server *Server) listUsers(ctx *gin.Context) {
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

type loginUserRequest struct {
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required,alphanum"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  db.User   `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err = util.CheckPassword(req.Password, user.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  user,
	}

	ctx.JSON(http.StatusOK, resp)
}
