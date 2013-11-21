package main

import (
	"fmt"
	"io/ioutil"
)

type Page struct {
	Title string
	Body  []byte
}

func (self *Page) save() error {
	filename := self.Title + ".txt"
	return ioutil.WriteFile(filename, self.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"

	if body, err := ioutil.ReadFile(filename); err != nil {
		return nil, fmt.Errorf("loadPage: %e", err)
	} else {
		return &Page{Title: title, Body: body}, nil
	}
}

func main() {
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
}
