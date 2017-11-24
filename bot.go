package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type bots struct {
	Content []groupBots `json: "content"`
}

type groupBots struct {
	FileName   string       `json: "fileName"`
	BotSpecies []botContent `json: "botSpecies"`
}

type botContent struct {
	BotName     string `json: "botName"`
	PackNumber  string `json: "packNumber"`
	FileSize    string `json: "fileSize"`
	FileName    string `json: "fileName"`
	MessageCall string `json: "messageCall"`
}

func botSearchMain() {
	// testlink
	// https://nibl.co.uk/bots.php?search=[HorribleSubs] Gamers! - 04 [720p]

	botLink := "https://nibl.co.uk/bots.php?search="
	// botQuery := "[HorribleSubs] Gamers! - 04 [720p]"
	// botQuery := "gamers 04"
	botQuery := "net-juu 08"
	combinedBotQuery := botLink + botQuery
	newCombinedQuery := strings.Replace(combinedBotQuery, " ", "%20", -1)
	// testCombine := "https://nibl.co.uk/bots.php?search=gamers"
	fmt.Println("this is replaced strings", newCombinedQuery)

	var botCollection []groupBots
	// format the combinedBotQuery
	// replace spaces with %20
	tempBotCollection := accessBotPage(newCombinedQuery, botCollection)

	slcT, _ := json.MarshalIndent(tempBotCollection, "", " ")
	fmt.Println(string(slcT))
}

func accessBotPage(combinedQuery string, collection []groupBots) []groupBots {
	// var wg sync.WaitGroup
	waitBot := make(chan []groupBots)

	// wg.Add(1)
	go scrapeBotPage(combinedQuery, collection, waitBot)
	// wg.Wait()
	result := <-waitBot
	return result
}

func scrapeBotPage(combinedQuery string, collection []groupBots, waitBot chan []groupBots) {
	fmt.Println("scrapeBotPage is running")

	// doc.find isn't running
	// bow := surf.NewBrowser()
	// err := bow.Open(combinedQuery)
	// if err != nil {
	// 	panic(err)
	// }

	doc, err := goquery.NewDocument(combinedQuery)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".botlistitem").Each(func(index int, item *goquery.Selection) {
		// use index to determine if jp event only/event character only
		// fmt.Println("botname is", item.Text())
		fmt.Println(index)

		botName := item.Find(".name").Text()
		packNumber := item.Find(".packnumber").Text()
		fileSize := item.Find(".filesize").Text()
		fileName := item.Find(".filename").Text()

		tempBot := botContent{
			BotName:     botName,
			PackNumber:  packNumber,
			FileSize:    fileSize,
			FileName:    formatFileName(fileName),
			MessageCall: createMessage(botName, packNumber),
		}

		fmt.Println("this is the filename", fileName)
		fmt.Println("this is the botName", botName)

		// add quality to tempGroupBot

		if index == 0 {
			var tempCollection []botContent
			tempGroupBot := groupBots{
				FileName:   tempBot.FileName,
				BotSpecies: append(tempCollection, tempBot),
			}
			collection = append(collection, tempGroupBot)
		} else {
			inCollection := false
			for i := range collection {
				// if same file name, FINE
				// keeps appending on new last one because it's still going over infinite collection
				if collection[i].FileName == tempBot.FileName {
					collection[i].BotSpecies = append(collection[i].BotSpecies, tempBot)
					inCollection = true
				}
			}

			if inCollection == false {
				var tempCollection []botContent
				tempGroupBot := groupBots{
					FileName:   tempBot.FileName,
					BotSpecies: append(tempCollection, tempBot),
				}
				collection = append(collection, tempGroupBot)
			}

			fmt.Println("this is the else statement")
		}
	})
	waitBot <- collection
}

func formatFileName(name string) string {
	tempString := strings.Replace(name, "\n", "", -1)
	tempString = strings.Replace(tempString, "[s]", "", -1)
	tempString = strings.TrimSpace(tempString)
	return tempString
}

func createMessage(bot, pack string) string {
	// /msg KareRaisu xdcc send #9924
	return fmt.Sprintf("/msg %s xdcc send #%s", bot, pack)
}
