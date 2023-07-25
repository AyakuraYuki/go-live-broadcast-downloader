package main

import (
	"fmt"
	"github.com/json-iterator/go/extra"
	cjson "go-live-broadcast-downloader/plugins/json"
	"go-live-broadcast-downloader/plugins/sqlnull2json"
	"time"
)

func main() {
	extra.RegisterFuzzyDecoders()

	//s := sqlnull2json.NullString{}
	//s.String = "a"
	//s.Valid = true
	//bs, _ := cjson.JSON.Marshal(s)
	////bs, _ := json.Marshal(mp)
	//fmt.Println(string(bs))
	//
	//s2 := sqlnull2json.NullString{}
	//cjson.JSON.UnmarshalFromString(string(bs), &s2)
	//fmt.Printf("%+v\n", s2)

	s := sqlnull2json.NullTimeUnixMilli{}
	s.Time = time.Now()
	s.Valid = true
	bs, _ := cjson.JSON.Marshal(s)
	fmt.Println(string(bs))

	s2 := sqlnull2json.NullTimeUnixMilli{}
	_ = cjson.JSON.UnmarshalFromString(string(bs), &s2)
	fmt.Printf("%+v\n", s2)
	fmt.Printf("%+v\n", s2.Time.UnixMilli())
}
