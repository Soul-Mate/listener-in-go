package main

import (
	"encoding/json"
	"io/ioutil"
	"errors"
	"log"
	"fmt"
)

var ErrorJsonFileRead = errors.New("read json file")

var ErrorJsonFileUnmarshal = errors.New("unmarshal json file")

var ErrorJsonUnmarshalNone = errors.New("json none data")

type Parser interface {
	ParseMatchFile(file string) (interface{}, error)
}

func BaseParse(file string) ([]interface{}, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, ErrorJsonFileRead
	}
	if content[0] == '"' || content[0] == '{' {
		result := make([]interface{}, 1)
		err = json.Unmarshal(content, &result[0])
		if err != nil {
			log.Fatal(err)
			return nil, ErrorJsonFileUnmarshal
		}
		return result, nil

	} else if content[0] == '[' {
		var result []interface{}

		err = json.Unmarshal(content, &result)
		if err != nil {
			return nil, ErrorJsonFileUnmarshal
		}
		return result, nil
	}
	return nil, ErrorJsonUnmarshalNone
}

func BaseSave(table string, fields []string, jsonKeys []string, data ... interface{}) {
	//db := GetMysqlConnect()
	//sql := GetInsertSql(table, fields...)
	for _, v := range data {
		vMap := v.(map[string]interface{})
		var args []interface{}
		for _, f := range jsonKeys {
			if elem, ok := vMap[f]; ok {
				if f == "Odds" {
					fmt.Println(json.Marshal(elem.([]map[string]interface{})))
				}
				args = append(args, elem)
			} else {
				args = append(args, nil)
			}
		}
		fmt.Println(args)
		//_, err := db.Exec(sql, args...)
		//Check(err)
	}
}
