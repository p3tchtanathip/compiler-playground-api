package main

import (
	"compiler-playground-api/internal/delivery/websocket"
	"compiler-playground-api/internal/repository"
	"compiler-playground-api/internal/usecase"
	"fmt"
	"log"
	"net/http"
)

func main() {
	tempFileRepo := repository.NewTempFileRepository()
	codeUseCase := usecase.NewCodeUseCase(tempFileRepo)
	wsHandler := websocket.NewWebSocketHandler(codeUseCase)

	// http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/ws", wsHandler.HandleWebsocket)

	fmt.Println("Server starting at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
