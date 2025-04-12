package databases

import (
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var RD *redis.Client

func ConnectRD() {
	addr := os.Getenv("RD_HOST")
	pass := os.Getenv("RD_PASS")
	dbid := os.Getenv("RD_DB")
	dbidx, err := strconv.Atoi(dbid)
	if err != nil {
		dbidx = 0
	}

	RD = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       dbidx,
	})
}
