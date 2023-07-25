package dateparse

import (
	"testing"
	"time"
)

func TestTsToDateYm(t *testing.T) {
	ts := time.Now()
	t.Log(TsToDateYmd(ts))
	t.Log(TsToDateYm(ts))
	t.Log(TsToDateYW(ts))
}

func TestSecondFriendly(t *testing.T) {
	t.Log(SecondFriendly(100))
	t.Log(SecondFriendly(1))
	t.Log(SecondFriendly(22))
	t.Log(SecondFriendly(3600))
	t.Log(SecondFriendly(60))
	t.Log(SecondFriendly(478))
	t.Log(SecondFriendly(3645))
	t.Log(SecondFriendly(360450010623))
}

func TestGetFirstDateOfWeek(t *testing.T) {
	mt := time.Now().AddDate(0, 0, +3)
	t.Log(mt)
	t.Log(GetFirstDateOfWeek(mt))
	t.Log(GetFirstDateOfNextWeek(mt))
}

func TestDayFriendly(t *testing.T) {
	ct, _ := ParseYmdHis("2022-05-12 02:07:28")

	sec := time.Now().Unix() - ct.Unix()
	t.Log("sec:", sec)
	t.Log(DayFriendly(sec))
}
