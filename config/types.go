package config

type HandbrakeConfig struct {
	Binary string          `yaml:"binary"`
	DVD    HandbrakePreset `yaml:"dvd"`
	Bluray HandbrakePreset `yaml:"bluray"`
}

type HandbrakePreset struct {
	Preset     string `yaml:"preset"`
	PresetName string `yaml:"preset_name"`
}

type OutputConfig struct {
	ShowsFilenameTemplate  string `yaml:"shows_filename_template"`
	MoviesFilenameTemplate string `yaml:"movies_filename_template"`
}

type MakeMKVConfig struct {
	ScratchDir string `yaml:"scratch_dir"`
}

type Config struct {
	Handbrake HandbrakeConfig `yaml:"handbrake"`
	Output    OutputConfig    `yaml:"output"`
	MakeMKV   MakeMKVConfig   `yaml:"makemkv"`
}

type Show struct {
	Name    string `yaml:"name"    json:"name"`
	Season  int    `yaml:"season"  json:"season"`
	Episode int    `yaml:"episode" json:"episode"`
	Title   int    `yaml:"title"   json:"title"`
}

type Movie struct {
	Name  string `yaml:"name"  json:"name"`
	Title int    `yaml:"title" json:"title"`
}

type DiscType string

const (
	DVD    DiscType = "dvd"
	Bluray DiscType = "bluray"
)

type JobDefinition struct {
	Drive     string        `yaml:"drive"                json:"drive"`
	DiscType  DiscType      `yaml:"disc_type,omitempty"  json:"discType"`
	OutputDir string        `yaml:"output_dir,omitempty" json:"outputDir"`
	Shows     []Show        `yaml:"shows,omitempty"      json:"shows"`
	Movies    []Movie       `yaml:"movies,omitempty"     json:"movies"`
}

type JobsConfig struct {
	OutputDir string          `yaml:"output_dir" json:"outputDir"`
	Jobs      []JobDefinition `yaml:"jobs"       json:"jobs"`
}