package rpc_client

import (
	// "errors"
	"encoding/json"
	"log"

	"github.com/bartmika/bleve-server/pkg/dtos"
)

func (s *BleveService) Register(filenames []string) error {
	req := &dtos.RegisterRequestDTO{
		Filenames: filenames,
	}
	var reply dtos.RegisterResponseDTO
	err := s.call("RPC.Register", req, &reply)
	if err != nil {
		log.Println("rpc_client | RPC.Registe | err", err)
		return err
	}
	return nil
}

func (s *BleveService) Index(filename string, identifier string, data any) error {
	dataBin, err := json.Marshal(data)
	if err != nil {
		log.Println("rpc_client | Marshal | err", err)
		return err
	}

	req := &dtos.IndexRequestDTO{
		Filename:   filename,
		Identifier: identifier,
		Data:       dataBin,
	}
	var reply dtos.IndexResponseDTO
	err = s.call("RPC.Index", req, &reply)
	if err != nil {
		log.Println("rpc_client | RPC.Index | err", err)
		return err
	}
	return nil
}

func (s *BleveService) Query(filename string, search string) ([]string, error) {
	req := &dtos.QueryRequestDTO{
		Filename: filename,
		Search:   search,
	}
	var reply dtos.QueryResponseDTO
	err := s.call("RPC.Query", req, &reply)
	if err != nil {
		log.Println("rpc_client | RPC.Query | err", err)
		return nil, err
	}
	return reply.UUIDs, nil
}

func (s *BleveService) Delete(filename string, identifier string) error {
	req := &dtos.DeleteRequestDTO{
		Filename:   filename,
		Identifier: identifier,
	}
	var reply dtos.DeleteResponseDTO
	err := s.call("RPC.Index", req, &reply)
	if err != nil {
		log.Println("rpc_client | RPC.Index | err", err)
		return err
	}
	return nil
}
