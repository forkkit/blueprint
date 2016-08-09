package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"github.com/wercker/auth/middleware"
	"github.com/wercker/{{package .Name}}/core"
	"golang.org/x/net/context"

	"gopkg.in/urfave/cli.v1"
)

var gatewayCommand = cli.Command{
	Name:   "gateway",
	Usage:  "Starts environment variable HTTP->gRPC gateway",
	Action: gatewayAction,
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:   "port, p",
			Value:  {{ .Gateway }},
			EnvVar: "HTTP_PORT",
		},
		cli.StringFlag{
			Name:   "host",
			Value:  "localhost:{{ .Port }}",
			EnvVar: "GRPC_HOST",
		},
	},
}

var gatewayAction = func(c *cli.Context) error {
	log.Println("Starting env var gateway")

	o, err := ParseGatewayOptions(c)
	if err != nil {
		log.Println(err)
		return cli.NewExitError("Invalid arguments", 3)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = core.Register{{class .Name}}HandlerFromEndpoint(ctx, mux, o.Host, opts)
	if err != nil {
		log.Println(err)
		return err
	}

	authMiddleware := middleware.AuthTokenMiddleware(mux)

	log.Printf("Listening on port %v\n", o.Port)
	http.ListenAndServe(fmt.Sprintf(":%v", o.Port), authMiddleware)

	return nil
}

func ParseGatewayOptions(c *cli.Context) (*GatewayOptions, error) {
	port := c.Int("port")
	if !validPortNumber(port) {
		return nil, ErrInvalidPortNumber
	}

	host := c.String("host")

	return &GatewayOptions{
		Port: port,
		Host: host,
	}, nil
}

type GatewayOptions struct {
	Port int
	Host string
}
