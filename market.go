package main

import (
	"io/ioutil"
	"encoding/json"
	"bytes"
	"strconv"
	"log"
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

func (mk Market) marshalOdds() string {
	bs, err := json.Marshal(mk.Odds)
	if err != nil {
		return ""
	}
	return string(bs)
}

func ParseMarketFile(file string) ([]Market, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if len(content) < 0 {
		return nil, nil
	}
	result := make([]Market, 32)
	if content[0] == '"' {
		var m Market
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

func ParseMarketSave(file string) {
	mks, err := ParseMarketFile(file)
	if err != nil {
		log.Fatal(err)
		return
	}
	if mks != nil {
		SaveMarketMysql(&mks)
	}
}

func SaveMarketMysql(mks *[]Market) {
	db := GetConnect()
	sql := saveMarketSql(mks)
	_, err := db.Exec(sql)
	Check(err)
}

func saveMarketSql(mks *[]Market) string {
	var (
		buf       bytes.Buffer
		fieldBuf  bytes.Buffer
		updateBuf bytes.Buffer
		fieldCnt  = len(MarketFields)
		mksCnt    = len(*mks)
	)
	buf.WriteString("INSERT INTO `radar_markets`(")

	for i, field := range MarketFields {
		// [a,b] -> `a`,`b`
		fieldBuf.WriteString("`")
		fieldBuf.WriteString(field)
		fieldBuf.WriteString("`")

		// [a,b] -> a=VALUES(a), b=VALUES(b)
		updateBuf.WriteString(field)
		updateBuf.WriteString("=VALUES(")
		updateBuf.WriteString(field)
		updateBuf.WriteString(")")
		if i != fieldCnt-1 {
			fieldBuf.WriteString(",")
			updateBuf.WriteString(",")
		}
	}

	buf.WriteString(fieldBuf.String())
	buf.WriteString(")VALUES")

	for i, mk := range *mks {
		buf.WriteString("(")
		buf.WriteString(strconv.Itoa(mk.Id))
		buf.WriteString(",")

		buf.WriteString(strconv.Quote(mk.Name))
		buf.WriteString(",")

		buf.WriteString(strconv.Itoa(mk.MatchId))
		buf.WriteString(",")

		buf.WriteString(BoolToStr(mk.Suspended))
		buf.WriteString(",")

		buf.WriteString(strconv.Itoa(mk.Status))
		buf.WriteString(",")

		buf.WriteString(BoolToStr(mk.IsLive))
		buf.WriteString(",")

		buf.WriteString(BoolToStr(mk.Visible))
		buf.WriteString(",")

		// odds map to string
		buf.WriteString(strconv.Quote(mk.marshalOdds()))
		buf.WriteString(")")

		if i != mksCnt-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString("ON DUPLICATE KEY UPDATE ")
	buf.WriteString(updateBuf.String())
	buf.WriteString(";")

	return buf.String()
}
