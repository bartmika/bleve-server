package dtos // Data Transfer Object

type QueryRequestDTO struct {
	Filename string `json:"instance_id"`
	Search     string `json:"search"`
}

type QueryResponseDTO struct {
	UUIDs []string
}
