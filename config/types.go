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
	CDFolderTemplate       string `yaml:"cd_foldername_template"`
	CDFilenameTemplate     string `yaml:"cd_filename_template"`
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
	Name    string `yaml:"name"`
	Season  int    `yaml:"season"`
	Episode int    `yaml:"episode"`
	Title   int    `yaml:"title"`
}

type Movie struct {
	Name  string `yaml:"name"`
	Title int    `yaml:"title"`
}

type CD struct {
	Musicbrainz string `yaml:"musicbrainz"`
	Offset      int    `yaml:"offset"`
}

type DiscType string

const (
	DVD     DiscType = "dvd"
	Bluray  DiscType = "bluray"
	CD_Disc DiscType = "cd"
)

type JobDefinition struct {
	Drive     string   `yaml:"drive"`
	DiscType  DiscType `yaml:"disc_type,omitempty"`
	OutputDir string   `yaml:"output_dir,omitempty"`
	Shows     []Show   `yaml:"shows,omitempty"`
	Movies    []Movie  `yaml:"movies,omitempty"`
	CD        []CD     `yaml:"cds,omitempty"`
}

type JobsConfig struct {
	OutputDir string          `yaml:"output_dir"`
	Jobs      []JobDefinition `yaml:"jobs"`
}
