package toolkit

import (
	"fmt"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
	i := uint64(1)
	mp := make(map[string]uint64)
	for {
		i++
		str, e := CreateInviteCode(11112222333, i)
		if e != nil {
			fmt.Println(str, i, e)
			return
		}
		//fmt.Println("~~~:", str, i)
		if len(str) > 4 {
			//fmt.Println("~~~:", str, i)
			//return
		}
		if k, ok := mp[str]; ok {
			fmt.Println("now:", str, i)
			fmt.Println("old:", str, k)
			return
		}

		mp[str] = i
		if len(mp) == 5000000 {
			fmt.Println(str, i)
			break
		}
	}
	fmt.Println(len(mp))
}

func TestArrSplit2(t *testing.T) {
	valmp := map[uint64]rune{
		0: 'A', 1: 'B', 2: 'C', 3: 'D', 4: 'E', 5: 'F', 6: 'G', 7: 'H', 8: 'I',
		9: 'J', 10: 'K', 11: 'L', 12: 'M', 13: 'N', 14: 'P', 15: 'Q', 16: 'R', 17: 'S',
		18: 'T', 19: 'U', 20: 'V', 21: 'W', 22: 'X', 23: 'Y', 24: 'Z', 25: '2', 26: '3', 27: '4',
		28: '5', 29: '6', 30: '7', 31: '8',
	}
	mp := map[string]string{}
	for i := 0; i < 100; i++ {
		list := make([]string, 0)
		for _, v := range valmp {
			list = append(list, string(v))
		}
		str := strings.Join(list, ",")
		if _, ok := mp[str]; ok {
			continue
		}
		mp[str] = ""
		fmt.Print("{")
		for i2 := 0; i2 < len(list); i2++ {
			fmt.Printf("%d:'%s',", i2, list[i2])
		}
		fmt.Println("}")
		if len(mp) == 33 {
			break
		}
	}

}
