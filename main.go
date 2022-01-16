package main

import (
	"PasswordService/api/handler"
	"PasswordService/app"
	"PasswordService/infrastructure/repository"
)

func main() {

	//
	// Create repos
	inMemDb := repository.InMemDB{}
	passwordRepo := repository.PasswordRepository{Db: &inMemDb}
	passwordRepo.Init()

	//
	// Create Services
	passwordService := app.NewPasswordService(&passwordRepo)
	handlerService := handler.NewHandlers(passwordService)

	//
	// Create http listner
	server := handler.NewServer(":8080", handlerService)
	//
	// Run
	server.ConfigureAndRun()
}
