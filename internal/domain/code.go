package domain

import (
	"context"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type CodeUseCase interface {
	ExecuteCode(ctx context.Context, msg Message, conn *websocket.Conn)
}
