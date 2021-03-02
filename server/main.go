package main

import (
	"log"

	"github.com/arganaphangquestian/usermanagement/server/repository"
	"github.com/arganaphangquestian/usermanagement/server/route"
	"github.com/arganaphangquestian/usermanagement/server/service"
)

func main() {
	repository := repository.New()
	service := service.New(repository)
	app := route.New(service)
	log.Fatalln(app.Listen(":8080"))
}
