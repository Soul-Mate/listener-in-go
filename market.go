package main

import (
	"io/ioutil"
	"encoding/json"
)

var MarketFields = []string{
	"number", "name", "match_id", "suspended",
	"status", "is_live", "visible", "odds",
}

type Market struct {
	Id        int                      `json:"Id"`
	Name      string                   `json:"Name"`
	MatchId   int                      `json:"MatchId"`
	Suspended bool                     `json:"Suspended"`
	Status    int                      `json:"Status"`
	IsLive    bool                     `json:"IsLive"`
	Visible   bool                     `json:"Visible"`
	Odds      []map[string]interface{} `json:"Odds"`
}

func ParseMarketFile(file string) ([]Market, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var result []Market
	if content[0] == '"' {
		var m Market
		err = json.Unmarshal(content, &m)
		if err != nil {
			return nil, err
		}
		result = append(result,m)
	} else if content[0] == '[' {
		err = json.Unmarshal(content, &result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (mk Market) SaveMysql() {
	db := GetConnect()
	sql := GetInsertSql("radar_markets", MarketFields...)
	odds, err := json.Marshal(mk.Odds)
	if err != nil {
		odds = nil
	}
	_, err = db.Exec(sql, mk.Id, mk.Name, mk.MatchId, mk.Suspended,
		mk.Status, mk.IsLive, mk.Visible, string(odds))
	Check(err)
}
