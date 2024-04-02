package main

import (
	"Hygieia/api"
	"Hygieia/database"
	"context"
	"database/sql"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Errorf("cannot connect to redis")
		return
	}
	db, err := sql.Open("mysql", "root:314159@/hygieia?parseTime=true")
	if err != nil {
		fmt.Errorf("cannot connect to mysql")
		return
	}

	store := database.NewMySqlStore(db)

	server, err := api.NewServer(rdb, store)
	if err != nil {
		fmt.Errorf("cannot create a new server")
		return
	}

	server.Start("0.0.0.0:8080")

}
