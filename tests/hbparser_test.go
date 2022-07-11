package hbparser_test

import (
	"github.com/jakjus/hbparser"
	"testing"
)

func TestGet(t *testing.T) {
	body, err := hbparser.GetData()
	if len(body) < 1 {
		t.Fatalf(`Data downloaded from Haxball API is empty.`)
	}
	if err != nil {
		t.Fatalf(`%v`, err)
	}
}

func TestGetAndParse(t *testing.T) {
	body, _ := hbparser.GetData()
	serverList := hbparser.Parse(body)
	for _, s := range serverList {
		if s.Name == nil {
			t.Fatalf(`Field "name" cannot be nil.`)
		}
		if s.Link == nil {
			t.Fatalf(`Field "link" cannot be nil.`)
		}
	}
}

func ExampleMain() {
	body, _ := hbparser.GetData()
	serverList := hbparser.Parse(body)
	for _, s := range serverList {
            s.Print()
	}
}
