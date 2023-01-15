package main

import (
	"os"

	"github.com/thavlik/t4vd/base/pkg/base"
	"go.uber.org/zap"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		base.DefaultLog.Error("main", zap.String("err", err.Error()))
		os.Exit(1)
	}
}
