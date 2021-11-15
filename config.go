package main

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	Login    *ScriptConfig `yaml:"login"`
	Playback *ScriptConfig `yaml:"playback"`
}

type ScriptConfig struct {
	File            string `yaml:"file"`
	ActionDelayTime int    `yaml:"actionDelayTime"`
}

type Script struct {
	StartUrl string         `json:"startUrl"`
	Actions  []ScriptAction `json:"actions"`
}

type ScriptAction struct {
	TabUrl string `json:"tabUrl"`
	Type   string `json:"type"`
	XPath  string `json:"xpath"`
	Value  string `json:"value,omitempty"`
}

func LoadConfig() Config {
	yamlFile, err := os.Open("config.yaml")
	if err != nil {
		panic(err)
	}
	defer yamlFile.Close()

	var config Config
	byteValue, _ := ioutil.ReadAll(yamlFile)
	err = yaml.Unmarshal(byteValue, &config)
	if err != nil {
		panic(err)
	}
	return config
}

func LoadScript(config *ScriptConfig) Script {
	jsonFile, err := os.Open(config.File)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	var script Script
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &script)
	if err != nil {
		panic(err)
	}
	return script
}
