package update

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tsinghua-cel/strategy-gen/types"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func GetCommand() *cobra.Command {
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update given strategy files slice to attacker service with the mode.",
		Run:   runCommand,
	}
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	setFlags(updateCmd)
	return updateCmd
}

func setFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(
		&params.interval,
		intervalFlag,
		1200,
		"the interval to update the strategy",
	)

	cmd.Flags().StringVar(
		&params.attacker,
		attackerFlag,
		"127.0.0.1:12001",
		"the attacker service to update",
	)

	cmd.Flags().StringSliceVar(
		&params.fileSlice,
		sliceFlag,
		[]string{},
		"the strategy files slice to update",
	)

	cmd.Flags().IntVar(
		&params.updateMode,
		modeFlag,
		0,
		"the update mode, 0: sorted loop 1: random loop",
	)

	cmd.Flags().IntVar(
		&params.loopCount,
		loopCountFlag,
		0,
		"the loop count, 0 is never stop",
	)
}

func runCommand(cmd *cobra.Command, _ []string) {
	if len(params.fileSlice) == 0 {
		cmd.Help()
		return
	}
	allStrategy := make(map[string]types.Strategy)
	log.Infof("update strategy files: %v", params.fileSlice)
	for _, file := range params.fileSlice {
		d, err := os.ReadFile(file)
		if err != nil {
			log.WithField("file", file).Errorf("ignore the file because of failed to read file: %s", err)
			continue
		}
		strategy := types.Strategy{}
		err = json.Unmarshal(d, &strategy)
		if err != nil {
			log.WithField("file", file).Errorf("ignore the file because of failed to unmarshal data: %s", err)
			continue
		}
		allStrategy[file] = strategy
	}
	stopCheck := func(count int) bool {
		if params.loopCount == 0 {
			return false
		}
		return count >= params.loopCount
	}
	total := 0
	for stopCheck(total) != true {
		var idx = 0
		if params.updateMode == 0 {
			// sorted loop
			idx = total % len(params.fileSlice)
		} else {
			idx = rand.Intn(len(params.fileSlice))
		}
		err := updateStrategy(params.attacker, allStrategy[params.fileSlice[idx]])
		if err != nil {
			log.WithField("file", params.fileSlice[idx]).Errorf("failed to update strategy: %s", err)
			time.Sleep(time.Second)
		} else {
			log.WithField("file", params.fileSlice[idx]).Info("update strategy success")
			time.Sleep(time.Duration(params.interval) * time.Second)
		}

		total++
	}
}

func updateStrategy(url string, strategy types.Strategy) error {
	d, err := json.Marshal(strategy)
	if err != nil {
		return err
	}

	res, err := http.Post(fmt.Sprintf("http://%s/v1/update-strategy", url), "application/json", bytes.NewReader(d))
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("failed to generate strategy: %s", res.Status)
	}
	return nil
}
