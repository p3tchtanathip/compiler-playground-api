package entity

import "time"

type Code struct {
	ID         string    `json:"id"`
	Language   string    `json:"language"`
	SourceCode string    `json:"source_code"`
	Input      string    `json:"input"`
	CreatedAt  time.Time `json:"created_at"`
}
