package main

import (
	"time"

	"github.com/thavlik/bjjvb/base/pkg/base"

	"github.com/spf13/cobra"
	compiler "github.com/thavlik/bjjvb/compiler/pkg/api"
	"github.com/thavlik/bjjvb/filter/pkg/labelstore"
	mongo_labelstore "github.com/thavlik/bjjvb/filter/pkg/labelstore/mongo"
	postgres_labelstore "github.com/thavlik/bjjvb/filter/pkg/labelstore/postgres"
	"github.com/thavlik/bjjvb/filter/pkg/server"
)

var serverArgs struct {
	base.ServerOptions
	db         base.DatabaseOptions
	compiler   base.ServiceOptions
	slideShow  base.ServiceOptions
	stackSize  int
	collection string
}

var serverCmd = &cobra.Command{
	Use: "server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.ServerEnv(&serverArgs.ServerOptions)
		base.DatabaseEnv(&serverArgs.db, true)
		base.ServiceEnv("compiler", &serverArgs.compiler)
		base.ServiceEnv("slide-show", &serverArgs.slideShow)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		go base.RunMetrics(serverArgs.MetricsPort, base.Log)
		base.RandomizeSeed()
		return server.Entry(
			serverArgs.Port,
			initLabelStore(),
			compiler.NewCompilerClientFromOptions(serverArgs.compiler),
			serverArgs.slideShow,
			serverArgs.stackSize,
			base.Log,
		)
	},
}

func initLabelStore() labelstore.LabelStore {
	switch serverArgs.db.Driver {
	case base.MongoDriver:
		return mongo_labelstore.NewMongoLabelStore(
			base.ConnectMongo(&serverArgs.db.Mongo))
	case base.PostgresDriver:
		return postgres_labelstore.NewPostgresLabelStore(
			base.ConnectPostgres(&serverArgs.db.Postgres))
	default:
		panic("unreachable branch detected")
	}
}

func init() {
	base.AddDatabaseFlags(serverCmd, &serverArgs.db)
	base.AddServerFlags(serverCmd, &serverArgs.ServerOptions)
	base.AddServiceFlags(serverCmd, "compiler", &serverArgs.compiler, 8*time.Second)
	base.AddServiceFlags(serverCmd, "slide-show", &serverArgs.slideShow, 20*time.Second)
	serverCmd.PersistentFlags().IntVar(&serverArgs.stackSize, "stack-size", 5, "number of markers to serve in a stack")
	serverCmd.PersistentFlags().StringVar(&serverArgs.collection, "collection", "filter", "collection name for the labels")
	ConfigureCommand(serverCmd)
}
