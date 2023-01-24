package app

import (
	"log"

	"sensQID/internal/pkg/cfg"
)

func Start() {
	config, err := cfg.Read()
	if err != nil {
		log.Panic(err.Error())
	}

	log.Println(config)
}
