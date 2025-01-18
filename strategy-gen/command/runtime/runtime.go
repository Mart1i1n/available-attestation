package runtime

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tsinghua-cel/strategy-gen/globalinfo"
	"github.com/tsinghua-cel/strategy-gen/library"
	"github.com/tsinghua-cel/strategy-gen/types"
	"github.com/tsinghua-cel/strategy-gen/utils"
	"os"
	"strings"
	"time"
)

func GetCommand() *cobra.Command {
	runtimeCmd := &cobra.Command{
		Use:   "runtime",
		Short: "Run a library strategy generator.",
		Run:   runCommand,
	}

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	setFlags(runtimeCmd)

	return runtimeCmd
}

func setlog(path string) func() {
	if path == "" {
		return nil
	}
	// logrus log to file
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	return func() {
		file.Close()
	}
}

func setFlags(cmd *cobra.Command) {

	cmd.Flags().StringVar(
		&params.logPath,
		logFlag,
		"",
		"the log path",
	)

	cmd.Flags().IntVar(
		&params.maxValidatorIndex,
		maxValidatorIndexFlag,
		-1,
		"the max hack validator index",
	)

	cmd.Flags().IntVar(
		&params.minValidatorIndex,
		minValidatorIndexFlag,
		0,
		"the min hack validator index",
	)

	cmd.Flags().StringVar(
		&params.attacker,
		attackerFlag,
		"127.0.0.1:12001",
		"the attacker service to update",
	)

	cmd.Flags().StringVar(
		&params.strategy,
		strategyFlag,
		"",
		"runtime strategy to used",
	)

	cmd.Flags().BoolVar(
		&params.listLibrary,
		listFlag,
		false,
		"list all library strategies",
	)
	cmd.Flags().IntVar(
		&params.duration,
		durationFlag,
		60,
		"duration to run for each strategy when strategy is all",
	)
}

func runCommand(cmd *cobra.Command, _ []string) {
	closeFunc := setlog(params.logPath)
	defer func() {
		if closeFunc != nil {
			closeFunc()
		}
	}()

	library.Init()

	if params.listLibrary {
		listLibrary()
		return
	}

	if params.strategy == "" {
		log.Fatal("strategy is required")
	}
	// get second per slot
	initFinished := false
	for !initFinished {
		base, err := utils.GetChainBaseInfo(params.attacker)
		if err != nil {
			log.WithError(err).Error("failed to get second per slot")
			time.Sleep(1 * time.Second)
			continue
		} else {
			log.Infof("get chain base info %s", base)
			initFinished = true
			globalinfo.Init(base)
		}
	}
	log.Info("init finished")
	filters := make(map[string]bool)
	if splits := strings.Split(params.strategy, ","); len(splits) > 1 {
		for _, split := range splits {
			filters[split] = true
		}
	}

	{
		strategy, ok := library.GetStrategy(params.strategy)
		if !ok {
			log.Fatalf("strategy %s not found", params.strategy)
		}

		strategy.Run(context.Background(), types.LibraryParams{
			Attacker:          params.attacker,
			MaxValidatorIndex: params.maxValidatorIndex,
			MinValidatorIndex: params.minValidatorIndex,
		})
	}
}

func listLibrary() {
	strategies := library.GetAllStrategies()
	for name, strategy := range strategies {
		fmt.Printf("%s:\n %s\n", name, strategy.Description())
	}
}
