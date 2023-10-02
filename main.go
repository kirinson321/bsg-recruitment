package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kirinson321/bsg-recruitment/pkg/config"
	"github.com/kirinson321/bsg-recruitment/pkg/downloader"
	"github.com/kirinson321/bsg-recruitment/pkg/exchange"
	"github.com/kirinson321/bsg-recruitment/pkg/output"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interval := flag.Uint("interval", 5, "sets the cadence in seconds in which checks should happen")
	numberOfChecks := flag.Uint("numberOfChecks", 10, "sets the number of checks that should run every interval time")
	flag.Parse()
	config := config.Config{
		RateCheckerInterval: (time.Duration(*interval) * time.Second),
		NumberOfChecks:      *numberOfChecks,
	}

	// Prepare services.
	c := http.DefaultClient

	downloader := downloader.NewDownloader(c)
	outputter := output.NewOutputter()
	exchangeService := exchange.NewService(downloader, outputter, config)

	// Prepare the log file.
	err := prepLogFile()
	if err != nil {
		panic(err)
	}

	// Listen for signals for graceful shutdown.
	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-quitSignal
		fmt.Printf("terminating due to signal %v\n", sig)
		cancel()
		os.Exit(1)
	}()

	// Start the exchange rates service.
	exchangeService.GetRates(ctx)
}

// prepLogFile verifies if a file with the same name as the log file exists, and if so, renames it to log.txt.old.
// Then it creates a new log file with name log.txt.
func prepLogFile() error {
	// prepare the log file
	// if a file with the same name exists, rename it to log.txt.old
	if logFileExists() {
		err := os.Rename(output.LogFileName, output.BackupLogFileName)
		if err != nil {
			panic(fmt.Errorf("error renaming the existing log file: %w", err))
		}
	}

	// create a new log file with name log.txt
	f, err := os.Create(output.LogFileName)
	if err != nil {
		panic(fmt.Errorf("error creating the log file: %w", err))
	}
	f.Close()

	return nil
}

func logFileExists() bool {
	_, err := os.Stat("log.txt")
	return !os.IsNotExist(err)
}
