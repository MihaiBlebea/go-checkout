package main

import (
	"log"

	"github.com/MihaiBlebea/go-checkout/cmd"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	cmd.Execute()
}
