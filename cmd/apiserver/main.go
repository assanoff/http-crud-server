package main

import (
	"fmt"
	"log"

	"github.com/assanoff/http-crud-server/internal/app/apiserver"
)

func main() {

	config := apiserver.NewConfig()
	fmt.Println(config.DatabaseURL)

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
