package main

import (
	"embed"
	"fmt"
	"log"
	"multiRip/config"
	"multiRip/ripper"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	startGui := func() {
		app := NewApp()

		err := wails.Run(&options.App{
			Title:  "MultiRip",
			Width:  1024,
			Height: 768,
			AssetServer: &assetserver.Options{
				Assets: assets,
			},
			// BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
			OnStartup: app.startup,
			Bind: []any{
				app,
			},
		})

		if err != nil {
			println("Error:", err.Error())
		}
	}

	args := os.Args[1:]
	if len(args) == 0 {
		startGui()
		return
	}

	command := args[0]
	switch command {
	case "help":
		printUsage()
		return
	case "init":
		if err := config.InitializeConfig(); err != nil {
			log.Fatalf("Failed to initialize config: %v", err)
		}
	case "verify":
		if err := config.VerifyConfig(); err != nil {
			log.Fatalf("Verification failed: %v", err)
		}
	case "run":
		if len(args) < 2 {
			log.Fatal("Usage: multiRip run <path_to_jobs.yml>")
		}
		jobsFile := args[1]

		appConfig, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("Failed to load app config: %v", err)
		}

		jobsConfig, err := config.LoadJobs(jobsFile)
		if err != nil {
			log.Fatalf("Failed to load jobs config: %v", err)
		}

		if err := ripper.RunJobs(appConfig, jobsConfig); err != nil {
			log.Fatalf("Job execution failed: %v", err)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Print(`Usage: multiRip <command> [arguments]
Commands:
  init         Create a default config file
  verify       Verify the configuration and paths
  run <path>   Execute rip jobs from a jobs file
`)
}
