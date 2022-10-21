package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mateoradman/tempus/config"
	db "github.com/mateoradman/tempus/internal/db/sqlc"
	"github.com/mateoradman/tempus/util"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := config.Config{
		TokenSymmetricKey:    util.RandomString(32),
		AccessTokenDuration:  10 * time.Minute,
		RefreshTokenDuration: 24 * time.Hour,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
