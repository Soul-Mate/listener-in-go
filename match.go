package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"log"
	"strings"
	"time"
)

var MatchJsonNone = errors.New("match json file none")

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

// 解析match文件
func ParseMatchFile(file string) ([]Match, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// read null
	if len(content) <= 0 {
		return nil, MatchJsonNone
	}

	result := make([]Match, 32)

	// one
	if content[0] == '"' {
		var m Match
		err = json.Unmarshal(content, &m)
		if err != nil {

			return nil, err
		}
		result = append(result, m)
		// array
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
		fmt.Println("save match")
	} else {
		log.Fatal(err)
	}
}

// 保存比赛到Mysql
func SaveMatchMysql(mas *[]Match) {
	if len(*mas) <= 0 {
		return
	}
	sql := saveMatchSql(mas)
	if sql != "" {
		db := GetMysqlConnect()
		res, err := db.Exec(sql)
		if err != nil {
			log.Fatal(err, res)
		}
	}
}

// 保存比赛的sql
func saveMatchSql(mas *[]Match) string {
	fieldsCnt := len(MatchFields)
	var (
		buf       bytes.Buffer
		fieldsBuf bytes.Buffer
		updateBuf bytes.Buffer
		valuesBuf bytes.Buffer
	)
	for i, matchField := range MatchFields {
		// [a,b,c] -> a,b,c
		fieldsBuf.WriteString("`")
		fieldsBuf.WriteString(matchField)
		fieldsBuf.WriteString("`")

		// [a,b,c] -> a=VALUES(a),b=VALUES(b),c=VALUES(c)
		updateBuf.WriteString("`")
		updateBuf.WriteString(matchField)
		updateBuf.WriteString("`=VALUES(`")
		updateBuf.WriteString(matchField)
		updateBuf.WriteString("`)")
		if i != fieldsCnt-1 {
			fieldsBuf.WriteString(",")
			updateBuf.WriteString(",")
		}
	}

	buf.WriteString("INSERT INTO `")
	buf.WriteString("radar_matches")
	buf.WriteString("`(")
	buf.WriteString(fieldsBuf.String())
	buf.WriteString(")VALUES")

	// (),(),()
	for _, ma := range *mas {
		valuesBuf.WriteString(ma.InsetValueSql())
	}
	values := valuesBuf.String()
	if values == "" {
		return ""
	}
	values = strings.TrimRight(values, ",")
	buf.WriteString(values)

	buf.WriteString("ON DUPLICATE KEY UPDATE ")
	buf.WriteString(updateBuf.String())
	buf.WriteString(";")
	return buf.String()
}

func (ma Match) Add8Hour() string {
	t, err := time.Parse("2006-01-02T15:04:05", ma.StartTime)
	if err != nil {
		return ""
	}
	then := time.Date(t.Year(), t.Month(), t.Day(), t.Hour()+8, t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	return then.Format("2006-01-02 15:04:05")
}

// 生成match插入格式的sql
func (ma Match) InsetValueSql() string {
	var buf bytes.Buffer
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
	buf.WriteString(strconv.Quote(ma.Add8Hour()))
	buf.WriteString(",")

	// "2008-08-02T11:22:32"
	buf.WriteString(strconv.Quote(ma.EndTime))
	buf.WriteString(") ")
	buf.WriteString(",")
	return buf.String()
}
