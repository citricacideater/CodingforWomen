package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Recipe struct {
	Id          string            `json:"Id"`
	RecipeName  string            `json:"Recipe Name"`
	Source      string            `json:"Source"`
	PrepTime    string            `json:"Preperation Time"`
	CookTime    string            `json:"Cook Time"`
	ServingSize int               `json:"Serving Size"`
	Ingredients map[string]string `json:"Ingredients"`
	Directions  map[int]string    `json:"Directions"`
	Tags        []string          `json:"Tags"`
}

var data []Recipe

func fetchAllRecipes() {
	info, err := http.Get("https://api.npoint.io/fff0f131782057b16a12")
	if err != nil {
		log.Fatal(err, "ERROR 500: Failed to find url")
	}
	defer info.Body.Close()
	body, err := ioutil.ReadAll(info.Body)
	if err != nil {
		log.Fatal(err, "Failed to read json")
	}

	json.Unmarshal(body, &data)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("static/html/index.html") //path is relative to the root of the project directort
	if err != nil {
		http.Error(w, "ERROR 500", 500)
		return
	}
	//result = result[:len(result)-1]

	ts.Execute(w, data)
}

func recipePage(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("static/html/recipe.html") //path is relative to the root of the project directort
	if err != nil {
		http.Error(w, "ERROR 500", 500)
		return
	}
	id := r.URL.Path[len("/recipe/"):]
	var single Recipe

	for _, v := range data {
		if v.Id == id {
			single = v
			break
		}
	}
	ts.Execute(w, single)

}

func tagPage(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("static/html/index.html")
	if err != nil {
		http.Error(w, "ERROR 500", 500)
		return
	}
	tag := r.URL.Path[len("/tag/"):]
	var collection []Recipe
	var single Recipe

	for _, v := range data {
		found := false
		for _, x := range v.Tags {
			if x == tag {
				single = v
				found = true
			}

		}
		if found {
			collection = append(collection, single)
		}

	}
	if len(collection) > 0 {
		ts.Execute(w, collection)
	} else {
		ts.Execute(w, "No Tag Found")
	}

}

func handleRequest() {
	fetchAllRecipes()
	server := http.NewServeMux()
	style := http.FileServer(http.Dir("static/css"))

	server.Handle("/static/css/", http.StripPrefix("/static/css/", style))

	server.HandleFunc("/", homePage)
	server.HandleFunc("/recipe/", recipePage)
	server.HandleFunc("/tag/", tagPage)

	log.Fatal(http.ListenAndServe(":8000", server))
}

func main() {
	handleRequest()
}
