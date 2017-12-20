package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

type GraphQL struct {
	Query     string `json:"query"`
	Variables string `json:"variables"`
}

type Raw struct {
	Data Data `json:"data"`
}

type Data struct {
	Page Page `json:"Page"`
}

type Media struct {
	ID                int               `json:"id"`
	CoverImage        CoverImage        `json:"coverImage"`
	Title             Title             `json:"title"`
	Studios           Studios           `json:"studios"`
	Description       string            `json:"description"`
	Type              string            `json:"type"`
	Format            string            `json:"format"`
	Episodes          int               `json:"episodes"`
	Genres            []string          `json:"genres"`
	AverageScore      int               `json:"averageScore"`
	Popularity        int               `json:"popularity"`
	StartDate         StartDate         `json:"startDate"`
	EndDate           EndDate           `json:"endDate"`
	Season            string            `json:"season"`
	NextAiringEpisode NextAiringEpisode `json:"nextAiringEpisode"`
}

type Page struct {
	PageInfo PageInfo `json:"pageInfo"`
	Media    []Media  `json:"media"`
}

type Title struct {
	UserPreferred string `json:"userPreferred"`
}

type Studios struct {
	Nodes []Nodes `json:"nodes"`
}

type Nodes struct {
	Name string `json:"name"`
}

type CoverImage struct {
	Large  string `json:"large"`
	Medium string `json:"medium"`
}

type StartDate struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"dat"`
}

type EndDate struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"dat"`
}

type NextAiringEpisode struct {
	AiringAt        int `json:"airingAt"`
	TimeUntilAiring int `json:"timeUntilAiring"`
	Episode         int `json:"episode"`
}

type PageInfo struct {
	Total       int  `json:"total"`
	PerPage     int  `json:"perPage"`
	CurrentPage int  `json:"currentPage"`
	LastPage    int  `json:"lastPage"`
	HasNextPage bool `json:"hasNextPage"`
}

func anilistMain() []Media {
	fmt.Println("this is anilistMain")

	queryGraph := `query (
		$page: Int,
		$type: MediaType,
		$format: MediaFormat,
		$startDate: String,
		$endDate: String,
		$season: MediaSeason,
		$genres: [String],
		$genresExclude: [String],
		$isAdult: Boolean = false, # Assign default value if isAdult is not included in our query variables 
		$sort: [MediaSort],
	  ) {
		Page (page: $page) {
		  pageInfo {
			total
			perPage
			currentPage
			lastPage
			hasNextPage
		  }
		  media (
			startDate_like: $startDate, # "2017%" will get all media starting in 2017, alternatively you could use the lesser & greater suffixes
			endDate_like: $endDate,
			season: $season,
			type: $type,
			format: $format,
			genre_in: $genres,
			genre_not_in: $genresExclude,
			isAdult: $isAdult,
			sort: $sort,
		  ) {
			id
			title {
			  userPreferred
			}
			coverImage {
			  large
			  medium
			}			
			studios {
				nodes{
					name
		  		}
			}
			description
			type
			format
			episodes
			chapters
			volumes
			genres
			averageScore
			popularity
			startDate {
			  year
			  month
			  day
			}
			endDate {
			  year
			  month
			  day
			}
			season
			nextAiringEpisode {
			  airingAt
			  timeUntilAiring
			  episode
			}
		  }
		}
	  }`

	variablesGraph := `{
		"startDate": "2017%",
		"season": "FALL"
	}`

	var tempMedia []Media

	tempMedia, totalPages := fetchInitialQuery(queryGraph, variablesGraph, tempMedia)
	fmt.Println("total pages", totalPages)

	tempMedia = fetchSubsequent(queryGraph, tempMedia, totalPages)

	return tempMedia
	// tempMedia now has all the media queries []Media
	// next, sort by until airing date
	//

	// slcT, _ := json.MarshalIndent(tempMedia, "", " ")
	// fmt.Println("response body", string(slcT))
}

func fetchInitialQuery(queryGraph string, variablesGraph string, tempMedia []Media) ([]Media, int) {
	url := "https://graphql.anilist.co"
	client := http.Client{}

	tempGraphQL := GraphQL{
		Query:     queryGraph,
		Variables: variablesGraph,
	}

	graphBytes, err := json.Marshal(tempGraphQL)
	if err != nil {
		panic(err)
	}

	fmt.Println("this is tempGraphQL", tempGraphQL)
	fmt.Println("this is graphBytes", graphBytes)
	// fmt.Println("this is string graphBytes", string(graphBytes))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(graphBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Unable to reach the server.")
	}

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

	// fmt.Println("response Body:", string(body))

	raw := Raw{}
	err = json.Unmarshal(body, &raw)
	if err != nil {
		panic(err)
	}

	// slcT, _ := json.MarshalIndent(raw, "", " ")
	// fmt.Println("response body", string(slcT))

	tempMedia = raw.Data.Page.Media
	totalPages := raw.Data.Page.PageInfo.LastPage

	// slcT, _ := json.MarshalIndent(tempMedia, "", " ")
	// fmt.Println("response body", string(slcT))

	return tempMedia, totalPages
}

func fetchSubsequent(queryGraph string, tempMedia []Media, totalPages int) []Media {
	// put all of this in for loop

	// inject template string
	// variablesGraph := `{
	// 	"startDate": "2017%",
	// 	"season": "FALL"
	// }`

	var wg sync.WaitGroup
	wg.Add(totalPages - 1)

	for i := 2; i <= totalPages; i++ {
		go func(i int) {
			defer wg.Done()
			bracket1 := "{"
			bracket2 := "}"
			startDate := `"startDate"`
			season := `"season"`
			currentPage := `"page"`
			colon := ":"
			singleQuote := `"`
			comma := ","

			startValue := "2017%"
			seasonValue := "FALL"
			currentPageValue := strconv.Itoa(i)

			variablesGraphSeason := bracket1 + "\n" +
				startDate + colon + singleQuote + startValue + singleQuote + comma + "\n" +
				season + colon + singleQuote + seasonValue + singleQuote + comma + "\n" +
				currentPage + colon + singleQuote + currentPageValue + singleQuote + "\n" +
				bracket2

			fmt.Println("this is variablesgraph", variablesGraphSeason)

			url := "https://graphql.anilist.co"
			client := http.Client{}

			tempGraphQL := GraphQL{
				Query:     queryGraph,
				Variables: variablesGraphSeason,
			}

			graphBytes, err := json.Marshal(tempGraphQL)
			if err != nil {
				panic(err)
			}

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(graphBytes))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Unable to reach the server.")
			}

			defer resp.Body.Close()

			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			body, _ := ioutil.ReadAll(resp.Body)

			fmt.Println("response Body:", string(body))

			raw := Raw{}
			err = json.Unmarshal(body, &raw)
			if err != nil {
				panic(err)
			}

			newMedia := raw.Data.Page.Media

			// slcT, _ := json.MarshalIndent(newMedia, "", " ")
			// fmt.Println("response body", string(slcT))

			tempMedia = append(tempMedia, newMedia...)
		}(i)
	}

	wg.Wait()

	return tempMedia
}
