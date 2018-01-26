package main

import (
	"io/ioutil"
	"encoding/json"
)

type League struct {
	Id           interface{}              `json:"Id"`
	CategoryName string                   `json:"CategoryName"`
	SportId      interface{}              `json:"SportId"`
	Tournaments  []map[string]interface{} `json:"Tournaments"`
}

var LeagueFields = []string{
	"number", "name", "sport_id",
}

var LeagueTournamentFields = []string{
	"number", "name", "league_number",
}

func ParseLeagueFile(file string) ([]League, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var result []League
	if content[0] == '"' {
		var l League
		err = json.Unmarshal(content, &l)
		if err != nil {
			return nil, err
		}
		result = append(result, l)
	} else if content[0] == '[' {
		err = json.Unmarshal(content, &result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (le League) SaveMysql() {
	db := GetConnect()
	leSql := GetInsertSql("radar_leagues", LeagueFields...)
	leTourSql := GetInsertSql("radar_league_tournaments", LeagueTournamentFields...)
	_, err := db.Exec(leSql, le.Id, le.CategoryName, le.SportId)
	Check(err)
	for _, v := range le.Tournaments {
		var leTourSlice []interface{}

		if elem, ok := v["Id"]; ok {
			leTourSlice = append(leTourSlice, elem)
		} else {
			leTourSlice = append(leTourSlice, nil)
		}

		if elem, ok := v["Name"]; ok {
			leTourSlice = append(leTourSlice, elem)
		} else {
			leTourSlice = append(leTourSlice, nil)
		}

		if elem, ok := v["CategoryId"]; ok {
			leTourSlice = append(leTourSlice, elem)
		} else {
			leTourSlice = append(leTourSlice, nil)
		}
		_, err := db.Exec(leTourSql, leTourSlice...)
		Check(err)
	}
}
