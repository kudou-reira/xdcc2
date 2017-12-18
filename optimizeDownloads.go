package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"sync"
)

type uniqueBot struct {
	FileName    string `json: "fileName"`
	BotToUse    string `json: "botToUse"`
	MessageCall string `json: "callToUse"`
	PackNumber  int    `json: "packNumber"`
}

type simplifiedBot struct {
	BotName       string   `json: "botName"`
	TemporaryBots []string `json: "temporaryBots"`
}

type kv struct {
	Key   string
	Value int
}

func optimizeDLMain(receivedBots []bots) []uniqueBot {
	collectionUnoptimized := findUniqueBots(receivedBots)

	slcT, _ := json.MarshalIndent(collectionUnoptimized, "", " ")
	fmt.Println(string(slcT))

	frequency := findBotFrequency(collectionUnoptimized)
	collectionOptimized := compareBots(collectionUnoptimized, frequency)
	botsOptimized := generateMessageCall(receivedBots, collectionOptimized)

	return botsOptimized
}

func generateMessageCall(receivedBots []bots, optimized []uniqueBot) []uniqueBot {
	for _, j := range receivedBots {
		messageToUse := ""
		for _, l := range j.BotOverall[0].BotSpecies {
			// l is just all the botspecies
			// fmt.Println(a)
			// fmt.Println(l)
			// fmt.Println(l.BotName)
			// fmt.Println("")
			// for _, n := range optimized {
			// 	if l.BotName == n.BotToUse && l.FileName == n.FileName {
			// 		messageToUse = l.MessageCall

			// 		break
			// 	}
			// 	fmt.Println(n)
			// }

			for m, n := range optimized {
				if l.BotName == n.BotToUse && l.FileName == n.FileName {
					fmt.Println("this is inside optimized", optimized[m])
					messageToUse = l.MessageCall
					optimized[m].MessageCall = messageToUse
					fmt.Println("this is what should be used", optimized[m].MessageCall)

					// convert packnumber to int
					intPack, err := strconv.Atoi(l.PackNumber)
					if err != nil {
						// handle error
						panic(err)
					}

					optimized[m].PackNumber = intPack
					break
				}
			}

			if len(messageToUse) > 0 {
				fmt.Println("this was l", l)
				fmt.Println("breaking out of messageToUse")
				break
			}
		}
	}
	// slcT, _ := json.MarshalIndent(optimized, "", " ")
	// fmt.Println(string(slcT))

	return optimized
}

func compareBots(unoptimized []simplifiedBot, frequency map[string]int) []uniqueBot {
	// do it for each item in unoptimized
	var wg sync.WaitGroup
	var optimizedBots []uniqueBot

	empty := createEmptyFrequency(frequency)
	newFrequency := frequency
	slcT, _ := json.MarshalIndent(newFrequency, "", " ")
	fmt.Println(string(slcT))
	fmt.Println("")

	// need to make a map of all available keys, but set their values to 0
	// need to keep track that the values of the specific key do not exceed 2
	// this is because only 2 ports can be opened on xdcc for a specific bot

	for _, b := range unoptimized {
		// should probably make this a goroutine
		// but there are race conditions, so you might not get optimal result
		// because the time differs in which the map gets updated
		// fmt.Println("this is b", b)

		wg.Add(1)

		go func(b simplifiedBot) {
			currentFrequencies := sortMap(newFrequency)

			//range over the lower currentFrequencies and do the check to see if it's available
			for _, d := range currentFrequencies {
				fmt.Println("this is d", d)

				botToUse := ""
				// now range over the temporarybots
				for _, f := range b.TemporaryBots {
					// fmt.Println("this is temporarybots", f)

					if d.Key == f && newFrequency[d.Key] != 0 && empty[d.Key] != 2 {
						tempVal := newFrequency[d.Key]
						tempVal--
						newFrequency[d.Key] = tempVal

						tempEmptyVal := empty[d.Key]
						tempEmptyVal++
						empty[d.Key] = tempEmptyVal

						fmt.Println("this is the file", b.BotName)
						fmt.Println("this bot should use", d.Key)
						botToUse = d.Key

						finalBot := uniqueBot{
							FileName: b.BotName,
							BotToUse: d.Key,
						}

						optimizedBots = append(optimizedBots, finalBot)

						// put this botToUse into a struct
						break
					}
					// stop ranging over temporary bots
				}

				if len(botToUse) > 0 {
					wg.Done()
					break
				}
				// stop ranging over current frequencies
			}
			fmt.Println("")
		}(b)

	}
	wg.Wait()
	fmt.Println("this is the empty map")
	slcT3, _ := json.MarshalIndent(empty, "", " ")
	fmt.Println(string(slcT3))

	fmt.Println("this is the newFrequency map")
	slcT2, _ := json.MarshalIndent(newFrequency, "", " ")
	fmt.Println(string(slcT2))

	return optimizedBots
}

func createEmptyFrequency(frequency map[string]int) map[string]int {
	empty := make(map[string]int)
	for i := range frequency {
		empty[i] = 0
	}

	return empty
}

func sortMap(frequency map[string]int) []kv {
	var ss []kv
	for k, v := range frequency {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value < ss[j].Value
	})

	// for _, kv := range ss {
	// 	fmt.Printf("%s, %d\n", kv.Key, kv.Value)
	// }

	return ss
}

func findBotFrequency(unoptimized []simplifiedBot) map[string]int {
	// go through all the arrays, then add the bot name to a new struct or data structure
	// if the botname already exists, increment ++ the botname
	a := make(map[string]int)
	for i, j := range unoptimized {
		fmt.Println(i)
		fmt.Println(j)
		for _, l := range j.TemporaryBots {
			if val, ok := a[l]; ok {
				//do something here
				if ok {
					a[l] = val + 1
				}
			} else {
				a[l] = 1
			}
		}
	}

	fmt.Println("this is findBotFrequency", a)

	// remove anything with NEW or v6

	newA := removeProblemBots(a)

	fmt.Println("this is newA", newA)

	return a
}

func removeProblemBots(a map[string]int) map[string]int {
	new := regexp.MustCompile("NEW")
	holland := regexp.MustCompile("HOLLAND")
	v6 := regexp.MustCompile("v6")

	var tempRemove []string

	for k := range a {
		// fmt.Printf("key[%s] value[%d]\n", k, v)
		tempNew := new.FindAllString(k, -1)
		tempHolland := holland.FindAllString(k, -1)
		v6 := v6.FindAllString(k, -1)

		// fmt.Println("this is temp", temp)
		if len(tempNew) > 0 || len(tempHolland) > 0 || len(v6) > 0 {
			tempRemove = append(tempRemove, k)
		}
	}

	for _, n := range tempRemove {
		delete(a, n)
	}

	fmt.Println("this is values to remove from map", tempRemove)
	fmt.Println("this is new a", a)

	return a
}

func findUniqueBots(receivedBots []bots) []simplifiedBot {
	var collectionUnoptimized []simplifiedBot
	for i := range receivedBots {
		var tempBotNames []string
		for j := range receivedBots[i].BotOverall[0].BotSpecies {
			tempBotNames = append(tempBotNames, receivedBots[i].BotOverall[0].BotSpecies[j].BotName)
		}
		tempSingleBot := simplifiedBot{
			BotName:       receivedBots[i].BotOverall[0].FileName,
			TemporaryBots: tempBotNames,
		}
		collectionUnoptimized = append(collectionUnoptimized, tempSingleBot)
	}
	return collectionUnoptimized
}
