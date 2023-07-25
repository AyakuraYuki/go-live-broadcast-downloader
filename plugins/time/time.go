package time

import (
	"database/sql/driver"
	"fmt"
	cjson "go-live-broadcast-downloader/plugins/json"
	"strconv"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05"
	timezone   = "Asia/Shanghai"
)

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	return cjson.JSON.Marshal(t.Time.UnixMilli())
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	val, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		now, err := time.ParseInLocation(`"`+timeFormat+`"`, string(data), time.Local)
		t.Time = now
		return err
	}
	t.Time = time.Unix(val/1000, (val%1000)*1000*1000)
	return nil
}

func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		t.Time = value
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}
