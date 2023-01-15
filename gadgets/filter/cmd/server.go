package main

import (
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/gadget"

	"github.com/spf13/cobra"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"
	mongo_labelstore "github.com/thavlik/t4vd/filter/pkg/labelstore/mongo"
	postgres_labelstore "github.com/thavlik/t4vd/filter/pkg/labelstore/postgres"
	"github.com/thavlik/t4vd/gadgets/filter/pkg/server"
)

var serverArgs struct {
	base.ServerOptions
	db                     base.DatabaseOptions
	gadgetID               string
	maxBatchSize           int
	collection             string
	initInputGadget        string
	initInputGadgetChannel string
}

var serverCmd = &cobra.Command{
	Use: "server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.ServerEnv(&serverArgs.ServerOptions)
		base.DatabaseEnv(&serverArgs.db, true)
		base.CheckEnv("COLLECTION", &serverArgs.collection)
		base.CheckEnv("GADGET_ID", &serverArgs.gadgetID)
		base.CheckEnvInt("MAX_BATCH_SIZE", &serverArgs.maxBatchSize)
		base.CheckEnv("INIT_INPUT_GADGET", &serverArgs.initInputGadget)
		base.CheckEnv("INIT_INPUT_GADGET_CHANNEL", &serverArgs.initInputGadgetChannel)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log := base.DefaultLog
		go base.RunMetrics(serverArgs.MetricsPort, log)
		base.RandomizeSeed()
		return server.Entry(
			serverArgs.ServerOptions.Port,
			initLabelStore(),
			serverArgs.gadgetID,
			gadget.NewDataRef(
				serverArgs.initInputGadget,
				serverArgs.initInputGadgetChannel,
			),
			serverArgs.maxBatchSize,
			log,
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
		panic(base.Unreachable)
	}
}

func init() {
	base.AddDatabaseFlags(serverCmd, &serverArgs.db)
	base.AddServerFlags(serverCmd, &serverArgs.ServerOptions)
	serverCmd.PersistentFlags().StringVar(&serverArgs.collection, "collection", "filter", "collection name for the labels")
	serverCmd.PersistentFlags().StringVar(&serverArgs.gadgetID, "gadget-id", "", "owner gadget uuid")
	serverCmd.PersistentFlags().IntVar(&serverArgs.maxBatchSize, "max-batch-size", 100, "maximum batch size for the labels")
	serverCmd.PersistentFlags().StringVar(&serverArgs.initInputGadget, "init-input-gadget", "", "initial input gadget name")
	serverCmd.PersistentFlags().StringVar(&serverArgs.initInputGadgetChannel, "init-input-gadget-channel", "", "initial input gadget channel name")
	ConfigureCommand(serverCmd)
}
