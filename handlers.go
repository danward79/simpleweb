package smplweb

import (
  "net/http"
  "strings"
)

func GeneralHandler(w http.ResponseWriter, r *http.Request) {
  page := &Page{}
  params := r.URL.Query()
  pageTitle := params.Get(":page")
  if pageTitle == "" {
    pageTitle = "index"
  }
  
  page.Title = strings.Title(pageTitle)
  RenderTemplate(w, pageTitle + ".tmpl", page)
}

func GeneralContentHandler(w http.ResponseWriter, r *http.Request) {
  page := &Page{}
  params := r.URL.Query()
  pageTitle := params.Get(":page")
  if pageTitle == "" {
    pageTitle = "index"
  }
  
  pageContents := params.Get(":contents")
  if pageContents == "" {
    pageTitle = "index"
  } else {
    var err error
    page, err = LoadPage(pageContents)
    if err != nil {
      http.Redirect(w, r, "/", http.StatusFound)
      return
    }
  }
  
  page.Title = strings.Title(pageTitle)
  RenderTemplate(w, pageTitle + ".tmpl", page)
}