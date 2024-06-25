package main

import (
	"log"

	"github.com/AriJaya07/go-rest-api/packages/config"
	"github.com/AriJaya07/go-rest-api/packages/config/db"
	"github.com/AriJaya07/go-rest-api/packages/models/store"
	api "github.com/AriJaya07/go-rest-api/packages/routes"
	"github.com/go-sql-driver/mysql"
)

func main() {
	cfg := mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	sqlStorage := db.NewMySQLStorage(cfg)

	db, err := sqlStorage.Init()
	if err != nil {
		log.Fatal(err)
	}

	store := store.NewStore(db)

	server := api.NewAPIServer(":3000", store)
	server.Serve()
}
