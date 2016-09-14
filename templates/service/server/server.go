package server

import (
	"github.com/wercker/blueprint/templates/service/core"
	"github.com/wercker/blueprint/templates/service/state"

	"golang.org/x/net/context"
)

// New instantiates the server
func New(store state.Store) (*BlueprintServer, error) {
	return &BlueprintServer{
		store: store,
	}, nil
}

// BlueprintServer is a grpc server
type BlueprintServer struct {
	store state.Store
}

// Action performs the task at hand
func (s *BlueprintServer) Action(ctx context.Context, req *core.ActionRequest) (*core.ActionResponse, error) {
	return &core.ActionResponse{}, nil
}

// Make sure that {{.Name}}Server implements the core.{{.Name}}Service interface.
var _ core.BlueprintServer = &BlueprintServer{}
