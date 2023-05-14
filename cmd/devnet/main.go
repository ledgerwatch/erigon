package main

import (
	"fmt"
	"sync"

	"github.com/ledgerwatch/erigon/cmd/devnet/commands"
	"github.com/ledgerwatch/erigon/cmd/devnet/devnetutils"
	"github.com/ledgerwatch/erigon/cmd/devnet/models"
	"github.com/ledgerwatch/erigon/cmd/devnet/node"
	"github.com/ledgerwatch/erigon/cmd/devnet/services"
	"os"
	"time"

	"github.com/ledgerwatch/erigon/cmd/utils"
	"github.com/ledgerwatch/erigon/params"
	erigoncli "github.com/ledgerwatch/erigon/turbo/cli"
	"github.com/ledgerwatch/erigon/turbo/debug"
	"github.com/ledgerwatch/erigon/turbo/logging"
	"github.com/urfave/cli/v2"
)

func main() {

	debug.RaiseFdLimit()

	app := cli.NewApp()
	app.Version = params.VersionWithCommit(params.GitCommit)
	app.Action = func(ctx *cli.Context) error {
		return action(ctx)
	}
	app.Flags = append(erigoncli.DefaultFlags, debug.Flags...) // debug flags are required
	app.Flags = append(app.Flags, utils.MetricFlags...)
	app.Flags = append(app.Flags, logging.Flags...)

	app.After = func(ctx *cli.Context) error {
		// unsubscribe from all the subscriptions made
		services.UnsubscribeAll()
		return nil
	}
	if err := app.Run([]string{}); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func action(ctx *cli.Context) error {
	dataDir := ctx.String("datadir")
	if dataDir == "" {
		return fmt.Errorf("missing --datadir flag - required for devnet tool")
	}
	fmt.Printf("datadir = %s\n", dataDir)
	logger := logging.SetupLoggerCtx("devnet", ctx, true /* rootLogger */)
	// clear all the dev files
	if err := devnetutils.ClearDevDB(logger); err != nil {
		return err
	}
	// wait group variable to prevent main function from terminating until routines are finished
	var wg sync.WaitGroup

	// remove the old logs from previous runs
	if err := devnetutils.DeleteLogs(logger); err != nil {
		return err
	}

	// start the first erigon node in a go routine
	node.Start(&wg, logger)

	// send a quit signal to the quit channels when done making checks
	node.QuitOnSignal(&wg)

	// sleep for seconds to allow the nodes fully start up
	time.Sleep(time.Second * 10)

	// start up the subscription services for the different sub methods
	services.InitSubscriptions([]models.SubMethod{models.ETHNewHeads})

	// execute all rpc methods amongst the two nodes
	commands.ExecuteAllMethods()

	// wait for all goroutines to complete before exiting
	wg.Wait()
	return nil
}
