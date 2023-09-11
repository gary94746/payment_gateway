package main

import (
	"log"

	"payment-processor.gary94746/main/server/rest"
)

func main() {
	server := rest.ApiRest{}
	err := server.Serve()
	if err != nil {
		log.Fatal("Not able to serve", err.Error())
	}
}
