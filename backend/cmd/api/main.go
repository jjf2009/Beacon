package main

import (
	"log"

	"your-module-name/internal/server"
)

func main() {
	srv := server.New()

	log.Println("Beacon API running on :8080")

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}