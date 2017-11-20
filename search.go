package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

type xdcc struct {
	Content []xdccContent `json: "content"`
}

type xdccContent struct {
	Name          string `json: "name"`
	Number        int    `json: "number"`
	BotId         int    `json: "botId"`
	Size          string `json: "size"`
	EpisodeNumber int    `json: "episodeNumber"`
	LastModified  string `json: "lastModified"`
}

func searchMain() {
	query := "gamers"
	episodeNumber := 3
	quality := "720p"
	temp := xdcc{}
	// make a slice for suggestions
	findPacklist(query, episodeNumber, quality, &temp)
	// form the query here once you figure out what the user wants
	// ex. gamers! 9 will return did you mean [HorribleSubs] Gamers! - 09[480].mkv?
	botMain(query, episodeNumber, quality, &temp)
}

func findPacklist(query string, episode int, quality string, x *xdcc) {
	queryString := fmt.Sprintf("https://api.nibl.co.uk:8080/nibl/search?query=%s&episodeNumber=%d", query, episode)
	getJSON(queryString, x)

	if len(x.Content) > 0 {
		fmt.Println("this is a valid query")
		createSuggestion(episode, quality, x)
	} else {
		fmt.Println("this is not a valid query")
		// return this is not a valid query message to the front end
	}
}

func createSuggestion(episode int, quality string, x *xdcc) {
	for _, j := range x.Content {
		// fmt.Println(j.Name)
		// regex the stuff here, use a function i guess
		// put all the expCheck values into a slice and send it back to the front end
		// combine the information from content with the slice that you'll create
		expCheck(j.Name, episode, quality, x)
	}
}

func expCheck(name string, episode int, quality string, x *xdcc) {
	if episode == -1 {
		buildRegex := "(?s)" + "\\" + "[" + "(.*)" + "\\" + "]" + " (.*) - " + "(.*)" + " " + "\\" + "[(" + quality + ")" + "\\" + "]"
		// re := regexp.MustCompile(`(?s)\[(.*)\] (.*) - (.*) \[(.*)\]`)
		re := regexp.MustCompile(buildRegex)
		m := re.FindAllStringSubmatch(name, -1)

		if len(m) > 0 {
			fmt.Printf("Capture value: %s", m[0][0])
			fmt.Println("")
		}
	} else {
		t := strconv.Itoa(episode)
		if episode < 10 {
			t = "0" + t
		}

		buildRegex := "(?s)" + "\\" + "[" + "(.*)" + "\\" + "]" + " (.*) - " + "(" + t + ")" + " " + "\\" + "[(" + quality + ")" + "\\" + "]"
		// re := regexp.MustCompile(`(?s)\[(.*)\] (.*) - (.*) \[(.*)\]`)
		re := regexp.MustCompile(buildRegex)
		m := re.FindAllStringSubmatch(name, -1)

		if len(m) > 0 {
			fmt.Printf("Capture value: %s", m[0][0])
			fmt.Println("")
		}
	}

}

func getJSON(url string, x *xdcc) {
	rs, err := http.Get(url)
	// Process response
	if err != nil {
		panic(err) // More idiomatic way would be to print the error and die unless it's a serious error
	}
	defer rs.Body.Close()

	bodyBytes, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bodyBytes, &x)

	// b, err := json.MarshalIndent(x, "", "  ")
	// if err != nil {
	// 	fmt.Println("error:", err)
	// }
	// // fmt.Println(string(b))
}
