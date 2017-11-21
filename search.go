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

type suggestedQuery struct {
	Suggestion        string      `json: "suggestion"`
	SuggestionContent xdccContent `json: "suggestionContent"`
}

func searchMain() {
	query := "gamers"
	episodeNumber := 4
	quality := "720p"
	temp := xdcc{}
	var collection []suggestedQuery
	// make a slice for suggestions
	tempSuggestion := findPacklist(query, episodeNumber, quality, &temp, collection)
	// pretty print tempSuggestion
	slcT, _ := json.MarshalIndent(tempSuggestion, "", " ")
	fmt.Println(string(slcT))

	// form the query here once you figure out what the user wants
	// ex. gamers! 9 will return did you mean [HorribleSubs] Gamers! - 09[480].mkv?
	botMain(query, episodeNumber, quality, &temp)
}

func findPacklist(query string, episode int, quality string, x *xdcc, collection []suggestedQuery) []suggestedQuery {
	queryString := fmt.Sprintf("https://api.nibl.co.uk:8080/nibl/search?query=%s&episodeNumber=%d", query, episode)
	getJSON(queryString, x)

	if len(x.Content) > 0 {
		fmt.Println("this is a valid query")
		collection = createSuggestion(episode, quality, x, collection)
	} else {
		fmt.Println("this is not a valid query")
		// return this is not a valid query message to the front end
	}

	return collection
}

func createSuggestion(episode int, quality string, x *xdcc, collection []suggestedQuery) []suggestedQuery {
	for _, j := range x.Content {
		// fmt.Println(j.Name)
		// regex the stuff here, use a function i guess
		// put all the expCheck values into a slice and send it back to the front end
		// combine the information from content with the slice that you'll create
		suggest := expCheck(j.Name, episode, quality)
		if len(suggest) > 0 {
			temp := suggestedQuery{
				Suggestion:        suggest,
				SuggestionContent: j,
			}
			collection = append(collection, temp)
		}
	}
	return collection
}

func expCheck(name string, episode int, quality string) string {
	queryString := ""
	if episode == -1 {
		buildRegex := "(?s)" + "\\" + "[" + "(.*)" + "\\" + "]" + " (.*) - " + "(.*)" + " " + "\\" + "[(" + quality + ")" + "\\" + "]"
		// re := regexp.MustCompile(`(?s)\[(.*)\] (.*) - (.*) \[(.*)\]`)
		re := regexp.MustCompile(buildRegex)
		m := re.FindAllStringSubmatch(name, -1)

		if len(m) > 0 {
			// fmt.Printf("Capture value: %s", m[0][0])
			// fmt.Println("")
			queryString = m[0][0]
		}

	} else {
		t := strconv.Itoa(episode)
		if episode < 10 {
			t = "0" + t
		}

		buildRegex := "(?s)" + "\\" + "[" + "(.*)" + "\\" + "]" + " (.*) - " + "(" + t + ")" + " " + "\\" + "[(" + quality + ")" + "\\" + "]"
		re := regexp.MustCompile(buildRegex)
		m := re.FindAllStringSubmatch(name, -1)

		if len(m) > 0 {
			// fmt.Printf("Capture value: %s", m[0][0])
			// fmt.Println("")
			queryString = m[0][0]
		}
	}
	return queryString
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
