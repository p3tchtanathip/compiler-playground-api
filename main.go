package main

import (
	"compiler-playground-api/internal/config"
	"compiler-playground-api/internal/infrastructure/persistence"
	"compiler-playground-api/internal/infrastructure/storage"
	"compiler-playground-api/internal/interface/controller"
	"compiler-playground-api/internal/usecase"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	cfg := config.LoadConfig()

	minioService, err := storage.NewMinioService(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.BucketName,
	)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}
	log.Println("MinIO setup complete")

	repo := persistence.NewCodeRepository()
	useCase := usecase.NewCodeUseCase(repo, minioService)
	handler := controller.NewCodeHandler(useCase)

	mux := http.NewServeMux()
	mux.HandleFunc("/submit_code", handler.SaveCode)
	mux.HandleFunc("/execute_code", handler.ExecuteCode)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
	}).Handler(mux)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
