package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
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

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]

	var p *Page = new(Page)
	var err error = nil

	if p, err = loadPage(title); err != nil {
		p = &Page{Title: title, Body: []byte(err.Error())}
	}

	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.ListenAndServe(":8080", nil)
}
