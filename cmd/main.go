package main

import (
	"genesis"
	"genesis/pkg/handler"
	"genesis/pkg/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @title GSES2 BTC application
// @description Сервіс з API, який дозволить: \n- дізнатись поточний курс біткоіну (BTC) у гривні (UAH) \n- підписати емейл на отримання інформації по зміні курсу \n- запит, який відправить всім підписаним користувачам актуальний курс.
// @version 1.0

// @host localhost:8000
// @BasePath /api
func main() {
	services := service.NewService()
	handlers := handler.NewHandler(services)
	srv := new(genesis.Server)

	go func() {
		if err := srv.Run(os.Getenv("PORT"), handlers.InitRoutes()); err != nil {
			log.Fatalf("Error occured while running http server: %s", err.Error())
		}
	}()

	log.Printf("App started up...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(); err != nil {
		log.Fatalf("Error occured on server shutting down: %s", err.Error())
	}
	log.Printf("App shutting down...")
}
