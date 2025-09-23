package main

import (
	"fmt"
	"log"
	"multiRip/config"
	"multiRip/ripper"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		return
	}

	command := args[0]
	switch command {
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
	case "scan":
		if len(args) < 2 {
			log.Fatal("Usage: multiRip scan <drive_path>")
		}
		drive := args[1]

		appConfig, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("Failed to load app config: %v", err)
		}

		if err := ripper.ScanDrive(drive, appConfig.Handbrake.Binary); err != nil {
			log.Fatalf("Scanning failed: %v", err)
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
  scan <drive> Scan a drive to find the 'Play All' title
`)
}
