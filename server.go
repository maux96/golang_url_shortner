package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type URLTransformer interface {
	SaveUrl(string) string
	GetUrl(string) (string, error)
}

func StartServer(urlTransformer URLTransformer) (err error) {

	convHtmlTemp, err := template.New("converterTemplate.html").ParseFiles("converterTemplate.html")
	if err != nil {
		log.Fatalf("An error ocurred related to parsing the templates.\n%s", err.Error())
	}

	http.HandleFunc("/", func(rw http.ResponseWriter, request *http.Request) {
		// http.ServeFile(rw, request, "converterTemplate.html")
		err = convHtmlTemp.Execute(rw, nil)
		if err != nil {
			log.Fatalf("An error ocurred related to execute the template.\n%s", err.Error())
		}
	})

	http.HandleFunc("POST /", func(rw http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		urlToEncode := req.PostForm.Get("url")

		ctx := map[string]string{
			"urlCode": fmt.Sprintf("https://%s/%s", req.Host, urlTransformer.SaveUrl(urlToEncode)),
		}
		log.Printf("%s-->%s", urlToEncode, ctx["urlCode"])

		err = convHtmlTemp.Execute(rw, ctx)
		if err != nil {
			log.Fatalf("An error ocurred related to execute the template.\n%s", err.Error())
		}
	})

	log.Println("Server started at http://0.0.0.0:8000")
	err = http.ListenAndServe(":8000", nil)

	return err
}
