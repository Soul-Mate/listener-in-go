package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"github.com/vmihailenco/msgpack"
)

var MarketJsonNone = errors.New("market file json none data")

var MarketFields = []string{
	"number", "name", "match_id", "suspended",
	"status", "is_live", "visible", "odds",
}

type Odds struct {
	Id        interface{}
	Status    int
	Suspended bool
	Visible   bool
	Name      interface{}
	Title     string
	Value     interface{}
	MarketId  interface{}
	MatchId   interface{}
}

type Market struct {
	Id        int    `json:"Id"`
	Name      string `json:"Name"`
	MatchId   int    `json:"MatchId"`
	Suspended bool   `json:"Suspended"`
	Status    int    `json:"Status"`
	IsLive    bool   `json:"IsLive"`
	Visible   bool   `json:"Visible"`
	Odds      []Odds `json:"Odds"`
	//Odds      []map[string]interface{} `json:"Odds"`
}

// 序列化odds
func (mk Market) MarshalOdds() string {
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

	if len(content) <= 0 {
		return nil, MarketJsonNone
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
	if err == nil {
		SaveMarketMysql(&mks)
		fmt.Println("save market")
	}
}

func SaveMarketMysql(mks *[]Market) {
	db := GetMysqlConnect()
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
		// 滚盘 且 比赛未结束
		if mk.IsLive && mk.Status != 3 {
			// 写入缓存
			mk.Set()
			continue
		} else {
			tmp := mk.Get()
			if tmp != nil {
				mk = *tmp.(*Market)
				mk.Del()
			}
		}

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
		buf.WriteString(strconv.Quote(mk.MarshalOdds()))
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

// 设置market数据存入缓存系统
// key 	 -> market.Id
// value -> msg pack序列化后的结果
func (mk Market) Set() error {
	// 滚盘
	if mk.IsLive {
		bs, err := msgpack.Marshal(mk)
		if err != nil {
			return err
		}
		client := GetRedisConnect()
		// cache to redis
		_, err = client.Set(strconv.Itoa(mk.Id), bs, 0).Result()
		if err != nil {
			return err
		}
	}
	return nil
}

// 从缓存系统中取出market数据
func (mk Market) Get() interface{} {
	client := GetRedisConnect()
	cmd := client.Get(strconv.Itoa(mk.Id))
	_, err := cmd.Result()
	if err != nil {
		return nil
	}
	bs, err := cmd.Bytes()
	if err != nil {
		return nil
	}
	m := new(Market)
	err = msgpack.Unmarshal(bs, m)
	if err != nil {
		return nil
	}
	return m
}

// 从缓存系统中删除market数据
func (mk Market) Del() error {
	client := GetRedisConnect()
	_, err := client.Del(strconv.Itoa(mk.Id)).Result()
	if err != nil {
		return err
	}
	return nil
}
