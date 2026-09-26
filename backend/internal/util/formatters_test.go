package util

import "testing"

func TestAppStatusText(t *testing.T) {
	cases := map[string]string{
		"planning": "规划中", "admitted": "已录取", "waitlisted": "候补名单", "x": "未知",
	}
	for in, want := range cases {
		if got := AppStatusText(in); got != want {
			t.Errorf("AppStatusText(%s) = %s, want %s", in, got, want)
		}
	}
}

func TestDocTypeText(t *testing.T) {
	if got := DocTypeText("ps"); got != "个人陈述" {
		t.Errorf("DocTypeText = %s", got)
	}
}

func TestMaterialStatusText(t *testing.T) {
	if got := MaterialStatusText("approved"); got != "已审核" {
		t.Errorf("MaterialStatusText = %s", got)
	}
}
