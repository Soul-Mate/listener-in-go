package main

import (
	"io/ioutil"
	"encoding/json"
	"bytes"
	"strconv"
	"log"
	"errors"
	"runtime/debug"
)

var MatchJsonNone = errors.New("match json none")

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

// 解析match文件
func ParseMatchFile(file string) ([]Match, error) {
	content, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, err
	}

	if len(content) <= 0 {
		return nil, MatchJsonNone
	}

	result := make([]Match, 32)
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

// 解析文件并保存
func ParseMatchSave(file string) {
	mas, err := ParseMatchFile(file)
	if err == nil {
		SaveMatchMysql(&mas)
	}
}

// 保存比赛到Mysql
func SaveMatchMysql(mas *[]Match) {
	if len(*mas) <= 0 {
		return
	}
	sql := saveMatchSql(mas)
	db := GetConnect()
	_, err := db.Exec(sql)
	Check(err)
}

// 保存比赛的sql
func saveMatchSql(mas *[]Match) string {
	masCnt := len(*mas)
	fieldsCnt := len(MatchFields)
	var fieldsBuf bytes.Buffer
	var updateBuf bytes.Buffer
	for i, matchField := range MatchFields {
		// [a,b,c] -> a,b,c
		fieldsBuf.WriteString(matchField)
		// [a,b,c] -> a=VALUES(a),b=VALUES(b),c=VALUES(c)
		updateBuf.WriteString(matchField)
		updateBuf.WriteString("=VALUES(")
		updateBuf.WriteString(matchField)
		updateBuf.WriteString(")")
		if i != fieldsCnt-1 {
			fieldsBuf.WriteString(",")
			updateBuf.WriteString(",")
		}
	}

	var buf bytes.Buffer
	buf.WriteString("INSERT INTO `")
	buf.WriteString("radar_matches")
	buf.WriteString("`(")
	buf.WriteString(fieldsBuf.String())
	buf.WriteString(")VALUES")

	// (),(),()
	for i, ma := range *mas {
		buf.WriteString("(")

		// "11221"
		buf.WriteString(strconv.Itoa(ma.Id))
		buf.WriteString(",")

		// "2:0"
		buf.WriteString(strconv.Quote(ma.Score))
		buf.WriteString(",")

		// "http://example.com"
		buf.WriteString(strconv.Quote(ma.StreamURL))
		buf.WriteString(",")

		buf.WriteString(strconv.Itoa(ma.Type))
		buf.WriteString(",")

		buf.WriteString(BoolToStr(ma.Visible))
		buf.WriteString(",")

		buf.WriteString(BoolToStr(ma.Suspended))
		buf.WriteString(",")

		buf.WriteString(strconv.Itoa(ma.Status))
		buf.WriteString(",")

		buf.WriteString(strconv.Itoa(ma.SportId))
		buf.WriteString(",")

		buf.WriteString(strconv.Itoa(ma.TournamentId))
		buf.WriteString(",")

		buf.WriteString(strconv.Itoa(ma.HomeTeamId))
		buf.WriteString(",")

		// "example"
		buf.WriteString(strconv.Quote(ma.HomeTeamName))
		buf.WriteString(",")

		buf.WriteString(strconv.Itoa(ma.AwayTeamId))
		buf.WriteString(",")

		// "example"
		buf.WriteString(strconv.Quote(ma.AwayTeamName))
		buf.WriteString(",")

		// "outright"
		buf.WriteString(strconv.Quote(ma.OutrightName))
		buf.WriteString(",")

		// "2008-08-02T11:22:32"
		buf.WriteString(strconv.Quote(ma.StartTime))
		buf.WriteString(",")

		// "2008-08-02T11:22:32"
		buf.WriteString(strconv.Quote(ma.EndTime))
		buf.WriteString(") ")
		if i != masCnt-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString("ON DUPLICATE KEY UPDATE ")
	buf.WriteString(updateBuf.String())
	buf.WriteString(";")
	return buf.String()
}
