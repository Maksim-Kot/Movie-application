package main

type config struct {
	API   APIConfig   `yaml:"api"`
	MySQL MySQLConfig `yaml:"mysql"`
}

type APIConfig struct {
	Port int    `yaml:"port"`
	Name string `yaml:"name"`
}

type MySQLConfig struct {
	Database string `yaml:"database"`
}
