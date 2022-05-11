package models

type Secrets struct {
	Name     string
	Internal []string `yaml:",omitempty"`
	External []string `yaml:",omitempty"`
}

type Conf struct {
	Flags []KV        `yaml:",omitempty"`
	Env   []KV        `yaml:",omitempty"`
	Files []ConfFiles `yaml:",omitempty"`
}

type KV struct {
	Name  string
	Value string
}

type ConfFiles struct {
	Name        string
	Path        string
	Filename    string `yaml:"fileName,omitempty"`
	Include     bool
	Constraints []KV `yaml:",omitempty"`
}

type CMD struct {
	Tag       string
	Platform  string
	ShaGit    string
	Version   string
	BuildDate string
	Push      bool
	Bin       string
}

type Dockerfile struct {
	Lang      string
	Name      string
	CMD       string
	Source    string
	Target    string
	Overwrite bool
	ShaGit    string
}
