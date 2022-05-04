package dtos // Data Transfer Object

type DeleteRequestDTO struct {
	Filename   string `json:"instance_id"`
	Identifier string `json:"identifier"`
}

type DeleteResponseDTO struct{}
