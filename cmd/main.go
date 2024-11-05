package main

import (
	"log"

	"app.flower.clip/src/routing"
)

func main() {
	err := routing.StartServer()
	if err != nil {
		log.Fatal(err)
	}
}
