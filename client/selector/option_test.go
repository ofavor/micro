package selector

import "testing"

func Test_VerToInt(t *testing.T) {
	v := "1.0.0"
	if verToInt(v) != 1000000 {
		t.Errorf("%s convert failed:%d", v, verToInt(v))
	}
	v = "5.43.27"
	if verToInt(v) != 5043027 {
		t.Errorf("%s convert failed:%d", v, verToInt(v))
	}
	v = "15.043.127"
	if verToInt(v) != 15043127 {
		t.Errorf("%s convert failed:%d", v, verToInt(v))
	}
}

func Test_IsVersionIn(t *testing.T) {
	if !isVersionIn("1.0.7", []string{}) {
		t.Error("check failed with empty versions")
	}
	if !isVersionIn("1.0.7", []string{"1.0.6"}) {
		t.Error("check failed with little version")
	}
	if !isVersionIn("1.0.7", []string{"", "1.0.8"}) {
		t.Error("check failed with big version")
	}
	if !isVersionIn("1.0.7", []string{"1.0.6", "1.0.8"}) {
		t.Error("check failed with two versions")
	}
}
