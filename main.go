package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	testurl := "https://www.youtube.com/watch?v=ajJaYGnKpu4" // Get your own LTT edition Noctua cooler! - Linus Tech Tips
	parsedurl, err := url.Parse(testurl)
	if err != nil {
		panic(fmt.Sprintln("There shouldn't be an error parsing the url:", err))
	}
	log.Println("The URL searched:", testurl)
	query := parsedurl.Query()
	log.Println("The extracted video id:", query.Get("v"))
	res, err := http.Get(testurl)
	if err != nil {
		panic(fmt.Sprintln("There shouldn't be an error getting the html:", err))
	}
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(fmt.Sprintln("There shouldn't be an error parsing the html:", err))
	}
	//log.Println("The html is: ", string(bs)) // was commented because it was really really unreadable
	prefix := bytes.Index(bs, []byte(".config = ")) + len([]byte(".config = ")) // After doing some debugging regarding the youtube scripts and a nice stack overflow post :)
	postfix := bytes.Index(bs,[]byte(";ytplayer.load"))
	between := bs[prefix:postfix]
	log.Println("The config (in json) is:", string(between))
	var parsedargs map[string]interface{}
	err = json.Unmarshal(between, &parsedargs)
	if err != nil {
		panic(fmt.Sprintln("There shouldn't be an error parsing the json: ", err))
	}
	cfgargs := parsedargs["args"].(map[string]interface{})
	log.Println("Final args:", cfgargs)
	adaptiveFmts := cfgargs["adaptive_fmts"]
	log.Println("adaptive_fmts:", adaptiveFmts)
	adptfmts := strings.Split(adaptiveFmts.(string), "&")
	for _, v := range adptfmts {
		if strings.HasPrefix(v, "url=") {
			ulr, err := url.PathUnescape(v[4:])
			if err != nil {
				panic(fmt.Sprintln("There shouldn't be an error unescaping the url: ", err))
			}
			log.Println("final url:", ulr)
			break
		}
	}

}
