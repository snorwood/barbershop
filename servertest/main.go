package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

var dataFolder = "data/"
var tmplFolder = "tmpl/"
var templates = template.Must(template.ParseFiles(tmplFolder+"edit.html", tmplFolder+"view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
var validRoot = regexp.MustCompile("^/([a-zA-Z0-9]*)$")
var parseFilePath = regexp.Compile("^[a-zA-Z0-9]+")

type Page struct {
	Title string
	Body  []byte
}

func (self *Page) save() error {
	filename := dataFolder + self.Title + ".txt"
	return ioutil.WriteFile(filename, self.Body, 0600)
}

func getPageNames(directory string) []string {
	files := ioutil.ReadDir(dataFolder)
	names := make([]string, len(files))

	for i, f := range files {
		names[i] = parseFilePath.FindStringSubmatch(f.Name())[0]
	}

	return names
}

func ConvertNamesToPageLinks(text string) {
	names := getPageNames(dataFolder)

	for _, name := range names {
		regexp.
	}
}

func loadPage(title string) (*Page, error) {
	filename := dataFolder + title + ".txt"

	if body, err := ioutil.ReadFile(filename); err != nil {
		return nil, fmt.Errorf("loadPage: %e", err)
	} else {
		return &Page{Title: title, Body: body}, nil
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)

		if m == nil {
			m = validRoot.FindStringSubmatch(r.URL.Path)
		}

		if m == nil {
			http.NotFound(w, r)
			return
		}

		fn(w, r, m[len(m)-1])
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	var p *Page = new(Page)
	var err error

	if p, err = loadPage(title); err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	var p *Page = new(Page)
	var err error

	if p, err = loadPage(title); err != nil {
		p = &Page{Title: title}
	}

	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")

	p := &Page{Title: title, Body: []byte(body)}
	if err := p.save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func rootHandler(w http.ResponseWriter, r *http.Request, title string) {
	if title == "" {
		title = "FrontPage"
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	if err := templates.ExecuteTemplate(w, tmpl+".html", p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/", makeHandler(rootHandler))

	http.ListenAndServe(":8080", nil)
}
