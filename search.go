package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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

type tempSuggested struct {
	Suggestion        string        `json: "suggestion"`
	SuggestionContent []xdccContent `json: "suggestionContent"`
}

// type compiledSuggest

func tempSearchMain(querySuggestion string) []tempSuggested {
	// query := "gamers"
	// need to make a function to check if there is a number or episode number in the query
	// then send it off
	episodeNumber := getEpisode(querySuggestion)
	fmt.Println("this is episode number", episodeNumber)
	fakeEpisodeNumber := 4
	// test case
	querySuggestion = strings.Replace(querySuggestion, " ", "%20", -1)
	fmt.Println("this is the querySuggestion", querySuggestion)
	quality := "720p"
	temp := xdcc{}
	var collection []tempSuggested
	// make a slice for suggestions

	// is it not working because of async?
	tempSuggestion := findPacklist(querySuggestion, fakeEpisodeNumber, quality, &temp, collection)
	// pretty print tempSuggestion

	// slcT, _ := json.MarshalIndent(tempSuggestion, "", " ")
	// fmt.Println(string(slcT))

	// form the query here once you figure out what the user wants
	// ex. gamers! 9 will return did you mean [HorribleSubs] Gamers! - 09[480].mkv?
	return tempSuggestion
}

func findPacklist(query string, episode int, quality string, x *xdcc, collection []tempSuggested) []tempSuggested {
	queryString := fmt.Sprintf("https://api.nibl.co.uk:8080/nibl/search?query=%s&episodeNumber=%d", query, episode)
	getJSON(queryString, x)

	if len(x.Content) > 0 {
		fmt.Println("this is a valid query")
		// slcT, _ := json.MarshalIndent(x.Content, "", " ")
		// fmt.Println("this is the collection of responses")
		// fmt.Println(string(slcT))
		collection = createSuggestion(episode, quality, x, collection)
	} else {
		fmt.Println("this is not a valid query")
		// return this is not a valid query message to the front end
	}

	return collection
}

func createSuggestion(episode int, quality string, x *xdcc, collection []tempSuggested) []tempSuggested {
	for i, j := range x.Content {
		// suggest := expCheck(j.Name, episode, quality)
		suggest := j.Name

		// check the loop
		// fmt.Println("hi", i)
		// creates a new collection of suggested names from x.content
		if len(suggest) > 0 {
			if len(collection) == 0 {

				var tempSuggestionContent []xdccContent
				temp := tempSuggested{
					Suggestion:        suggest,
					SuggestionContent: append(tempSuggestionContent, j),
				}
				collection = append(collection, temp)

				fmt.Println("this is the length of the collection", len(collection))
				fmt.Println("we are on index", i)
				fmt.Println("this is the initial if")
				// this gets called first so there's already an existing container
				// but the else statement down below doesn't have an existing container YET
			} else if len(collection) > 0 {
				// initialize a counter for unique

				run := false
				for k := range collection {
					if suggest == collection[k].Suggestion {
						collection[k].SuggestionContent = append(collection[k].SuggestionContent, j)
						// fmt.Println(k)
						// slcT, _ := json.MarshalIndent(collection, "", " ")
						// fmt.Println("this is the current collection")
						// fmt.Println(string(slcT))
					} else {
						run = true
					}
				}

				if run {
					var tempSuggestionContent []xdccContent
					temp := tempSuggested{
						Suggestion:        suggest,
						SuggestionContent: append(tempSuggestionContent, j),
					}
					collection = append(collection, temp)

					// fmt.Println("this is the length of the collection", len(collection))
					// fmt.Println("we are on index", i)
					// fmt.Println("this is the else if")
				}
			}
		}
	}
	collection = groupDuplicates(collection)
	// do a final group here in case
	return collection
}

func groupDuplicates(c []tempSuggested) []tempSuggested {
	for i := 0; i < len(c); i++ {
		for j := i + 1; j < len(c); j++ {
			if c[i].Suggestion == c[j].Suggestion {
				// remove duplicate
				c[i].SuggestionContent = append(c[i].SuggestionContent, c[j].SuggestionContent...)
				c = append(c[:j], c[j+1:]...)
				j--
			}
		}
	}
	return c
}

func getEpisode(name string) string {
	var episodeNumbers string

	// use go routines and channels to get back stuff
	arrCont, errMsg := matchContinuous(name)
	fmt.Println("this is arrCont", arrCont)
	fmt.Println("there is an error", errMsg)

	if len(arrCont) < 1 {
		singleEpisode := regexp.MustCompile(" [0-9]+")
		s := singleEpisode.FindAllStringSubmatch(name, -1)
		if len(s) > 0 {
			fmt.Println("this is m in episode", s)
			episodeNumbers = s[0][0]
		}

		fmt.Println("this is the episode number", episodeNumbers)
	}

	return episodeNumbers
}

func matchContinuous(name string) ([]string, bool) {
	var tempEpisodes []string
	var errMsg bool
	multipleEpisodes := regexp.MustCompile(" [0-9]+-[0-9]+")
	m := multipleEpisodes.FindAllStringSubmatch(name, -1)
	if len(m) > 0 {
		tempString := strings.Split(strings.TrimSpace(m[0][0]), "-")
		tempEpisodes = append(tempEpisodes, tempString[0])
		tempEpisodes = append(tempEpisodes, tempString[1])

		x, err := strconv.Atoi(tempString[0])
		if err != nil {
			panic(err)
		}
		y, err := strconv.Atoi(tempString[1])
		if err != nil {
			panic(err)
		}

		if y < x || y == x {
			errMsg = true
		}
	}
	return tempEpisodes, errMsg
}

func matchSingle(name string) {

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

// func expCheck(name string, episode int, quality string) string {
// 	queryString := ""
// 	if episode == -1 {
// 		buildRegex := "(?s)" + "\\" + "[" + "(.*)" + "\\" + "]" + " (.*) - " + "(.*)" + " " + "\\" + "[(" + quality + ")" + "\\" + "]"
// 		// re := regexp.MustCompile(`(?s)\[(.*)\] (.*) - (.*) \[(.*)\]`)
// 		re := regexp.MustCompile(buildRegex)
// 		m := re.FindAllStringSubmatch(name, -1)

// 		// also do a regex to find if the numbers 480, 720, or 1080 are in the string, or all
// 		// there are some edge cases in which they aren't in brackets
// 		// [DameDesuYo] Blend S - 04 (1920x1080 10bit AAC) [7CA7EB0F].mkv

// 		if len(m) > 0 {
// 			// fmt.Printf("Capture key: %s", m[0][0])
// 			// fmt.Println("")
// 			queryString = m[0][0]
// 		}

// 	} else {
// 		t := strconv.Itoa(episode)
// 		if episode < 10 {
// 			t = "0" + t
// 		}

// 		buildRegex := "(?s)" + "\\" + "[" + "(.*)" + "\\" + "]" + " (.*) - " + "(" + t + ")" + " " + "\\" + "[(" + quality + ")" + "\\" + "]"
// 		re := regexp.MustCompile(buildRegex)
// 		m := re.FindAllStringSubmatch(name, -1)

// 		if len(m) > 0 {
// 			// fmt.Printf("Capture key: %s", m[0][0])
// 			// fmt.Println("")
// 			queryString = m[0][0]
// 		}
// 	}
// 	return queryString
// }
