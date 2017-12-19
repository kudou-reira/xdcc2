package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	ID    int    `json:"id"`
	Title Title  `json:"title"`
	Type  string `json:"type"`
}

type Title struct {
	UserPreferred string `json:"userPreferred"`
}

type Page struct {
	PageInfo PageInfo `json:"pageInfo"`
	Media    []Media  `json:"media"`
}

type PageInfo struct {
	Total       int  `json:"total"`
	PerPage     int  `json:"perPage"`
	CurrentPage int  `json:"currentPage"`
	LastPage    int  `json:"lastPage"`
	HasNextPage bool `json:"hasNextPage"`
}

func anilistMain() {
	fmt.Println("this is anilistMain")

	url := "https://graphql.anilist.co"

	client := http.Client{}

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
		"startDate": "2017%"
	}`

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

	fmt.Println("response Body:", string(body))

	raw := Raw{}
	err = json.Unmarshal(body, &raw)
	if err != nil {
		panic(err)
	}

	fmt.Println("this is data", raw)

	// slcT, _ := json.MarshalIndent(string(body), "", " ")
	// fmt.Println("response body", string(slcT))
}

// func getPages(body []byte) (*page, error) {
// 	var p = new(page)
// 	err := json.Unmarshal(body, &p)
// 	if err != nil {
// 		fmt.Println("whoops:", err)
// 	}
// 	return p, err
// }
