package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Article struct {
	Title string
}

type Response struct {
	Data []Article
}

func homePage(w http.ResponseWriter, r *http.Request) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	req, _ := http.NewRequest("GET",
		"https://newsapi.org/v2/top-headlines?country=us", nil)

	req.Header.Add("X-Api-Key", os.Getenv("NEWS-APIKEY"))

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf(err.Error())
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	response := string(body)
	responseBytes := []byte(response)
	var jsonRes map[string]interface{}
	_ = json.Unmarshal(responseBytes, &jsonRes)

	arr := jsonRes["articles"].([]interface{})
	var articles []Article

	for i := range arr {
		cur := arr[i].(map[string]interface{})
		a := Article{cur["title"].(string)}
		articles = append(articles, a)
	}

	respArticles := Response{articles}
	b, err := json.Marshal(respArticles)
	_, err = fmt.Fprintf(w, string(b))
	if err != nil {
		return
	}
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	handleRequests()
}
