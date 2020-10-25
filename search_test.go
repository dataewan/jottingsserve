package main

import (
	"net/http"
	"testing"
)

var ind BleveSeachIndex

type Testfile struct {
	Title string
}

func (t Testfile) ToHTML(w http.ResponseWriter) {
}

func (t Testfile) GetTitle() string {
	return t.Title
}
func (t Testfile) GetText() string {
	return t.Title
}
func (t Testfile) GetPath() string {
	return t.Title
}

func TestGetIndex(t *testing.T) {
	ind = GetIndex()
}

func TestPut(t *testing.T) {
	ind.Put("a", Testfile{Title: "alpha alpha"})
	ind.Put("b", Testfile{Title: "alpha"})
	ind.Put("c", Testfile{Title: "pinkie"})
	ind.Put("d", Testfile{Title: "pinky"})

}

func TestSearch(t *testing.T) {
	results := ind.Search("alpha")

	if len(results) != 2 {
		t.Errorf("Should have more results %v", results)
	}
}
