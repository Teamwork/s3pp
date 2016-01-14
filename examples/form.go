package main

import (
	"html/template"
	"net/http"
	"time"

	"github.com/pborman/uuid"
	"github.com/teamwork/s3pp"
)

const (
	accessKeyID     = "<YOUR-ACCESS-KEY-ID>"
	secretAccessKey = "<YOUR-SECRET-ACCESS-KEY>"
	region          = "us-east-1"
	bucket          = "mybucket"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("form").Parse(templ)
	if err != nil {
		panic(err)
	}

	form, err := s3pp.New(s3pp.Config{
		AWSCredentials: s3pp.AWSCredentials{
			AccessKeyID:     accessKeyID,
			SecretAccessKey: secretAccessKey,
		},
		Bucket:  bucket,
		Region:  region,
		Expires: 10 * time.Minute,
		Key:     s3pp.Match("key", uuid.New()),
		Conditions: []s3pp.Condition{
			s3pp.Match("acl", "public-read"),
			s3pp.Match("success_action_status", "201"),
		},
	})
	if err != nil {
		panic(err)
	}

	t.Execute(w, form)
}

const templ = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>Form Upload Example</title>
    <style>input { display: block; width: 100%; padding: 5px; margin: 5px; }</style>
  </head>
  <body>
    <form action="{{.Action}}" method="POST" enctype="multipart/form-data">
      {{range $name, $value := .Fields}}
	<input type="text" name="{{$name}}" value="{{$value}}" readonly>
      {{end}}

      <input type="file" name="file">
      <input type="submit" value="Upload">
    </form>
  </body>
</html>
`
