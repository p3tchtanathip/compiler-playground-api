package websocket

import (
	"compiler-playground-api/internal/domain"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WebsocketHandler struct {
	upgrader    websocket.Upgrader
	codeUseCase domain.CodeUseCase
}

func NewWebSocketHandler(codeUseCase domain.CodeUseCase) *WebsocketHandler {
	return &WebsocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		codeUseCase: codeUseCase,
	}
}

func (h *WebsocketHandler) HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		var msg domain.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("JSON unmarshal error: %v", err)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		h.codeUseCase.ExecuteCode(ctx, msg, conn)
		conn.Close()
	}
}
