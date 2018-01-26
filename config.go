package main

import (
	"path/filepath"
	"io/ioutil"
	"encoding/json"
)

type ListenerConf struct {
	RootPath string `json:"root"`
	StaticFiles struct {
		League      string `json:"league"`
		MatchFull   string `json:"match_full"`
		MarketFull  string `json:"market_full"`
		LeaguesFull string `json:"leagues_full"`
	} `json:"static_files"`
}

type DbConf struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Charset  string `json:"charset"`
}

type RedisConf struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Config struct {
	Db       DbConf       `json:"db"`
	Redis    RedisConf    `json:"redis"`
	Listener ListenerConf `json:"listener"`
}

var conf *Config

func NewConfig(c Config) (*Config, error) {
	path, err := filepath.Abs(filepath.Dir("./"))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(path + "/config.json")
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func GetConfig() (*Config) {
	var err error
	if conf == nil {
		conf, err = NewConfig(Config{})
		if err != nil {
			panic(err)
		}
	}
	return conf
}

