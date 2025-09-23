package ripper

import (
	"bytes"
	"fmt"
	"log/slog"
	"multiRip/config"
	"multiRip/util"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Job struct {
	ID     int
	Name   string
	Device string
	Cmd    *exec.Cmd
}

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func RunJobs(appConfig *config.Config, jobsConfig *config.JobsConfig) error {
	jobs := make(map[string][]Job)

	for _, driveJob := range jobsConfig.Jobs {
		// Determine output directory
		outputDir := driveJob.OutputDir
		if outputDir == "" {
			outputDir = jobsConfig.OutputDir
		}

		// Process shows
		for _, show := range driveJob.Shows {
			data := map[string]interface{}{
				"title":   show.Title,
				"season":  show.Season,
				"episode": show.Episode,
				"track":   show.TrackNumber,
				"drive":   driveJob.Drive,
			}
			filename, err := util.GenerateFilename(appConfig.Output.ShowsFilenameTemplate, data)
			if err != nil {
				logger.Error("Could not generate filename, skipping", "error", err)
				continue
			}
			outputPath := filepath.Join(outputDir, filename+".mkv")

			cmd := makeRip(driveJob.Drive, outputPath, appConfig.Handbrake.Preset, appConfig.Handbrake.PresetName, appConfig.Handbrake.Binary, show.TrackNumber)
			jobs[driveJob.Drive] = append(jobs[driveJob.Drive], Job{
				ID:     show.TrackNumber,
				Name:   filename,
				Device: driveJob.Drive,
				Cmd:    cmd,
			})
		}

		// Process movies
		for _, movie := range driveJob.Movies {
			data := map[string]interface{}{
				"title": movie.Title,
				"track": movie.TrackNumber,
				"drive": driveJob.Drive,
			}
			filename, err := util.GenerateFilename(appConfig.Output.MoviesFilenameTemplate, data)
			if err != nil {
				logger.Error("Could not generate filename, skipping", "error", err)
				continue
			}
			outputPath := filepath.Join(outputDir, filename+".mkv")

			cmd := makeRip(driveJob.Drive, outputPath, appConfig.Handbrake.Preset, appConfig.Handbrake.PresetName, appConfig.Handbrake.Binary, movie.TrackNumber)
			jobs[driveJob.Drive] = append(jobs[driveJob.Drive], Job{
				ID:     movie.TrackNumber,
				Name:   filename,
				Device: driveJob.Drive,
				Cmd:    cmd,
			})
		}
	}

	if len(jobs) == 0 {
		logger.Warn("No jobs found to process.")
		return nil
	}

	var workerGroup sync.WaitGroup
	workerGroup.Add(len(jobs))

	for device, jobList := range jobs {
		go func(device string, jobs []Job, wg *sync.WaitGroup) {
			defer wg.Done()

			logFilename := filepath.Base(device) + ".log"
			logFile, err := os.Create(logFilename)
			if err != nil {
				logger.Error("Error creating log file", "file", logFilename, "error", err)
				return
			}
			defer logFile.Close()

			for _, job := range jobs {
				logger.Info("Worker started job", "device", device, "job_id", job.ID, "name", job.Name)
				job.Cmd.Stdout = logFile
				job.Cmd.Stderr = logFile

				if err := job.Cmd.Run(); err != nil {
					logger.Error("Error while transcoding", "device", device, "job_id", job.ID, "error", err)
					continue
				}

				logger.Info("Worker finished job", "device", device, "job_id", job.ID)
			}
		}(device, jobList, &workerGroup)
	}

	workerGroup.Wait()
	logger.Info("All jobs completed.")
	return nil
}

func makeRip(device, filename, presetFile, presetName, binary string, title int) *exec.Cmd {
	args := []string{
		"--preset-import-file", presetFile,
		"--preset", presetName,
		"-i", device,
		"-t", strconv.Itoa(title),
		"-o", filename,
	}

	return exec.Command(binary, args...)
}

// ScanDrive executes a scan on the specified drive to find the "Play All" title.
func ScanDrive(drivePath string, handbrakeBinary string) error {
	logger.Info("Scanning drive to find 'Play All' title...", "drive", drivePath)

	cmd := exec.Command(handbrakeBinary, "-i", drivePath, "--title", "0")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		// HandBrakeCLI exits with error code 1 after a scan, which is expected.
		// We just log the error and continue with parsing the output.
		logger.Warn("HandBrakeCLI scan finished with an error (this is often expected)", "error", err)
	}

	logger.Info("Full HandBrakeCLI scan output:", "output", out.String())

	title, err := findPlayAllTitle(out.String())
	if err != nil {
		return fmt.Errorf("could not determine 'Play All' title: %w", err)
	}

	logger.Info("Found 'Play All' title", "title_number", title)
	fmt.Printf("The 'Play All' title is: %d\n", title)
	return nil
}

// findPlayAllTitle parses the HandBrakeCLI scan output to find the title with the longest duration.
func findPlayAllTitle(scanOutput string) (int, error) {
	logger.Info("Parsing HandBrakeCLI output...", "raw_output", scanOutput)
	titleRegex := regexp.MustCompile(`\+ title (\d+):\s+.*duration: (\d{2}:\d{2}:\d{2})`)
	matches := titleRegex.FindAllStringSubmatch(scanOutput, -1)
	logger.Info("Found title matches", "count", len(matches))

	if len(matches) == 0 {
		return 0, fmt.Errorf("no titles found in scan output")
	}

	var longestDuration time.Duration
	var playAllTitle int

	for _, match := range matches {
		titleNum, _ := strconv.Atoi(match[1])
		durationStr := match[2]

		parts := strings.Split(durationStr, ":")
		hours, _ := strconv.Atoi(parts[0])
		minutes, _ := strconv.Atoi(parts[1])
		seconds, _ := strconv.Atoi(parts[2])

		duration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second

		if duration > longestDuration {
			longestDuration = duration
			playAllTitle = titleNum
		}
	}

	if playAllTitle == 0 {
		return 0, fmt.Errorf("could not identify a title with the longest duration")
	}

	return playAllTitle, nil
}
