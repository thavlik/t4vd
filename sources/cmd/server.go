package main

import (
	"time"

	"github.com/thavlik/t4vd/base/cmd/iam"
	"github.com/thavlik/t4vd/base/pkg/base"

	"github.com/spf13/cobra"
	compiler "github.com/thavlik/t4vd/compiler/pkg/api"
	seer "github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/sources/pkg/server"
	"github.com/thavlik/t4vd/sources/pkg/store"
	mongo_store "github.com/thavlik/t4vd/sources/pkg/store/mongo"
	postgres_store "github.com/thavlik/t4vd/sources/pkg/store/postgres"
)

var serverArgs struct {
	base.ServerOptions
	iam      base.IAMOptions
	seer     base.ServiceOptions
	compiler base.ServiceOptions
	gateway  base.ServiceOptions
	db       base.DatabaseOptions
}

var serverCmd = &cobra.Command{
	Use: "server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.IAMEnv(&serverArgs.iam, false)
		base.ServerEnv(&serverArgs.ServerOptions)
		base.ServiceEnv("compiler", &serverArgs.compiler)
		base.ServiceEnv("seer", &serverArgs.seer)
		base.DatabaseEnv(&serverArgs.db, true)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log := base.Log
		go base.RunMetrics(serverArgs.MetricsPort, log)
		return server.Entry(
			serverArgs.Port,
			iam.InitIAM(&serverArgs.iam),
			initStore(&serverArgs.db),
			seer.NewSeerClientFromOptions(serverArgs.seer),
			compiler.NewCompilerClientFromOptions(serverArgs.compiler),
			log,
		)
	},
}

func initStore(opts *base.DatabaseOptions) store.Store {
	switch opts.Driver {
	case base.MongoDriver:
		return mongo_store.NewMongoStore(
			base.ConnectMongo(&opts.Mongo))
	case base.PostgresDriver:
		return postgres_store.NewPostgresStore(
			base.ConnectPostgres(&opts.Postgres))
	default:
		panic(base.Unreachable)
	}
}

func init() {
	base.AddIAMFlags(serverCmd, &serverArgs.iam)
	base.AddServerFlags(serverCmd, &serverArgs.ServerOptions)
	base.AddServiceFlags(serverCmd, "seer", &serverArgs.seer, 15*time.Second)
	base.AddServiceFlags(serverCmd, "compiler", &serverArgs.compiler, 8*time.Second)
	base.AddDatabaseFlags(serverCmd, &serverArgs.db)
	ConfigureCommand(serverCmd)
}
