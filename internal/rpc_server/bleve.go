package rpc_server

import (
	"github.com/bartmika/bleve-server/pkg/dtos"
)

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
