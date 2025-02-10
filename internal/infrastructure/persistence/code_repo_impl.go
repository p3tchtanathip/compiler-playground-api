package persistence

import "log"

type CodeRepository struct{}

func NewCodeRepository() *CodeRepository {
	return &CodeRepository{}
}

func (r *CodeRepository) SaveCodeMetadata(id, language, createdAt string) error {
	log.Printf("Saving Code Metadata: ID=%s, Language=%s, CreatedAt=%s\n", id, language, createdAt)
	return nil
}
