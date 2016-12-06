package main

import (
	"fmt"
	"net"

	"gopkg.in/mgo.v2"
	"gopkg.in/urfave/cli.v1"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/wercker/blueprint/templates/service/core"
	"github.com/wercker/blueprint/templates/service/server"
	"github.com/wercker/blueprint/templates/service/state"
	"google.golang.org/grpc"
)

var serverCommand = cli.Command{
	Name:   "server",
	Usage:  "start gRPC server",
	Action: serverAction,
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:   "port",
			Value:  666,
			EnvVar: "PORT",
		},
		cli.StringFlag{
			Name:   "mongo",
			Value:  "mongodb://localhost:27017",
			EnvVar: "MONGODB_URI",
		},
		cli.StringFlag{
			Name:  "mongo-database",
			Value: "blueprint",
		},
		cli.StringFlag{
			Name:  "state-store",
			Usage: "storage driver, currently supported [mongo]",
			Value: "mongo",
		},
		cli.StringFlag{
			Name:   "service-key",
			Usage:  "Hex encoded service key to use",
			EnvVar: "WERCKER_SERVICE_KEY",
		},
	},
}

var serverAction = func(c *cli.Context) error {
	log.Info("Starting blueprint server")

	log.Debug("Parsing server options")
	o, err := parseServerOptions(c)
	if err != nil {
		log.WithError(err).Error("Unable to validate arguments")
		return errorExitCode
	}

	store, err := getStore(o)
	if err != nil {
		log.WithError(err).Error("Unable to create state store")
		return errorExitCode
	}
	defer store.Close()

	log.Debug("Creating server")
	server, err := server.New(store)
	if err != nil {
		log.WithError(err).Error("Unable to create server")
		return errorExitCode
	}

	s := grpc.NewServer()
	core.RegisterBlueprintServer(s, server)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", o.Port))
	if err != nil {
		log.WithField("port", o.Port).WithError(err).Error("Failed to listen")
		return errorExitCode
	}

	log.WithField("port", o.Port).Info("Starting server")
	err = s.Serve(lis)
	if err != nil {
		log.WithError(err).Error("Failed to serve gRPC")
		return errorExitCode
	}

	return nil
}

type serverOptions struct {
	MongoDatabase string
	MongoURI      string
	Port          int
	ServiceKey    []byte
	StateStore    string
}

func parseServerOptions(c *cli.Context) (*serverOptions, error) {
	port := c.Int("port")
	if !validPortNumber(port) {
		return nil, fmt.Errorf("Invalid port number: %d", port)
	}

	serviceKey, err := parseServiceKey(c)
	if err != nil {
		return nil, errors.Wrap(err, "invalid wercker service key")
	}

	return &serverOptions{
		MongoDatabase: c.String("mongo-database"),
		MongoURI:      c.String("mongo"),
		Port:          port,
		ServiceKey:    serviceKey,
		StateStore:    c.String("state-store"),
	}, nil
}

func getStore(o *serverOptions) (state.Store, error) {
	switch o.StateStore {
	case "mongo":
		return getMongoStore(o)
	default:
		return nil, fmt.Errorf("Invalid store: %s", o.StateStore)
	}
}

func getMongoStore(o *serverOptions) (*state.MongoStore, error) {
	log.Info("Creating MongoDB store")

	log.WithField("MongoURI", o.MongoURI).Debug("Dialing the MongoDB cluster")
	session, err := mgo.Dial(o.MongoURI)
	if err != nil {
		return nil, errors.Wrap(err, "Dialing the MongoDB cluster failed")
	}

	log.WithField("MongoDatabase", o.MongoDatabase).Debug("Creating MongoDB store")
	store, err := state.NewMongoStore(session, o.MongoDatabase)
	if err != nil {
		return nil, errors.Wrap(err, "Creating MongoDB store failed")
	}

	return store, nil
}
