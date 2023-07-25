package sqlnull2json

import (
	"database/sql"
	cjson "go-live-broadcast-downloader/plugins/json"
	"time"
)

type NullString struct {
	sql.NullString
}

func (v NullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return cjson.JSON.Marshal(v.String)
	} else {
		return cjson.JSON.Marshal(nil)
	}
}

func (v *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := cjson.JSON.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.String = *s
	} else {
		v.Valid = false
	}
	return nil
}

type NullTime struct {
	sql.NullTime
}

func (v NullTime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return cjson.JSON.Marshal(v.Time)
	} else {
		return cjson.JSON.Marshal(nil)
	}
}

func (v *NullTime) UnmarshalJSON(data []byte) error {
	var s *time.Time
	if err := cjson.JSON.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.Time = *s
	} else {
		v.Valid = false
	}
	return nil
}

type NullTimeUnixMilli struct {
	sql.NullTime
}

func (v NullTimeUnixMilli) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return cjson.JSON.Marshal(v.Time.UnixMilli())
	} else {
		return cjson.JSON.Marshal(nil)
	}
}

func (v *NullTimeUnixMilli) UnmarshalJSON(data []byte) error {
	var s *int64
	if err := cjson.JSON.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.Time = time.Unix(*s/1000, (*s%1000)*1000*1000)
	} else {
		v.Valid = false
	}
	return nil
}

type NullBool struct {
	sql.NullBool
}

func (v NullBool) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return cjson.JSON.Marshal(v.Bool)
	} else {
		return cjson.JSON.Marshal(nil)
	}
}

func (v *NullBool) UnmarshalJSON(data []byte) error {
	var s *bool
	if err := cjson.JSON.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.Bool = *s
	} else {
		v.Valid = false
	}
	return nil
}

type NullByte struct {
	sql.NullByte
}

func (v NullByte) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return cjson.JSON.Marshal(v.Byte)
	} else {
		return cjson.JSON.Marshal(nil)
	}
}

func (v *NullByte) UnmarshalJSON(data []byte) error {
	var s *byte
	if err := cjson.JSON.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.Byte = *s
	} else {
		v.Valid = false
	}
	return nil
}

type NullFloat64 struct {
	sql.NullFloat64
}

func (v NullFloat64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return cjson.JSON.Marshal(v.Float64)
	} else {
		return cjson.JSON.Marshal(nil)
	}
}

func (v *NullFloat64) UnmarshalJSON(data []byte) error {
	var s *float64
	if err := cjson.JSON.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.Float64 = *s
	} else {
		v.Valid = false
	}
	return nil
}

type NullInt16 struct {
	sql.NullInt16
}

func (v NullInt16) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return cjson.JSON.Marshal(v.Int16)
	} else {
		return cjson.JSON.Marshal(nil)
	}
}

func (v *NullInt16) UnmarshalJSON(data []byte) error {
	var s *int16
	if err := cjson.JSON.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.Int16 = *s
	} else {
		v.Valid = false
	}
	return nil
}

type NullInt32 struct {
	sql.NullInt32
}

func (v NullInt32) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return cjson.JSON.Marshal(v.Int32)
	} else {
		return cjson.JSON.Marshal(nil)
	}
}

func (v *NullInt32) UnmarshalJSON(data []byte) error {
	var s *int32
	if err := cjson.JSON.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.Int32 = *s
	} else {
		v.Valid = false
	}
	return nil
}

type NullInt64 struct {
	sql.NullInt64
}

func (v NullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return cjson.JSON.Marshal(v.Int64)
	} else {
		return cjson.JSON.Marshal(nil)
	}
}

func (v *NullInt64) UnmarshalJSON(data []byte) error {
	var s *int64
	if err := cjson.JSON.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.Int64 = *s
	} else {
		v.Valid = false
	}
	return nil
}
