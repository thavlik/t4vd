package main

import (
	"os"

	"github.com/thavlik/bjjvb/base/pkg/base"
	"go.uber.org/zap"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		base.Log.Error("main", zap.String("err", err.Error()))
		os.Exit(1)
	}
}
