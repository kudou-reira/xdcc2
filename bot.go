package main

func botMain(query string, episode int, quality string, x *xdcc) {
	fetchBotList(query, episode, quality, x)
}

func fetchBotList(query string, episode int, quality string, x *xdcc) {
	if episode == -1 {
		// search all strings
	}
}
