package main

import (
	"log"

	"sensQID/internal/app"
)

func main() {
	log.Print("sensitive QID anonymizer started")

	app.Start()
}
