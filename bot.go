package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf"
)

type bots struct {
	Content []botContent `json: "content"`
}

type botContent struct {
	BotName    string `json: "botName"`
	PackNumber string `json: "packNumber"`
	FileSize   string `json: "fileSize"`
	FileName   string `json: "fileName"`
}

func botSearchMain() {
	// testlink
	// https://nibl.co.uk/bots.php?search=[HorribleSubs] Gamers! - 04 [720p]

	botLink := "https://nibl.co.uk/bots.php?search="
	botQuery := "[HorribleSubs] Gamers! - 04 [720p]"
	combinedBotQuery := botLink + botQuery
	newCombinedQuery := strings.Replace(combinedBotQuery, " ", "%20", -1)
	// testCombine := "https://nibl.co.uk/bots.php?search=gamers"
	fmt.Println("this is replaced strings", newCombinedQuery)

	var botCollection []botContent
	// format the combinedBotQuery
	// replace spaces with %20
	tempBotCollection := accessBotPage(newCombinedQuery, botCollection)

	slcT, _ := json.MarshalIndent(tempBotCollection, "", " ")
	fmt.Println(string(slcT))
}

func accessBotPage(combinedQuery string, collection []botContent) []botContent {
	// var wg sync.WaitGroup
	// wg.Add(1)
	return scrapeBotPage(combinedQuery, collection)
	// wg.Wait()
}

func scrapeBotPage(combinedQuery string, collection []botContent) []botContent {
	fmt.Println("scrapeBotPage is running")
	// doc.find isn't running
	bow := surf.NewBrowser()
	err := bow.Open(combinedQuery)
	if err != nil {
		panic(err)
	}

	bow.Dom().Find(".botlistitem").Each(func(index int, item *goquery.Selection) {
		// use index to determine if jp event only/event character only
		// fmt.Println("botname is", item.Text())
		fmt.Println(index)

		botName := item.Find(".name").Text()
		packNumber := item.Find(".packnumber").Text()
		fileSize := item.Find(".filesize").Text()
		fileName := item.Find(".filename").Text()

		tempBot := botContent{
			BotName:    botName,
			PackNumber: packNumber,
			FileSize:   fileSize,
			FileName:   formatFileName(fileName),
		}

		collection = append(collection, tempBot)

		// fmt.Println("botname is", botName)
		// fmt.Println("packNumber is", packNumber)
		// fmt.Println("filesize is", fileSize)
		// fmt.Println("filename is", fileName)

	})
	return collection
}

func formatFileName(name string) string {
	tempString := strings.Replace(name, "\n", "", -1)
	tempString = strings.Replace(tempString, "[s]", "", -1)
	tempString = strings.TrimSpace(tempString)
	return tempString
}

// func botSearchMain(query string, episode int, quality string, x *xdcc) {
// 	fetchBotList(query, episode, quality, x)
// }

// func fetchBotList(query string, episode int, quality string, x *xdcc) {
// 	if episode == -1 {
// 		// search all strings
// 	}
// }
