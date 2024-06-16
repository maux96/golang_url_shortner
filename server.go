package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type URLTransformer interface {
	SaveUrl(string) (string, error)
	GetUrl(string) (string, error)
}

func StartServer(urlTransformer URLTransformer) (err error) {

	convHtmlTemp, err := template.New("converterTemplate.html").ParseFiles("converterTemplate.html")
	if err != nil {
		log.Fatalf("An error ocurred related to parsing the templates.\n%s", err.Error())
	}

	http.HandleFunc("GET /", func(rw http.ResponseWriter, request *http.Request) {
		// http.ServeFile(rw, request, "converterTemplate.html")
		err = convHtmlTemp.Execute(rw, nil)
		if err != nil {
			log.Fatalf("An error ocurred related to execute the template.\n%s", err.Error())
		}
	})

	http.HandleFunc("POST /", func(rw http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		urlToEncode := req.PostForm.Get("url")

		code, err := urlTransformer.SaveUrl(urlToEncode)
		if err != nil {
			log.Println(err.Error())
			rw.WriteHeader(500)
			return
		}

		ctx := map[string]string{
			"urlCode": fmt.Sprintf("http://%s/%s", req.Host, code),
		}
		log.Printf("%s --> %s", urlToEncode, ctx["urlCode"])

		err = convHtmlTemp.Execute(rw, ctx)
		if err != nil {
			log.Println(err.Error())
			rw.WriteHeader(500)
			return
		}
	})

	http.HandleFunc("GET /{code}", func(rw http.ResponseWriter, req *http.Request) {
		urlCode := req.PathValue("code")

		originalUrl, err := urlTransformer.GetUrl(urlCode)
		if !strings.HasPrefix(originalUrl, "https://") && !strings.HasPrefix(originalUrl, "http://") {
			originalUrl = "https://" + originalUrl
		}

		if err != nil {
			log.Println(err.Error())
			rw.WriteHeader(404)
		} else {
			http.Redirect(rw, req, originalUrl, http.StatusSeeOther)
		}
	})

	log.Println("Server started at http://0.0.0.0:8000")
	err = http.ListenAndServe(":8000", nil)

	return err
}
