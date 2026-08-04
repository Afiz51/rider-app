package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/env"
)

var (
	httpAddr = env.GetString("TRIPSERVICE_ADDR", ":8083")
)

func main() {
	fmt.Println(httpAddr)
	ctx := context.Background()
	inmemRepo := repository.NewInmemRepository()

	fare := &domain.RideFareModel{
		UserID: "42",
	}
	svc := service.NewService(inmemRepo)
	t, err := svc.CreateTrip(ctx, fare)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(t)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /preview", func(w http.ResponseWriter, r *http.Request) {
		var body any
		json.NewDecoder(r.Body).Decode(&body)

		fmt.Println(body)
		trip, err := svc.CreateTrip(ctx, fare)

		if err != nil {
			http.Error(w, "Error creating trip", http.StatusNotFound)
		}

		json.NewEncoder(w).Encode(trip)
	})

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP server error: %v", err)
	}
}
