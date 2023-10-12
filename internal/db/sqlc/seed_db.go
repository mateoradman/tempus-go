package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/mateoradman/tempus/internal/config"
	"github.com/mateoradman/tempus/internal/util"
)

func seedSuperUser(ctx context.Context, q *Queries, config config.Config) error {
	_, err := q.GetUserByEmail(ctx, config.SuperUserEmail)
	if err == nil {
		// this means that the user already exists
		return nil
	} else if err != pgx.ErrNoRows {
		// if error is not that the user doesn't exist, return it
		return err
	} 
	// else create a new user
	
	hashedPassword, err := util.HashPassword(config.SuperUserPassword)
	if err != nil {
		return err
	}

	arg := CreateUserParams{
		Username:  config.SuperUserUsername,
		Email:     config.SuperUserEmail,
		Name:      "Superuser",
		Surname:   "Admin",
		Password:  hashedPassword,
		BirthDate: time.Now(),
	}

	_, err = q.CreateUser(ctx, arg)
	return err
}
