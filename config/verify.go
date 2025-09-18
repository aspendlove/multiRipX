package config

import (
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

// VerifyConfig loads and validates the configuration.
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

	// Check HandBrake binary
	if err := checkExecutable(cfg.Handbrake.Binary); err != nil {
		fmt.Printf("  - Handbrake Binary: %s (Error: %v)\n", cfg.Handbrake.Binary, err)
		valid = false
	} else {
		fmt.Printf("  - Handbrake Binary: %s (OK)\n", cfg.Handbrake.Binary)
	}

	// Check preset file
	if err := checkFile(cfg.Handbrake.Preset); err != nil {
		fmt.Printf("  - Preset File: %s (Error: %v)\n", cfg.Handbrake.Preset, err)
		valid = false
	} else {
		fmt.Printf("  - Preset File: %s (OK)\n", cfg.Handbrake.Preset)
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
	// First, try to stat the file directly.
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return fmt.Errorf("path is a directory, not a file")
		}
		// On Unix-like systems, check for execute permissions.
		if info.Mode().Perm()&0111 != 0 {
			return nil // File exists and is executable.
		}
	}

	// If stat fails or it's not executable, check the system PATH.
	if _, err := exec.LookPath(path); err != nil {
		return fmt.Errorf("not found in PATH and not an executable file")
	}

	return nil // Found in PATH.
}
