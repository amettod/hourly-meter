package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/amettod/hourly-meter/internal/app"
)

//go:embed template
var pages embed.FS

func main() {
	addr := flag.Int("addr", 8008, "server port address")
	flag.Parse()

	http.Handle("/", allowMethod(setForm, http.MethodGet))
	http.Handle("/run", allowMethod(runApp, http.MethodPost))

	log.Printf("go to http://localhost:8008/")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *addr), nil))
}

func setForm(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	templateParse(w, nil, "form_page.tmpl", "base_layout.tmpl")
}

func runApp(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		httpError(w, http.StatusBadRequest, err)
		return
	}

	file, header, err := r.FormFile("filename")
	if err != nil {
		httpError(w, http.StatusBadRequest, err)
		return
	}
	defer file.Close()

	f, err := os.CreateTemp("", header.Filename)
	if err != nil {
		httpError(w, http.StatusInternalServerError, err)
		return
	}
	defer os.Remove(f.Name())

	_, err = io.Copy(f, file)
	if err != nil {
		httpError(w, http.StatusInternalServerError, err)
		return
	}

	contract := r.PostForm.Get("contract")
	name := r.PostForm.Get("name")
	meter := r.PostForm.Get("meter")
	coefficient, err := strconv.ParseFloat(r.PostForm.Get("coefficient"), 64)
	if err != nil {
		httpError(w, http.StatusBadRequest, err)
		return
	}

	app, err := app.New(f.Name(), contract, name, meter, coefficient)
	if err != nil {
		httpError(w, http.StatusInternalServerError, err)
		return
	}

	err = app.Run()
	if err != nil {
		httpError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	result := map[string]string{
		"Total":  fmt.Sprintf("%.2f", app.Total),
		"Values": fmt.Sprintf("%d", len(app.Rows)),
	}

	templateParse(w, result, "result_page.tmpl", "base_layout.tmpl")
}

func allowMethod(next http.HandlerFunc, method string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.Header().Set("Allow", method)
			httpError(w, http.StatusMethodNotAllowed, errors.New(http.StatusText(http.StatusMethodNotAllowed)))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func httpError(w http.ResponseWriter, status int, err error) {
	log.Println(err)
	http.Error(w, http.StatusText(status), status)
}

func addTemplatesPrefix(file []string) []string {
	var files []string
	for _, f := range file {
		files = append(files, "template/"+f)
	}

	return files
}

func templateParse(w http.ResponseWriter, data interface{}, files ...string) {
	temp, err := template.ParseFS(pages, addTemplatesPrefix(files)...)
	if err != nil {
		httpError(w, http.StatusInternalServerError, err)
		return
	}

	err = temp.Execute(w, data)
	if err != nil {
		httpError(w, http.StatusInternalServerError, err)
	}
}
