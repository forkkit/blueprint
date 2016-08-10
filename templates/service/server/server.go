package server

import (
	"github.com/wercker/blueprint/templates/service/core"
	"github.com/wercker/blueprint/templates/service/state"

	"golang.org/x/net/context"
)

func New(store state.Store) (*BlueprintServer, error) {
	return &BlueprintServer{
		store: store,
	}, nil
}

type BlueprintServer struct {
	store state.Store
}

func (s *BlueprintServer) Action(ctx context.Context, req *core.ActionRequest) (*core.ActionResponse, error) {
	return &core.ActionResponse{}, nil
}

// Make sure that {{.Name}}Server implements the core.{{.Name}}Service interface.
var _ core.BlueprintServer = &BlueprintServer{}
