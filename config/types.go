package config

type HandbrakeConfig struct {
	Binary     string `yaml:"binary"`
	Preset     string `yaml:"preset"`
	PresetName string `yaml:"preset_name"`
}

type OutputConfig struct {
	ShowsFilenameTemplate  string `yaml:"shows_filename_template"`
	MoviesFilenameTemplate string `yaml:"movies_filename_template"`
}

type Config struct {
	Handbrake HandbrakeConfig `yaml:"handbrake"`
	Output    OutputConfig    `yaml:"output"`
}

type Show struct {
	Title       string `yaml:"title"`
	Season      int    `yaml:"season"`
	Episode     int    `yaml:"episode"`
	TrackNumber int    `yaml:"track_number"`
}

type Movie struct {
	Title       string `yaml:"title"`
	TrackNumber int    `yaml:"track_number"`
}

type JobDefinition struct {
	Drive     string  `yaml:"drive"`
	OutputDir string  `yaml:"output_dir,omitempty"`
	Shows     []Show  `yaml:"shows,omitempty"`
	Movies    []Movie `yaml:"movies,omitempty"`
}

type JobsConfig struct {
	OutputDir string          `yaml:"output_dir"`
	Jobs      []JobDefinition `yaml:"jobs"`
}
