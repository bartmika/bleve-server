package dtos // Data Transfer Object

type RegisterRequestDTO struct {
	Filenames []string `json:"filenames"`
}

type RegisterResponseDTO struct{}
