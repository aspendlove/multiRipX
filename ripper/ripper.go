package ripper

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"multiRip/config"
	"multiRip/util"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
)

type JobExec func(*os.File) error

type Job struct {
	ID     int
	Name   string
	Device string
	Cmd    JobExec
}

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func RunJobs(appConfig *config.Config, jobsConfig *config.JobsConfig) error {
	jobs, err := GetJobs(appConfig, jobsConfig)
	if err != nil {
		return err
	}
	if err = ExecuteJobs(jobs, nil); err != nil {
		return err
	}

	logger.Info("All jobs completed.")
	return nil
}

func GetJobs(appConfig *config.Config, jobsConfig *config.JobsConfig) (map[string][]Job, error) {
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

			cmd, err := makeRip(driveJob.Drive, outputPath, appConfig, show.Title, driveJob.DiscType)
			if err != nil {
				logger.Error("Could not generate job, skipping", "error", err)
				continue
			}
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

			cmd, err := makeRip(driveJob.Drive, outputPath, appConfig, movie.Title, driveJob.DiscType)
			if err != nil {
				logger.Error("Could not generate job, skipping", "error", err)
				continue
			}
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
		return nil, nil
	}

	return jobs, nil
}

func ExecuteJobs(jobs map[string][]Job, onLog func(driveId, message string)) error {
	var workerGroup sync.WaitGroup
	workerGroup.Add(len(jobs))
	var jobError error = nil

	for device, jobList := range jobs {
		go func(device string, jobs []Job, wg *sync.WaitGroup) {
			defer wg.Done()

			logFilename := filepath.Base(device) + ".log"
			logFile, err := os.Create(logFilename)
			if err != nil {
				logger.Error("Error creating log file", "file", logFilename, "error", err)
				jobError = fmt.Errorf("Error creating log file")
				return
			}
			defer logFile.Close()

			pr, pw, _ := os.Pipe()

			logFinished := make(chan bool)

			go func() {
				tee := io.TeeReader(pr, logFile)

				scanner := bufio.NewScanner(tee)
				for scanner.Scan() {
					if onLog != nil {
						onLog(device, scanner.Text()+"\n")
					}
				}
			}()

			for _, job := range jobs {
				logger.Info("Worker started job", "device", device, "job_id", job.ID, "name", job.Name)

				if err := job.Cmd(pw); err != nil {
					logger.Error("Error while transcoding", "device", device, "job_id", job.ID, "error", err)
					jobError = fmt.Errorf("Error while transcoding")
					continue
				}

				logger.Info("Worker finished job", "device", device, "job_id", job.ID)
			}

			pw.Close()

			<-logFinished
			pr.Close()
		}(device, jobList, &workerGroup)
	}

	workerGroup.Wait()
	return jobError
}

func makeRip(device, filename string, appConfig *config.Config, title int, diskType config.DiscType) (JobExec, error) {
	switch diskType {
	case config.DVD:
		return makeDvdRip(device, filename, appConfig, title), nil
	case config.Bluray:
		return makeBlurayRip(device, filename, appConfig, title), nil
	}
	return nil, fmt.Errorf("Invalid Disk Type")
}

func makeDvdRip(device, filename string, appConfig *config.Config, title int) JobExec {
	args := []string{
		"--preset-import-file", appConfig.Handbrake.DVD.Preset,
		"--preset", appConfig.Handbrake.DVD.PresetName,
		"-i", device,
		"-t", strconv.Itoa(title),
		"-o", filename,
	}

	var command JobExec = func(file *os.File) error {
		cmd := exec.Command(appConfig.Handbrake.Binary, args...)
		cmd.Stdout = file
		cmd.Stderr = file
		return cmd.Run()
	}

	return command
}

func makeBlurayRip(device, filename string, appConfig *config.Config, title int) JobExec {
	var command JobExec = func(file *os.File) error {
		tempDir, err := os.MkdirTemp(appConfig.MakeMKV.ScratchDir, "multirip_bluray_*")
		if err != nil {
			return fmt.Errorf("Cannot create temporary directory: %w", err)
		}
		defer os.RemoveAll(tempDir)
		args := []string{
			"mkv",
			"dev:" + device,
			strconv.Itoa(title),
			tempDir,
		}
		ripCmd := exec.Command("makemkvcon", args...)
		args = []string{
			"--preset-import-file", appConfig.Handbrake.Bluray.Preset,
			"--preset", appConfig.Handbrake.Bluray.PresetName,
			"-i", filepath.Join(tempDir, "temp.mkv"),
			"-o", filename,
		}
		transcodeCmd := exec.Command(appConfig.Handbrake.Binary, args...)
		ripCmd.Stdout = file
		ripCmd.Stderr = file
		transcodeCmd.Stdout = file
		transcodeCmd.Stderr = file
		err = ripCmd.Run()
		if err != nil {
			return fmt.Errorf("Cannot Rip BluRay: %w", err)
		}
		err = transcodeCmd.Run()
		if err != nil {
			return fmt.Errorf("Cannot Transcode BluRay: %w", err)
		}
		return nil
	}
	return command
}
