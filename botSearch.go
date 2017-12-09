package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type bots struct {
	BotOverall []groupBots `json: "content"`
}

type groupBots struct {
	FileName    string       `json: "fileName"`
	FileQuality string       `json: "fileQuality"`
	BotSpecies  []botContent `json: "botSpecies"`
}

type botContent struct {
	BotName     string `json: "botName"`
	PackNumber  string `json: "packNumber"`
	FileSize    string `json: "fileSize"`
	FileName    string `json: "fileName"`
	MessageCall string `json: "messageCall"`
}

func botSearchMain(stack []string) []bots {
	// testlink
	// https://nibl.co.uk/bots.php?search=[HorribleSubs] Gamers! - 04 [720p]

	fmt.Println("this is the stack in botsearchmain", stack)

	tempBotLinks := createBotLinks(stack)
	fmt.Println("this is tempBotLinks", tempBotLinks)

	// format the combinedBotQuery
	// replace spaces with %20

	// newCombinedQuery := "https://nibl.co.uk/bots.php?search=gamers"
	temp := startBotSearch(tempBotLinks)

	slcT, _ := json.MarshalIndent(temp, "", " ")
	fmt.Println(string(slcT))

	return temp
}

func startBotSearch(tempBotLinks []string) []bots {
	var wg sync.WaitGroup
	var allBots []bots
	for _, singleQuery := range tempBotLinks {
		wg.Add(1)
		var botCollection []groupBots
		go func(singleQuery string) {
			tempBotCollection := accessBotPage(singleQuery, botCollection)
			tempBot := bots{
				BotOverall: tempBotCollection,
			}
			allBots = append(allBots, tempBot)
			wg.Done()
		}(singleQuery)
	}
	wg.Wait()

	return allBots
}

func createBotLinks(stack []string) []string {
	var botLinks []string
	for _, j := range stack {
		baseLink := "https://nibl.co.uk/bots.php?search="
		// botQuery := "[HorribleSubs] Gamers! - 04 [720p]"
		// botQuery := "gamers 04"
		combinedBotQuery := baseLink + j
		newCombinedQuery := strings.Replace(combinedBotQuery, " ", "%20", -1)
		// testCombine := "https://nibl.co.uk/bots.php?search=gamers"
		fmt.Println("this is replaced strings", newCombinedQuery)
		botLinks = append(botLinks, newCombinedQuery)
	}
	return botLinks
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

	doc, err := goquery.NewDocument(combinedQuery)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".botlistitem").Each(func(index int, item *goquery.Selection) {
		// use index to determine if jp event only/event character only
		// fmt.Println("botname is", item.Text())

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

		// fmt.Println("this is the filename", fileName)
		// fmt.Println("this is the botName", botName)

		// add quality to tempGroupBot
		// fmt.Println("this is the tempBot before going", tempBot)
		// fmt.Println("this is the collection before going", collection)

		if index == 0 {
			var tempCollection []botContent
			tempGroupBot := groupBots{
				FileName:    tempBot.FileName,
				FileQuality: extractQuality(tempBot.FileName),
				BotSpecies:  append(tempCollection, tempBot),
			}
			collection = append(collection, tempGroupBot)
		} else {
			inCollection := false
			for i := range collection {
				// if same file name, FINE
				// kept appending on new last one because it's still going over infinite collection
				if collection[i].FileName == tempBot.FileName {
					collection[i].BotSpecies = append(collection[i].BotSpecies, tempBot)
					inCollection = true
				}
			}

			if inCollection == false {
				var tempCollection []botContent
				tempGroupBot := groupBots{
					FileName:    tempBot.FileName,
					FileQuality: extractQuality(tempBot.FileName),
					BotSpecies:  append(tempCollection, tempBot),
				}
				collection = append(collection, tempGroupBot)
			}
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

func extractQuality(name string) string {
	var quality string
	re := regexp.MustCompile(`\[(.*?)\]`)
	m := re.FindAllStringSubmatch(name, -1)
	// fmt.Println("this is m", m)
	if len(m) > 1 {
		fmt.Println("this is extract quality", m)
		quality = m[1][1]
	} else {
		quality = "Not a media file"
	}
	return quality
}
