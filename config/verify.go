package config

import (
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

func VerifyConfig() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	fmt.Println("Current Configuration:")
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("could not marshal config for display: %w", err)
	}
	fmt.Println(string(data))

	fmt.Println("Verifying paths...")
	valid := true

	if err := checkExecutable(cfg.Handbrake.Binary); err != nil {
		fmt.Printf("  - Handbrake Binary: %s (Error: %v)\n", cfg.Handbrake.Binary, err)
		valid = false
	} else {
		fmt.Printf("  - Handbrake Binary: %s (OK)\n", cfg.Handbrake.Binary)
	}

	if err := checkFile(cfg.Handbrake.Bluray.Preset); err != nil {
		fmt.Printf("  - Preset File: %s (Error: %v)\n", cfg.Handbrake.Bluray.Preset, err)
		valid = false
	} else {
		fmt.Printf("  - Preset File: %s (OK)\n", cfg.Handbrake.Bluray.Preset)
	}

	if err := checkFile(cfg.Handbrake.DVD.Preset); err != nil {
		fmt.Printf("  - Preset File: %s (Error: %v)\n", cfg.Handbrake.DVD.Preset, err)
		valid = false
	} else {
		fmt.Printf("  - Preset File: %s (OK)\n", cfg.Handbrake.DVD.Preset)
	}

	if cfg.MakeMKV.ScratchDir != "" {
		if err := checkDirectory(cfg.MakeMKV.ScratchDir); err != nil {
			fmt.Printf("  - Scratch Directory: %s (Error: %v)\n", cfg.MakeMKV.ScratchDir, err)
			valid = false
		} else {
			fmt.Printf("  - Scratch Directory: %s (OK)\n", cfg.MakeMKV.ScratchDir)
		}
	}

	if !valid {
		return fmt.Errorf("configuration verification failed")
	}

	fmt.Println("\nConfiguration OK.")
	return nil
}

func checkFile(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("file does not exist")
	}
	if err != nil {
		return fmt.Errorf("could not stat file: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file")
	}
	return nil
}

func checkExecutable(path string) error {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return fmt.Errorf("path is a directory, not a file")
		}
		if info.Mode().Perm()&0111 != 0 {
			return nil
		}
	}

	if _, err := exec.LookPath(path); err != nil {
		return fmt.Errorf("not found in PATH and not an executable file")
	}

	return nil
}

func checkDirectory(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist")
	}
	if err != nil {
		return fmt.Errorf("could not stat path: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is a file, not a directory")
	}
	return nil
}
