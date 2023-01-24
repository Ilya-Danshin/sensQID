package app

import (
	"context"
	"log"

	"sensQID/internal/pkg/cfg"
	"sensQID/internal/pkg/database"
)

func Start() {
	config, err := cfg.Read()
	if err != nil {
		log.Panic(err.Error())
	}

	ctx := context.Background()
	db := database.NewDB()
	err = db.Init(ctx, config.DB)
	if err != nil {
		log.Panic(err.Error())
	}

	info := NewAnonInfo()
	err = info.getInfo(db)
	if err != nil {
		log.Panic(err.Error())
	}

}
