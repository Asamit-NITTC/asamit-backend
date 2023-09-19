package test

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func PostTest() {
	v := url.Values{}

	v.Set("id_token", "")
	v.Set("client_id", "")

	req, err := http.PostForm("https://api.line.me/oauth2/v2.1/verify", v)
	if err != nil {
		log.Fatal(err)
	}

	defer req.Body.Close()

	fmt.Println(req)
}
