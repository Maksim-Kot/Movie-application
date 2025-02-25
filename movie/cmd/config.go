package main

type config struct {
	API APIConfig `yaml:"api"`
}

type APIConfig struct {
	Port int    `yaml:"port"`
	Name string `yaml:"name"`
}
