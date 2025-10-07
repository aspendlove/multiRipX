package ripper

import (
	"log/slog"
	"multiRip/config"
	"multiRip/util"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
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
		outputDir := driveJob.OutputDir
		if outputDir == "" {
			outputDir = jobsConfig.OutputDir
		}

		for _, show := range driveJob.Shows {
			data := map[string]any{
				"title":   show.Name,
				"season":  show.Season,
				"episode": show.Episode,
				"track":   show.Title,
				"drive":   driveJob.Drive,
			}
			filename, err := util.GenerateFilename(appConfig.Output.ShowsFilenameTemplate, data)
			if err != nil {
				logger.Error("Could not generate filename, skipping", "error", err)
				continue
			}
			outputPath := filepath.Join(outputDir, filename+".mkv")

			cmd := makeRip(driveJob.Drive, outputPath, appConfig.Handbrake.Preset, appConfig.Handbrake.PresetName, appConfig.Handbrake.Binary, show.Title)
			jobs[driveJob.Drive] = append(jobs[driveJob.Drive], Job{
				ID:     show.Title,
				Name:   filename,
				Device: driveJob.Drive,
				Cmd:    cmd,
			})
		}

		for _, movie := range driveJob.Movies {
			data := map[string]any{
				"title": movie.Name,
				"track": movie.Title,
				"drive": driveJob.Drive,
			}
			filename, err := util.GenerateFilename(appConfig.Output.MoviesFilenameTemplate, data)
			if err != nil {
				logger.Error("Could not generate filename, skipping", "error", err)
				continue
			}
			outputPath := filepath.Join(outputDir, filename+".mkv")

			cmd := makeRip(driveJob.Drive, outputPath, appConfig.Handbrake.Preset, appConfig.Handbrake.PresetName, appConfig.Handbrake.Binary, movie.Title)
			jobs[driveJob.Drive] = append(jobs[driveJob.Drive], Job{
				ID:     movie.Title,
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
