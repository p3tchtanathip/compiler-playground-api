package usecase

import (
	"compiler-playground-api/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"

	"github.com/gorilla/websocket"
)

type CodeUseCase struct {
	tempFileRepo domain.TempFileRepository
}

func NewCodeUseCase(tempFileRepo domain.TempFileRepository) *CodeUseCase {
	return &CodeUseCase{
		tempFileRepo: tempFileRepo,
	}
}

func (c *CodeUseCase) ExecuteCode(ctx context.Context, msg domain.Message, conn *websocket.Conn) {
	if msg.Type != "code" {
		log.Printf("Expected code message, got %s", msg.Type)
		return
	}

	fileName, err := c.tempFileRepo.Create(msg.Content)
	if err != nil {
		log.Printf("Error creating temp file: %v", err)
		return
	}
	defer c.tempFileRepo.Delete(fileName)

	cmd := exec.CommandContext(ctx, "python3", fileName)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("Error creating stdin pipe: %v", err)
		return
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error creating stdout pipe: %v", err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Error creating stderr pipe: %v", err)
		return
	}

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				return
			}

			var msg struct {
				Type    string `json:"type"`
				Content string `json:"content"`
			}
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Printf("JSON unmarshal error: %v", err)
				continue
			}

			if msg.Type == "input" {
				if _, err := stdin.Write([]byte(msg.Content + "\n")); err != nil {
					log.Printf("Error writing to stdin: %v", err)
					return
				}
				fmt.Println("input:", msg.Content)
			}
		}
	}()

	if err := cmd.Start(); err != nil {
		log.Printf("Error starting command: %v", err)
		return
	}

	go c.readPipe(stdout, conn, "output")
	go c.readPipe(stderr, conn, "error")

	if err := cmd.Wait(); err != nil {
		log.Printf("Command execution error: %v", err)
		return
	}
}

func (c *CodeUseCase) readPipe(pipe io.Reader, conn *websocket.Conn, msgType string) {
	buffer := make([]byte, 1024)
	for {
		n, err := pipe.Read(buffer)
		if n > 0 {
			content := string(buffer[:n])
			message := domain.Message{Type: msgType, Content: content}
			messageJson, _ := json.Marshal(message)
			if err := conn.WriteMessage(websocket.TextMessage, messageJson); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
		}
		if err != nil {
			break
		}
	}
}
