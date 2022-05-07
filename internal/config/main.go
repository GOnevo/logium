package config

import (
	"encoding/json"
	"github.com/milad-abbasi/gonfig"
	"io/ioutil"
	"log"
	"os"
)

// FILE is default name of config file
const FILE = "config.json"

type server struct {
	Port int               `json:"port" default:"3333"`
	Last int               `json:"last" default:"5"`
	Auth map[string]string `json:"auth"`
}

// Source type
type Source struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// Config type
type Config struct {
	Server  server   `json:"server"`
	Sources []Source `json:"sources"`
}

// Get returns filled configuration
func Get() *Config {
	var c Config
	err := gonfig.Load().FromFile(FILE).Into(&c)
	if err != nil {
		log.Printf("\n%s is broken", FILE)
		os.Exit(1)
	}

	return &c
}

// Template creates template config file
func Template() {
	c := &Config{}
	c.Server.Port = 1234
	c.Server.Last = 5
	c.Sources = append(c.Sources, Source{Name: "test log", Path: "/path/to/dir/test.log"})
	j, _ := json.MarshalIndent(c, "", "  ")
	err := ioutil.WriteFile(FILE, j, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
