package main

import (
	"timekeeper/cmd"
	"timekeeper/tracker"
	"timekeeper/tracker/storage"
)

func main() {
	historyStorage := storage.New()
	timeTracker := tracker.New(historyStorage)
	cmd.Execute(timeTracker)
}
