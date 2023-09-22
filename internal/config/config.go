package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Config struct {
	Postgres struct {
		ConnStr string `json:"connect_string"`
		Port    string `json:"port"`
	}
}

var (
	_config     *Config
	_onceConfig sync.Once
)

// GetConfig получение объекта конфига
func GetConfig() *Config {
	_onceConfig.Do(func() { _config = &Config{}; _config.load() })
	return _config
}

func (cnf *Config) load() {
	cnf.loadCNF("../../configPostgres.json", &cnf.Postgres)
}

func (cnf *Config) loadCNF(filename string, data interface{}) {
	if file, err := os.Open(filename); err != nil {
		fmt.Println("не могу открыть файл с конфигурацией ", filename, " ", err)
		os.Exit(1)
	} else {
		defer file.Close()

		d := json.NewDecoder(file)
		if err = d.Decode(data); err != nil {
			fmt.Println("не могу загрузить конфигурацию ", filename, err)
			os.Exit(1)
		} else {
			fmt.Println("конфигурация загружена ", filename)
		}
	}
}
