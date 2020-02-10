package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/go-redis/redis/v7"
	// "log"
	"path"
	"time"
	
	"encoding/json"
)

type SiteInfo struct {
	LongUrl  string  `redis:"longurl" json:"longurl"`
	Count  int  `redis:"count" json:"count"`

}

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:         ":6379",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/getstats/{id}", handler).Methods("GET")
	http.ListenAndServe(":8080", r)

}


func handler(w http.ResponseWriter, r *http.Request) {

	urlSuffix := path.Base(r.URL.String())
	fmt.Println(urlSuffix)

	s, _ := rdb.Get(urlSuffix).Result()
	if (s != ""){
		var resInfo SiteInfo
		json.Unmarshal([]byte(s), &resInfo)
		countres := resInfo.Count + 1
		json, err := json.Marshal(SiteInfo{LongUrl: resInfo.LongUrl, Count: countres})
		if err != nil {
			fmt.Println(err)
		}
		err = rdb.Set(urlSuffix, json, 0).Err()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		json, err := json.Marshal(SiteInfo{LongUrl: "Elliot", Count: 1})
		if err != nil {
			fmt.Println(err)
		}
		err = rdb.Set(urlSuffix, json, 0).Err()
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Fprintf(w, s)
}
