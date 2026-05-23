package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const configDirName = "multiRip"
const configFileName = "config.yml"

func GetConfigFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("could not get user config directory: %w", err)
	}
	return filepath.Join(configDir, configDirName, configFileName), nil
}

func InitializeConfig() error {
	path, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		fmt.Printf("Config file already exists: %s\n", path)
		return nil
	}

	defaultConfig := Config{
		Handbrake: HandbrakeConfig{
			Binary: "/path/to/HandBrakeCLI",
			DVD: HandbrakePreset{
				Preset:     "/path/to/dvd-presets.json",
				PresetName: "My_DVD_Preset",
			},
			Bluray: HandbrakePreset{
				Preset:     "/path/to/bluray-presets.json",
				PresetName: "My_Bluray_Preset",
			},
		},
		Output: OutputConfig{
			ShowsFilenameTemplate:  "S{season}E{episode} - {title}",
			MoviesFilenameTemplate: "{title}",
			CDFolderTemplate:       "%A - %d",
			CDFilenameTemplate:     "%t - %n",
		},
		MakeMKV: MakeMKVConfig{
			ScratchDir: "",
		},
	}

	data, err := yaml.Marshal(&defaultConfig)
	if err != nil {
		return fmt.Errorf("could not marshal default config: %w", err)
	}

	configDir := filepath.Dir(path)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	fmt.Printf("Created default config file: %s\n", path)
	return nil
}

func LoadConfig() (*Config, error) {
	path, err := GetConfigFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w. Please run 'multiRip init'", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)
	}

	return &cfg, nil
}

func LoadJobs(path string) (*JobsConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read jobs file: %w", err)
	}

	var jobsCfg JobsConfig
	if err := yaml.Unmarshal(data, &jobsCfg); err != nil {
		return nil, fmt.Errorf("could not parse jobs file: %w", err)
	}

	for i := range jobsCfg.Jobs {
		if jobsCfg.Jobs[i].DiscType == "" {
			jobsCfg.Jobs[i].DiscType = DVD
		}
	}

	return &jobsCfg, nil
}
