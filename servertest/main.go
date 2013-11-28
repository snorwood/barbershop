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
var parseFilePath = regexp.MustCompile("^[a-zA-Z0-9]+")

type Page struct {
	Title string
	Body  template.HTML
}

func (self *Page) save() error {
	filename := dataFolder + self.Title + ".txt"
	return ioutil.WriteFile(filename, []byte(self.Body), 0600)
}

func getPageNames(directory string) ([]string, error) {
	files, err := ioutil.ReadDir(dataFolder)

	if err != nil {
		return nil, fmt.Errorf("getPageNames: %v", err)
	}

	names := make([]string, len(files))

	for i, f := range files {
		names[i] = parseFilePath.FindStringSubmatch(f.Name())[0]
	}

	return names, nil
}

func sortByLength(strings []string) {
	sorted := false
	count := 0

	for !sorted {
		sorted = true
		for i, s := range strings {
			if len(s) < len(strings[count]) {
				sorted = false
				temp := strings[count]
				strings[count] = strings[i]
				strings[i] = temp
			}
		}
		count += 1
	}
}

func ConvertPageNameToLink(name []byte) []byte {
	return []byte(fmt.Sprintf("<a href=\"/view/%s\">%s</a>", name, name))
}

func ConvertPageNamesToLinks(text []byte, ignore []string) ([]byte, error) {
	names, err := getPageNames(dataFolder)
	if err != nil {
		return nil, fmt.Errorf("ConvertPageNamesToLinks: ", err)
	}
	sortByLength(names)
	filteredNames := make([]string, 0)

	for _, name := range names {
		valid := true

		for _, i := range ignore {
			if name == i {
				valid = false
			}
		}

		if valid {
			filteredNames = append(filteredNames, name)
		}
	}

	for _, name := range filteredNames {

		findName := regexp.MustCompile(name)
		text = findName.ReplaceAllFunc(text, ConvertPageNameToLink)
	}

	return text, nil
}

func loadPage(title string) (*Page, error) {
	filename := dataFolder + title + ".txt"

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("loadPage: %e", err)
	}

	return &Page{Title: title, Body: template.HTML(body)}, nil
}

func loadHTMLPage(title string) (*Page, error) {
	filename := dataFolder + title + ".txt"

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("loadPage: %e", err)
	}

	ignore := []string{title}
	body, err = ConvertPageNamesToLinks(body, ignore)
	if err != nil {
		return nil, fmt.Errorf("loadPage: %e", err)
	}

	return &Page{Title: title, Body: template.HTML(body)}, nil
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

	if p, err = loadHTMLPage(title); err != nil {
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

	p := &Page{Title: title, Body: template.HTML(body)}
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

func resourcesHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Println(path)
	http.ServeFile(w, r, path)
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
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.HandleFunc("/", makeHandler(rootHandler))
	http.ListenAndServe(":8080", nil)
}
