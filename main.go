package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kirinson321/bsg-recruitment/pkg/downloader"
	"github.com/kirinson321/bsg-recruitment/pkg/exchange"
	"github.com/kirinson321/bsg-recruitment/pkg/output"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	downloader := downloader.NewDownloader()
	outputter := output.NewOutputter()
	exchangeService := exchange.NewService(downloader, outputter)

	err := prepLogFile()
	if err != nil {
		panic(err)
	}

	err = exchangeService.GetRates(ctx)
	if err != nil {
		panic(err)
	}

	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-quitSignal
		fmt.Printf("terminating due to signal %v\n", sig)
		cancel()
		os.Exit(1)
	}()

	go exchangeService.GetRates(context.Background())
}

func prepLogFile() error {
	// prepare the log file
	// check if a file with the same name exists
	if logFileExists() {
		err := os.Rename(output.LogFileName, output.BackupLogFileName)
		if err != nil {
			panic(fmt.Errorf("error renaming the existing log file: %w", err))
		}
	}

	// create a new log file
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
