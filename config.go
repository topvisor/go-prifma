package proxy

import (
	"encoding/json"
	"os"
)

type Config struct {
	AccessLogFile *string     `json:"accessLogFile,omitempty"`
	ErrorLogFile  *string     `json:"errorLogFile,omitempty"`
	ListenPort    int        `json:"listenPort"`
	ListenType    ListenType `json:"listenType"`
	Filters       *[]Filter   `json:"filters,omitempty"`
}

func NewConfig(filename string) (*Config, error) {
	var config Config
	err := config.LoadFromFile(filename)

	return &config, err
}

func (t *Config) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(t)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (t *Config) ListenFile(filename string) {

}
