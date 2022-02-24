//go:build wireinject
// +build wireinject

package main

import (
	"github.com/CPCTF2022/ssh-separator/api"
	"github.com/CPCTF2022/ssh-separator/repository"
	"github.com/CPCTF2022/ssh-separator/repository/badger"
	"github.com/CPCTF2022/ssh-separator/service"
	"github.com/CPCTF2022/ssh-separator/ssh"
	"github.com/CPCTF2022/ssh-separator/store"
	"github.com/CPCTF2022/ssh-separator/store/gomap"
	"github.com/CPCTF2022/ssh-separator/workspace"
	"github.com/CPCTF2022/ssh-separator/workspace/docker"
	"github.com/google/wire"
)

var (
	transactionBind         = wire.Bind(new(repository.ITransaction), new(*badger.Transaction))
	storeWorkspaceBind      = wire.Bind(new(store.IWorkspace), new(*gomap.Workspace))
	repositoryUserBind      = wire.Bind(new(repository.IUser), new(*badger.User))
	workspaceBind           = wire.Bind(new(workspace.IWorkspace), new(*docker.Workspace))
	workspaceConnectionBind = wire.Bind(new(workspace.IWorkspaceConnection), new(*docker.WorkspaceConnection))
	serviceUserBind         = wire.Bind(new(service.IUser), new(*service.User))
	servicePipeBind         = wire.Bind(new(service.IPipe), new(*service.Pipe))
)

type Server struct {
	*service.Setup
	*api.API
	*ssh.SSH
}

func NewServer(setup *service.Setup, a *api.API, s *ssh.SSH) (*Server, error) {
	return &Server{
		Setup: setup,
		API:   a,
		SSH:   s,
	}, nil
}

func InjectServer() (*Server, func(), error) {
	wire.Build(
		NewServer,
		api.NewAPI,
		api.NewUser,
		gomap.NewWorkspace,
		badger.NewDB,
		badger.NewTransaction,
		badger.NewUser,
		service.NewSetup,
		service.NewUser,
		service.NewPipe,
		ssh.NewSSH,
		docker.NewWorkspace,
		docker.NewWorkspaceConnection,
		transactionBind,
		storeWorkspaceBind,
		repositoryUserBind,
		workspaceBind,
		workspaceConnectionBind,
		serviceUserBind,
		servicePipeBind,
	)

	return nil, nil, nil
}
