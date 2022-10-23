package db

import (
	"context"
	"time"

	"github.com/mateoradman/tempus/internal/config"
	"github.com/mateoradman/tempus/internal/util"
)

func seedSuperUser(ctx context.Context, q *Queries, config config.Config) error {
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
