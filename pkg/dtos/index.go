package dtos // Data Transfer Object

type IndexRequestDTO struct {
	Filename string `json:"instance_id"`
	Identifier string `json:"identifier"`
	Data       []byte `json:"data"`
}

type IndexResponseDTO struct {
	UUIDs []string
}
