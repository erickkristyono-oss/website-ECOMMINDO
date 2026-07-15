package main

import (
	"log"
	"net/http"

	"ecommindo/internal/config"
	httphandler "ecommindo/internal/handlers/http"
	dbrepo "ecommindo/internal/repository/db"
	"ecommindo/internal/usecase"
)

func main() {
	db, err := config.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	jwtSecret := config.JWTSecret()

	userRepo := dbrepo.NewUserRepository(db)
	serviceRepo := dbrepo.NewServiceRepository(db)
	cartRepo := dbrepo.NewCartRepository(db)
	orderRepo := dbrepo.NewOrderRepository(db)

	authUC := usecase.NewAuthUsecase(userRepo, jwtSecret)
	serviceUC := usecase.NewServiceUsecase(serviceRepo)
	cartUC := usecase.NewCartUsecase(cartRepo, serviceRepo)
	orderUC := usecase.NewOrderUsecase(orderRepo, cartUC, cartRepo)

	router := httphandler.NewRouter(authUC, serviceUC, cartUC, orderUC, jwtSecret, config.StaticDir())

	addr := ":" + config.Port()
	log.Printf("Ecommindo Jaya Persada server running at http://localhost%s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
