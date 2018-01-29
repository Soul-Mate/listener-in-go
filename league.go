package main

import (
	"io/ioutil"
	"encoding/json"
	"bytes"
	"strconv"
	"errors"
	"log"
)

var LeagueJsonNone = errors.New("league file none")

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

func (le League) SaveMysql() {
	db := GetMysqlConnect()
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

func ParseLeagueFile(file string) ([]League, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if len(content) <= 0 {
		return nil, LeagueJsonNone
	}

	result := make([]League, 320)
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

func ParseLeagueSave(file string) {
	les, err := ParseLeagueFile(file)
	if err == nil {
		SaveLeagueMysql(&les)
	} else {
		log.Fatal(err)
	}
}

func SaveLeagueMysql(les *[]League) {
	db := GetMysqlConnect()
	leSql, tourSql := saveLeagueSql(les)
	_, err := db.Exec(leSql)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(tourSql)
	if err != nil {
		log.Fatal(err)
	}
}

func saveLeagueSql(les *[]League) (string, string) {
	var (
		leBuf       bytes.Buffer
		leFieldBuf  bytes.Buffer
		leUpdateBuf bytes.Buffer
		leFieldCnt  = len(LeagueFields)

		leTourBuf       bytes.Buffer
		leTourFieldBuf  bytes.Buffer
		leTourUpdateBuf bytes.Buffer
		leTourFieldCnt  = len(LeagueTournamentFields)
		lesCnt          = len(*les)
	)
	// leagues
	leBuf.WriteString("INSERT INTO `radar_leagues`(")
	for i, field := range LeagueFields {
		// [a,b] -> `a`,`b`
		leFieldBuf.WriteString("`")
		leFieldBuf.WriteString(field)
		leFieldBuf.WriteString("`")

		// [a,b] -> a=VALUES(a), b=VALUES(b)
		leUpdateBuf.WriteString(field)
		leUpdateBuf.WriteString("=VALUES(")
		leUpdateBuf.WriteString(field)
		leUpdateBuf.WriteString(")")
		if i != leFieldCnt-1 {
			leFieldBuf.WriteString(",")
			leUpdateBuf.WriteString(",")
		}
	}
	leBuf.WriteString(leFieldBuf.String())
	leBuf.WriteString(")VALUES")

	// league tournaments
	leTourBuf.WriteString("INSERT INTO `radar_league_tournaments`(")
	for i, field := range LeagueTournamentFields {
		// [a,b] -> `a`,`b`
		leTourFieldBuf.WriteString("`")
		leTourFieldBuf.WriteString(field)
		leTourFieldBuf.WriteString("`")

		// [a,b] -> a=VALUES(a), b=VALUES(b)
		leTourUpdateBuf.WriteString(field)
		leTourUpdateBuf.WriteString("=VALUES(")
		leTourUpdateBuf.WriteString(field)
		leTourUpdateBuf.WriteString(")")
		if i != leTourFieldCnt-1 {
			leTourFieldBuf.WriteString(",")
			leTourUpdateBuf.WriteString(",")
		}
	}
	leTourBuf.WriteString(leTourFieldBuf.String())
	leTourBuf.WriteString(")VALUES")

	for i, le := range *les {
		leBuf.WriteString("(")
		leBuf.WriteString(InterfaceToStr(le.Id))
		leBuf.WriteString(",")

		leBuf.WriteString("\"")
		leBuf.WriteString(le.CategoryName)
		leBuf.WriteString("\",")

		leBuf.WriteString(InterfaceToStr(le.SportId))
		leBuf.WriteString(")")

		// league tournaments
		leTourBuf.WriteString(tournamentsSql(le.Tournaments))

		if i != lesCnt-1 {
			leBuf.WriteString(",")
			leTourBuf.WriteString(",")
			leBuf.WriteString("\n")
		}
	}
	// leagues
	leBuf.WriteString("ON DUPLICATE KEY UPDATE ")
	leBuf.WriteString(leUpdateBuf.String())
	leBuf.WriteString(";")

	// league tournaments
	leTourBuf.WriteString("ON DUPLICATE KEY UPDATE ")
	leTourBuf.WriteString(leTourUpdateBuf.String())
	leTourBuf.WriteString(";")
	return leBuf.String(), leTourBuf.String()
}

func tournamentsSql(tours []map[string]interface{}) string {
	var leTourBuf bytes.Buffer
	toursCnt := len(tours)
	for i, tour := range tours {
		leTourBuf.WriteString("(")

		if elem, ok := tour["Id"]; ok {
			leTourBuf.WriteString(InterfaceToStr(elem))
		}
		leTourBuf.WriteString(",")

		if elem, ok := tour["Name"]; ok {
			leTourBuf.WriteString(strconv.Quote(InterfaceToStr(elem)))
		}
		leTourBuf.WriteString(",")

		if elem, ok := tour["CategoryId"]; ok {
			leTourBuf.WriteString(InterfaceToStr(elem))
		}
		leTourBuf.WriteString(")")

		if i != toursCnt-1 {
			leTourBuf.WriteString(",")
			leTourBuf.WriteString("\n")
		}
	}
	return leTourBuf.String()
}
