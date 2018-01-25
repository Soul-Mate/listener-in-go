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

func (c *Config) NewConfig() error {
	path, err := filepath.Abs(filepath.Dir("./"))
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(path + "/config.json")
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, c); err != nil {
		return err
	}
	return nil
}
