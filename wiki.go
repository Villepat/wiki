package main

import (
    "html/template"
    "os"
    "net/http"
	"log"
)
// Let's start by defining the data structures.
// A wiki consists of a series of interconnected pages,
// each of which has a title and a body (the page content).
// Here, we define Page as a struct with two fields representing the title and body.
type Page struct {
    Title string
    Body  []byte
}
// The Body element is a []byte rather than string because that is the type expected
// by the io libraries we will use, as you'll see below.

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func (p *Page) save() error { //method named save
    filename := p.Title + ".txt"
    return os.WriteFile(filename, p.Body, 0600) //creates a file with p.Title as filename
}
// This is a method named save that takes as its receiver p, a pointer to Page .
// It takes no parameters, and returns a value of type error.
// This method will save the Page's Body to a text file.
// For simplicity, we will use the Title as the file name.
func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}
func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/view/"):]
    p, err := loadPage(title) //return *Page, err
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound) //nonexistant file -> edit
        return
    }
    renderTemplate(w, "view", p)
}
func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl+".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
func saveHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/save/"):]
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
    // p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
    // p1.save()
    // p2, _ := loadPage("TestPage")
    // fmt.Println(string(p2.Body))
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/edit/", editHandler)
    http.HandleFunc("/save/", saveHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}