package main

import (
	"io/ioutil"
	"encoding/json"
	"bytes"
	"strconv"
	"fmt"
)

type Match struct {
	Id           int    `json:"Id"`
	Score        string `json:"Score"`
	StreamURL    string `json:"StreamURL"`
	Type         int    `json:"Type"`
	Visible      bool   `json:"Visible"`
	Suspended    bool   `json:"Suspended"`
	Status       int    `json:"Status"`
	SportId      int    `json:"SportId"`
	TournamentId int    `json:"TournamentId"`
	HomeTeamId   int    `json:"HomeTeamId"`
	HomeTeamName string `json:"HomeTeamName"`
	AwayTeamId   int    `json:"AwayTeamId"`
	AwayTeamName string `json:"AwayTeamName"`
	OutrightName string `json:"OutrightName"`
	StartTime    string `json:"StartTime"`
	EndTime      string `json:"EndTime"`
}

var MatchFields = []string{
	"number", "score", "stream_url", "type",
	"visible", "suspended", "status", "sport_id",
	"tournament_id", "home_team_id", "home_team_name", "away_team_id",
	"away_team_name", "outright_name", "start_time", "end_time",
}

func (ma Match) SaveMysql() {
	db := GetConnect()
	sql := GetInsertSql("radar_matches", MatchFields...)
	_, err := db.Exec(sql, ma.Id, ma.Score, ma.StreamURL, ma.Type,
		ma.Visible, ma.Suspended, ma.Status, ma.SportId,
		ma.TournamentId, ma.HomeTeamId, ma.HomeTeamName, ma.AwayTeamName,
		ma.AwayTeamName, ma.OutrightName, ma.StartTime, ma.EndTime)
	Check(err)
}

func SaveMysql(mas []Match) {
	masCnt := len(mas)
	fieldsCnt := len(MatchFields)
	var fieldsBuf bytes.Buffer
	var updateBuf bytes.Buffer
	for i,matchField := range MatchFields {
		// [a,b,c] -> a,b,c
		fieldsBuf.WriteString(matchField)
		// [a,b,c] -> a=VALUES(a),b=VALUES(b),c=VALUES(c)
		updateBuf.WriteString(matchField)
		updateBuf.WriteString("=VALUES(")
		updateBuf.WriteString(matchField)
		updateBuf.WriteString(")")
		if i != fieldsCnt {
			fieldsBuf.WriteString(",")
			updateBuf.WriteString(",")
		}
	}

	var buf bytes.Buffer
	buf.WriteString("INSERT INTO `")
	buf.WriteString("radar_matches")
	buf.WriteString("`(id,")
	buf.WriteString(fieldsBuf.String())
	buf.WriteString(")VALUES")
	for i,ma := range mas {
		buf.WriteString("(")
		buf.WriteString(strconv.Itoa(ma.Id))
		buf.WriteString(",")
		buf.WriteString(ma.Score)
		buf.WriteString(",")
		buf.WriteString(ma.StreamURL)
		buf.WriteString(",")
		buf.WriteString(strconv.Itoa(ma.Type))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatBool(ma.Visible))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatBool(ma.Suspended))
		buf.WriteString(",")
		buf.WriteString(strconv.Itoa(ma.Status))
		buf.WriteString(",")
		buf.WriteString(strconv.Itoa(ma.SportId))
		buf.WriteString(",")
		buf.WriteString(strconv.Itoa(ma.TournamentId))
		buf.WriteString(",")
		buf.WriteString(strconv.Itoa(ma.HomeTeamId))
		buf.WriteString(",")
		buf.WriteString(ma.HomeTeamName)
		buf.WriteString(",")
		buf.WriteString(strconv.Itoa(ma.AwayTeamId))
		buf.WriteString(",")
		buf.WriteString(ma.AwayTeamName)
		buf.WriteString(",")
		buf.WriteString(ma.OutrightName)
		buf.WriteString(",")
		buf.WriteString(ma.StartTime)
		buf.WriteString(",")
		buf.WriteString(ma.EndTime)
		buf.WriteString(")")
		if i!=masCnt{
			buf.WriteString(",")
		}
	}
	buf.WriteString("ON DUPLICATE KEY UPDATE ")
	buf.WriteString(updateBuf.String())
	buf.WriteString(";")
	fmt.Println(buf.String())
	}

func saveSql(mas []*Match) {

}

func ParseMatchFile(file string) ([]Match, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var result []Match
	if content[0] == '"' {
		var m Match
		err = json.Unmarshal(content, &m)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	} else if content[0] == '[' {
		err = json.Unmarshal(content, &result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func ParseMatchSave(file string) {
	mas, err := ParseMatchFile(file)
	Check(err)
	for _, ma := range mas {
		ma.SaveMysql()
	}
}
