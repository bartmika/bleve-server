package dtos // Data Transfer Object

type QueryRequestDTO struct {
	Filename string `json:"filename"`
	Search   string `json:"search"`
}

type QueryResponseDTO struct {
	UUIDs []string
}
