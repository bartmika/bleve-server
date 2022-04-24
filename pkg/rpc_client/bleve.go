package rpc_client

import (
	// "errors"
	"log"

	"github.com/bartmika/bleve-server/pkg/dtos"
)

func (s *BleveService) Index(filename string, identifier string, data []byte) error {
	req := &dtos.IndexRequestDTO{
		Filename: filename,
		Identifier: identifier,
		Data:       data,
	}
	var reply dtos.IndexResponseDTO
	err := s.call("RPC.Index", req, &reply)
	if err != nil {
		log.Println("RPC CLIENT RPC CLIENT ERROR | BleveService | Index | err", err)
		return err
	}
	return nil
}

func (s *BleveService) Query(filename string, search string) (*[]string, error) {
	req := &dtos.QueryRequestDTO{
		Filename: filename,
		Search:     search,
	}
	var reply dtos.QueryResponseDTO
	err := s.call("RPC.Query", req, &reply)
	if err != nil {
		log.Println("RPC CLIENT RPC CLIENT ERROR | BleveService | Query | err", err)
		return nil, err
	}
	return &reply.UUIDs, nil
}
