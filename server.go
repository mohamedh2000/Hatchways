package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
)

type Success struct {
	Success bool `json:"success"`
}
type Error struct {
	Error string `json:"error"`
}
type Post struct {
	Id         int      `json:"id"`
	Author     string   `json:"author"`
	AuthorId   int      `json:"authorId"`
	Likes      int      `json:"likes"`
	Popularity float64  `json:"popularity"`
	Reads      int      `json:"reads"`
	Tags       []string `json:"tags"`
}
type Posts struct {
	Posts []Post `json:"posts"`
}

func (a Posts) Len() int { return len(a.Posts) }
func (a Posts) Less(i, j int) bool {
	return a.Posts[i].Id < a.Posts[j].Id
}
func (a Posts) Swap(i, j int) { a.Posts[i], a.Posts[j] = a.Posts[j], a.Posts[i] }

const URL = "https://api.hatchways.io/assessment/blog/posts"

func main() {

	http.HandleFunc("/api/posts", ApiPost)

	http.HandleFunc("/api/ping", ApiPing)

	port := ":8080"

	fmt.Println("server is running on port" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func ApiPing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	temp := Success{Success: true}
	ret, _ := json.Marshal(temp)
	w.Write(ret)
}

func ApiPost(w http.ResponseWriter, r *http.Request) {
	queryVals := r.URL.Query()
	var sortBy string = "id"
	var direction string = "asc"

	if !queryVals.Has("tags") {
		w.WriteHeader(400)
		var response, _ = json.Marshal(Error{Error: "Tags parameter is required"})
		w.Write(response)
		return
	}
	if queryVals.Has("sortBy") {
		switch queryVals.Get("sortBy") {
		case "id":
			break
		case "reads":
			sortBy = "reads"
			break
		case "likes":
			sortBy = "likes"
			break
		case "popularity":
			sortBy = "popularity"
			break
		default:
			w.WriteHeader(400)
			var response, _ = json.Marshal(Error{Error: "sortBy parameter is invalid"})
			w.Write(response)
			return

		}
	}
	if queryVals.Has("direction") {
		switch queryVals.Get("direction") {
		case "asc":
			break
		case "desc":
			sortBy = "desc"
			break
		default:
			w.WriteHeader(400)
			var response, _ = json.Marshal(Error{Error: "direction parameter is invalid"})
			w.Write(response)
			return
		}
	}

	temp := strings.Split(queryVals.Get("tags"), ",")
	returnMap := make(map[int]Post)
	for i, s := range temp {
		resp, err := http.Get(URL + "?tag=" + string(temp[i]))
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var tempRet Posts

		json.Unmarshal([]byte(bytes), &tempRet)
		var allRet []Post = tempRet.Posts

		for k := range allRet {
			newEntry := allRet[k]
			returnMap[newEntry.Id] = newEntry //removes duplicates for multiple tags
		}
		fmt.Println(s)
	}

	var retArray []Post

	for _, v := range returnMap {
		retArray = append(retArray, v)
	}

	sort.Slice(retArray, func(i, j int) bool {
		switch sortBy {
		case "id":
			if direction == "asc" {
				return retArray[i].Id < retArray[j].Id
			}
			return retArray[i].Id > retArray[j].Id
		case "reads":
			if direction == "asc" {
				return retArray[i].Reads < retArray[j].Reads
			}
			return retArray[i].Reads > retArray[j].Reads
		case "popularity":
			if direction == "asc" {
				return retArray[i].Popularity < retArray[j].Popularity
			}
			return retArray[i].Popularity > retArray[j].Popularity
		case "likes":
			if direction == "asc" {
				return retArray[i].Likes < retArray[j].Likes
			}
			return retArray[i].Likes > retArray[j].Likes
		default:
			return false
		}
	})

	jsonStr, err := json.Marshal(Posts{Posts: retArray})

	fmt.Println(string(jsonStr))
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonStr)
}
