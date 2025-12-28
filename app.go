package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"multiRip/config"
	"multiRip/ripper"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx       context.Context
	appConfig *config.Config
}

func NewApp(cfg *config.Config) *App {
	return &App{
		appConfig: cfg,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	r, w, _ := os.Pipe()
	processStdout := os.Stdout

	os.Stdout = w
	os.Stderr = w

	go func() {
		tee := io.TeeReader(r, processStdout)
		scanner := bufio.NewScanner(tee)

		for scanner.Scan() {
			line := scanner.Text()
			runtime.EventsEmit(a.ctx, "global-log", line+"\n")
		}
	}()
}

func (a *App) RunRipper(jobsFile string) error {
	jobsConfig, err := config.LoadJobs(jobsFile)
	if err != nil {
		return fmt.Errorf("Failed to load jobs config: %v", err)
	}

	go func() {
		fmt.Println("Starting Jobs")

		onLog := func(driveId string, message string) {
			runtime.EventsEmit(a.ctx, "log-update", map[string]string{
				"driveId": driveId,
				"message": message,
			})
		}

		jobs, err := ripper.GetJobs(a.appConfig, jobsConfig)
		if err != nil {
			fmt.Printf("Orchestration error: %v\n", err)
			runtime.EventsEmit(a.ctx, "job-error", err.Error())
			return
		}
		if err = ripper.ExecuteJobs(jobs, onLog); err != nil {
			fmt.Printf("Orchestration error: %v\n", err)
			runtime.EventsEmit(a.ctx, "job-error", err.Error())
			return
		}

		fmt.Println("All drive operations finished successfully.")
		runtime.EventsEmit(a.ctx, "jobs-complete", true)
	}()
	return nil
}
