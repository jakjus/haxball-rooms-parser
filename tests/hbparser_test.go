package hbparser_test

import (
	"github.com/jakjus/hbparser/cmd"
	"testing"
)

func TestGet(t *testing.T) {
	body, err := cmd.GetData()
	if len(body) < 1 {
		t.Fatalf(`Data downloaded from Haxball API is empty.`)
	}
	if err != nil {
		t.Fatalf(`%v`, err)
	}
}

func TestGetAndParse(t *testing.T) {
	body, _ := cmd.GetData()
	serverList := cmd.Parse(body)
	for _, s := range serverList {
		if s.Name == "" {
			t.Fatalf(`Field "name" cannot be nil.`)
		}
		if len(s.Link) < 1 {
			t.Fatalf(`Field "link" cannot be nil.`)
		}
	}
}
