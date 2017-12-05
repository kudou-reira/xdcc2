package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type episodesSuggested struct {
	Compilation   []tempSuggested `json: "compilation"`
	EpisodeNumber int             `json: "groupedEpisode"`
}

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

type errorMessage struct {
	Error        bool   `json: "errorExist"`
	ErrorMessage string `json: "errorMessage"`
}

// type compiledSuggest

func tempSearchMain(querySuggestion string) ([]episodesSuggested, errorMessage) {
	var compiledSuggestions []episodesSuggested
	var err errorMessage

	querySuggestion = strings.TrimSpace(querySuggestion)
	// query this with match query
	// 3-gatsu no lion 4, 10
	// before match query, remove the quality tag [720] for example
	quality, newQuery := getQuality(querySuggestion)
	episodeNumbers, errMsg, arrType := getEpisode(newQuery)
	querySuggestion = matchQuery(newQuery)

	fmt.Println("this is episode number", episodeNumbers)
	fmt.Println("this is the error message", errMsg)
	fmt.Println("this is the quality", quality)

	if errMsg {
		err = errorMessage{
			Error:        errMsg,
			ErrorMessage: "Your query episode numbers are inverted!",
		}
	} else {
		// test case
		querySuggestion = strings.Replace(querySuggestion, " ", "%20", -1)
		fmt.Println("this is the querySuggestion", querySuggestion)
		temp := xdcc{}
		var compilation []episodesSuggested
		var collection []tempSuggested
		// make a slice for suggestions

		compiledSuggestions = findPacklist(querySuggestion, episodeNumbers, quality, arrType, &temp, collection, compilation)
		// pretty print tempSuggestion

		// slcT, _ := json.MarshalIndent(compiledSuggestions, "", " ")
		// fmt.Println(string(slcT))

		// form the query here once you figure out what the user wants
		// ex. gamers! 9 will return did you mean [HorribleSubs] Gamers! - 09[480].mkv?

		err = errorMessage{
			Error:        errMsg,
			ErrorMessage: "",
		}
	}

	return compiledSuggestions, err
}

func findPacklist(query string, episode []int, quality string, arrType string, x *xdcc, collection []tempSuggested, compilation []episodesSuggested) []episodesSuggested {
	var tempEpisodeHold []int
	var wg sync.WaitGroup
	if arrType == "continuous" {
		for i := episode[0]; i <= episode[1]; i++ {
			tempEpisodeHold = append(tempEpisodeHold, i)
		}
		episode = tempEpisodeHold
	}

	fmt.Println("this is in packlist before searching query")
	if len(episode) < 1 {
		unlimitedEps := -1
		queryString := fmt.Sprintf("https://api.nibl.co.uk:8080/nibl/search?query=%s&episodeNumber=%d", query, unlimitedEps)
		fmt.Println("this is the unique query string", queryString)
		getJSON(queryString, x)

		if len(x.Content) > 0 {
			fmt.Println("this is a valid query")
			// slcT, _ := json.MarshalIndent(x.Content, "", " ")
			// fmt.Println("this is the collection of responses")
			// fmt.Println(string(slcT))
			collection = createSuggestion(x, collection, quality)

			tempCollection := episodesSuggested{
				Compilation: collection,
			}

			compilation = append(compilation, tempCollection)
		}

	} else {
		// ranging up to a nonexistent value like 1-29 for blen

		var newCollection []tempSuggested
		for _, singleEP := range episode {

			wg.Add(1)
			go func(singleEP int) {
				queryString := fmt.Sprintf("https://api.nibl.co.uk:8080/nibl/search?query=%s&episodeNumber=%d", query, singleEP)
				fmt.Println("this is the unique query string", queryString)
				fmt.Println("this is the quality", quality)
				getJSON(queryString, x)

				// had to create new array for collection, it was reusing the old collection
				// hence the duplicate values
				if len(x.Content) > 0 {
					fmt.Println("this is a valid query")
					// slcT, _ := json.MarshalIndent(x.Content, "", " ")
					// fmt.Println("this is the collection of responses")
					// fmt.Println(string(slcT))
					newCollection := createSuggestion(x, newCollection, quality)
					tempCollection := episodesSuggested{
						Compilation:   newCollection,
						EpisodeNumber: singleEP,
					}

					// slcT, _ := json.MarshalIndent(tempCollection, "", " ")
					// fmt.Println(string(slcT))

					compilation = append(compilation, tempCollection)
					wg.Done()
				}
			}(singleEP)
		}
		wg.Wait()
		sortByEpisode(compilation)
	}

	return compilation
}

func sortByEpisode(compilation []episodesSuggested) []episodesSuggested {
	sort.Slice(compilation, func(i, j int) bool {
		return compilation[i].EpisodeNumber < compilation[j].EpisodeNumber
	})

	return compilation
}

func createSuggestion(x *xdcc, collection []tempSuggested, quality string) []tempSuggested {
	for _, j := range x.Content {
		// suggest := expCheck(j.Name, episode, quality)
		suggest := j.Name
		qualityExist := checkQuality(suggest, quality)
		if len(suggest) > 0 && qualityExist {
			run := true
			for k := range collection {
				if suggest == collection[k].Suggestion {
					collection[k].SuggestionContent = append(collection[k].SuggestionContent, j)
					// fmt.Println(k)
					// slcT, _ := json.MarshalIndent(collection, "", " ")
					// fmt.Println("this is the current collection")
					// fmt.Println(string(slcT))

					run = false
				}
			}

			if run {
				var tempSuggestionContent []xdccContent
				temp := tempSuggested{
					Suggestion:        suggest,
					SuggestionContent: append(tempSuggestionContent, j),
				}
				collection = append(collection, temp)

				// 	// fmt.Println("this is the length of the collection", len(collection))
				// 	// fmt.Println("we are on index", i)
				// 	// fmt.Println("this is the else if")

			}
		}
	}
	// collection = groupDuplicates(collection)
	// do a final group here in case
	return collection
}

func checkQuality(name string, quality string) bool {
	var isQuality bool

	if len(quality) > 0 {
		check := regexp.MustCompile(quality)
		c := check.FindAllStringSubmatch(name, -1)
		if len(c) > 0 {
			isQuality = true
		}
	} else if len(quality) == 0 {
		isQuality = true
	}

	return isQuality
}

func getQuality(name string) (string, string) {
	var tempQualityArr string
	var tempQuality string
	quality := regexp.MustCompile(`\[(.*?)\]`)
	q := quality.FindAllStringSubmatch(name, -1)
	if len(q) > 0 {
		fmt.Println("this is tempquality", q)
		tempQualityArr = q[0][0]
		tempQuality = q[0][1]
	}
	newName := strings.Replace(name, tempQualityArr, "", -1)

	return tempQuality, newName
}

func getEpisode(name string) ([]int, bool, string) {
	var episodeNumbers []string
	var arrType string

	queryOnly := matchQuery(name)
	fmt.Println("this is query only", queryOnly)
	// name = strings.Replace(name, queryOnly, "", -1)
	// fmt.Println("this is the new name", name)
	// use go routines and channels to get back values
	cont1 := make(chan []string)
	cont2 := make(chan bool)
	single1 := make(chan []string)
	multiple1 := make(chan []string)

	go matchContinuous(name, cont1, cont2)
	go matchSingle(name, single1)
	go matchMultiple(name, multiple1)

	arrCont := <-cont1
	errMsg := <-cont2
	arrSingle := <-single1
	arrMult := <-multiple1

	fmt.Println("this is arrCont", arrCont)
	fmt.Println("there is an error", errMsg)
	fmt.Println("this is arrSingle", arrSingle)
	fmt.Println("this is arr multiple", arrMult)

	if len(arrCont) == len(arrMult) && len(arrCont) > 0 {
		episodeNumbers = arrCont
		arrType = "continuous"
	} else if len(arrMult) > len(arrCont) && len(arrMult) > len(arrSingle) {
		episodeNumbers = arrMult
		arrType = "multiple"
	} else {
		episodeNumbers = arrSingle
		arrType = "single"
	}

	episodeNumbersInt := convertAndSort(episodeNumbers)

	return episodeNumbersInt, errMsg, arrType
}

func convertAndSort(nums []string) []int {
	var tempNumbers []int
	for _, x := range nums {
		y, err := strconv.Atoi(x)
		if err != nil {
			panic(err)
		}
		tempNumbers = append(tempNumbers, y)
	}
	sort.Ints(tempNumbers)
	return tempNumbers
}

func matchContinuous(name string, cont1 chan []string, cont2 chan bool) {
	var tempEpisodes []string
	var errMsg bool
	continuousEpisodes := regexp.MustCompile(" [0-9]+-[0-9]+")
	c := continuousEpisodes.FindAllStringSubmatch(name, -1)
	if len(c) > 0 {
		tempString := strings.Split(strings.TrimSpace(c[0][0]), "-")
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
	cont1 <- tempEpisodes
	cont2 <- errMsg
}

func matchSingle(name string, single1 chan []string) {
	var tempEpisodes []string

	singleEpisode := regexp.MustCompile(" [0-9]+")
	s := singleEpisode.FindAllStringSubmatch(name, -1)
	if len(s) > 0 {
		tempEpisodes = append(tempEpisodes, strings.TrimSpace(s[0][0]))
	}

	single1 <- tempEpisodes
}

func matchMultiple(name string, multiple1 chan []string) {
	var tempEpisodes []string

	// process these values next, send back as an array of strings for episode numbers in getEpisode method
	tempName := strings.Replace(name, matchQuery(name), "", -1)
	multipleEpisodes := regexp.MustCompile(" ?[0-9]+,?")
	m := multipleEpisodes.FindAllStringSubmatch(tempName, -1)

	if len(m) > 0 {
		for _, s := range m {
			dupeExist := false
			tempVal := strings.TrimSpace(strings.Replace(s[0], ",", "", -1))
			if len(tempEpisodes) > 0 {
				for _, u := range tempEpisodes {
					if tempVal == u {
						dupeExist = true
					}
				}
			}

			if dupeExist == false {
				tempEpisodes = append(tempEpisodes, tempVal)
			}
		}
	}

	multiple1 <- tempEpisodes
}

func matchQuery(name string) string {
	// have to be able to find something before a space, number, then comma
	// ex: space9,
	var onlyQuery string
	var cutPoint int

	for i := 0; i < len(name); i++ {
		if name[i] == 32 && 48 <= name[i+1] && name[i+1] <= 57 {
			cutPoint = i
			break
		}
	}
	newName := name[:cutPoint]
	fmt.Println("this is newName", newName)

	if len(newName) < 1 {
		newName = strings.TrimSpace(name)
	}

	onlyQuery = newName

	return onlyQuery
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

// func groupDuplicates(c []tempSuggested) []tempSuggested {
// 	for i := 0; i < len(c); i++ {
// 		for j := i + 1; j < len(c); j++ {
// 			if c[i].Suggestion == c[j].Suggestion {
// 				// remove duplicate
// 				c[i].SuggestionContent = append(c[i].SuggestionContent, c[j].SuggestionContent...)
// 				c = append(c[:j], c[j+1:]...)
// 				j--
// 			}
// 		}
// 	}
// 	return c
// }

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
