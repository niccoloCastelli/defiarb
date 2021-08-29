package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	NodeUrl         string
	ContractAddress string
	ExecuteTrades   bool
	Db              DbConfig
}

type DbConfig struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
	SslMode  bool
}

func (d DbConfig) ConnectionString() string {
	var sslmode string
	if d.SslMode {
		sslmode = "require"
	} else {
		sslmode = "disable"
	}
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", d.Host, d.Port, d.User, d.DbName,
		d.Password, sslmode)
}
func GetDefaultConfig() Config {
	return Config{
		NodeUrl:         "https://bsc-dataseed.binance.org/",
		ContractAddress: "",
		Db: DbConfig{
			Host:     "localhost",
			Port:     "5432",
			DbName:   "bscbot",
			User:     "bscbot",
			Password: "bscbot",
			SslMode:  false,
		},
	}
}

func ReadConfig(fileName string) (*Config, error) {
	cfg := Config{}
	if fileName == "" {
		fileName = GetDefaultConfigLocation()
	}
	File, err := os.OpenFile(fileName, os.O_RDONLY, 0775)
	if err != nil {
		return nil, err
	}
	defer func() { _ = File.Close() }()
	FileContent, err := ioutil.ReadAll(File)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(FileContent, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func GetDefaultConfigLocation() string {
	return "config.json"
}

func WriteConfig(fileName string) error {
	defaultConfig := GetDefaultConfig()
	if fileName == "" {
		fileName = GetDefaultConfigLocation()
	}
	if f, err := ioutil.ReadFile(fileName); err == nil {
		if err := json.Unmarshal(f, &defaultConfig); err != nil {
			return err
		}
	}

	File, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		return err
	}
	defer func() { _ = File.Close() }()

	content, err := json.MarshalIndent(&defaultConfig, "", "\t")
	_, err = File.Write(content)
	return err
}
