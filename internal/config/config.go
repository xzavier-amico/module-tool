package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Job struct {
	Name    string   `yaml:"name"`
	Source  string   `yaml:"source"`
	From    string   `yaml:"from"`
	Targets []string `yaml:"targets"`
	Icon    string   `yaml:"icon"`
	Desc    string   `yaml:"desc"`
	Author  string   `yaml:"author"`
	OutDir  string   `yaml:"out_dir"`
}

type Config struct {
	SrcDir string `yaml:"src_dir"`
	DistDir string `yaml:"dist_dir"`
	Jobs   []Job  `yaml:"jobs"`
}

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	if c.SrcDir == "" {
		c.SrcDir = "src"
	}
	if c.DistDir == "" {
		c.DistDir = "dist"
	}
	return &c, nil
}

func Ext(format string) string {
	switch format {
	case "surge":
		return ".sgmodule"
	case "loon":
		return ".plugin"
	case "qx":
		return ".conf"
	case "egern":
		return ".yaml"
	case "stash":
		return ".stoverride"
	}
	return ".txt"
}
