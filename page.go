package smplweb

import (
  "github.com/russross/blackfriday"
  "html/template"
  "io/ioutil"
  "net/http"
  "errors"
  "os"
  "path/filepath"
  "log"
  "strings"
)

// Page structure for documents
type Page struct {
  Title string
  Body template.HTML
}

// Read file contents
func readFile(path string) ([]byte, error){
  body, err := ioutil.ReadFile (path)
  if err != nil {
    return nil, err
  }
  return body, nil
}

// Load a Markdown file into a Page struct
func loadMarkDownPage(title string) (*Page, error) {
  body, err := readFile(paths.Static + title + ".md")
  if err != nil {
    return nil, err
  }
  bodyHtml := template.HTML(blackfriday.MarkdownCommon(body))
  return &Page{Title: strings.Title(title), Body:bodyHtml}, nil
}

// Load a TXT file into a Page struct
func loadTxtPage(title string) (*Page, error) {
  body, err := readFile(paths.Static + title + ".txt")
  if err != nil {
    return nil, err
  }
  
  bodyHtml := template.HTML(string(body))
  return &Page{Title: strings.Title(title), Body:bodyHtml}, nil
}

// Load a HTML file into a Page struct
func loadHTMLPage(title string) (*Page, error) {
  body, err := readFile(paths.Static + title + ".html")
  if err != nil {
    return nil, err
  }
  
  bodyHtml := template.HTML(string(body))
  return &Page{Title: strings.Title(title), Body:bodyHtml}, nil
}

// Handlers should call LoadPage, which should manage the page load regardless of storage type. Which could be HTML, txt or markdown
func LoadPage(title string) (*Page, error){
  path := paths.Static + title

  if _, err := os.Stat(path + ".html"); err == nil {
    return loadHTMLPage(title)
  } else if _, err = os.Stat(path + ".md"); err == nil {
    return loadMarkDownPage(title)
  } else if _, err = os.Stat(path + ".txt"); err == nil {
    return loadTxtPage(title) 
  }
  return nil, errors.New("Invalid Page Document")
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
  m := GetPath(r.URL.Path)
  if m == nil {
    http.NotFound(w, r)
    return "", errors.New("Invalid Page Title")
  }
  return m[2], nil
}

var templates map[string]*template.Template

// From the layouts and includes builde templates.
func CreateTemplates () {
  if templates == nil {
    templates = make(map[string]*template.Template)
  }
  
  layouts, err := filepath.Glob(paths.Template + "layouts/*.tmpl")
  if err != nil {
    log.Fatal(err)
  }
    
  includes, err := filepath.Glob(paths.Template + "includes/*.tmpl")
  if err != nil {
    log.Fatal(err)
  }
  
  for _, layout := range layouts {
    files := append(includes, layout)
    templates[filepath.Base(layout)] = template.Must(template.ParseFiles(files...))
  }  
}

// Render given template on response
func RenderTemplate(w http.ResponseWriter, tmpl string, p *Page) {  
  if t, ok := templates[tmpl]; ok {
    err := t.ExecuteTemplate(w, "base", p)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }
  }
  
  return
}