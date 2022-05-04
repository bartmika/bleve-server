package rpc_server

import (
	"github.com/bartmika/bleve-server/pkg/dtos"
)

func (rpc *RPC) Register(req *dtos.RegisterRequestDTO, res *dtos.RegisterResponseDTO) error {
	err := rpc.Controller.Register(req.Filenames)
	if err != nil {
		return err
	}

	*res = dtos.RegisterResponseDTO{}
	return nil
}

func (rpc *RPC) Index(req *dtos.IndexRequestDTO, res *dtos.IndexResponseDTO) error {
	err := rpc.Controller.Index(req.Filename, req.Identifier, string(req.Data))
	if err != nil {
		return err
	}

	*res = dtos.IndexResponseDTO{}
	return nil
}

func (rpc *RPC) Query(req *dtos.QueryRequestDTO, res *dtos.QueryResponseDTO) error {
	uuids, err := rpc.Controller.Query(req.Filename, req.Search)
	if err != nil {
		return err
	}

	*res = dtos.QueryResponseDTO{
		UUIDs: uuids,
	}
	return nil
}

func (rpc *RPC) Delete(req *dtos.DeleteRequestDTO, res *dtos.DeleteResponseDTO) error {
	err := rpc.Controller.Delete(req.Filename, req.Identifier)
	if err != nil {
		return err
	}

	*res = dtos.DeleteResponseDTO{}
	return nil
}
