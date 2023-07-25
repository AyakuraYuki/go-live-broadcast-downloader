package toolkit

import "testing"

func TestVersionCompare(t *testing.T) {
	v1 := ""
	v2 := "1.2.3"
	b, t1, t2 := VersionCompare(v1, v2)
	switch b {
	case VersionGt:
		t.Log("v1 > v2")
	case VersionLt:
		t.Log("v1 < v2")
	case VersionEq:
		t.Log("v1 = v2")
	}

	t.Log(t1)
	t.Log(t2)
}
