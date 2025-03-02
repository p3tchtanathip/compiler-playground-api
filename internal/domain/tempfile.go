package domain

type TempFileRepository interface {
	Create(content string) (string, error)
	Delete(name string) error
}
