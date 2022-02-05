package storage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"time"
	"timekeeper/tracker"
)

type TimeHistory []tracker.Item

type history struct {
	historyFilePath    string
	currentKeyFilePath string
}

func (h history) openHistoryFile() (TimeHistory, error) {
	jsonBytes, err := ioutil.ReadFile(h.historyFilePath)
	if err != nil {
		return nil, err
	}

	var timeHistory TimeHistory
	err = json.Unmarshal(jsonBytes, &timeHistory)
	if err != nil {
		return nil, err
	}
	return timeHistory, nil
}

func (h history) StopLastHistoryItem() error {
	timeHistory, err := h.openHistoryFile()
	if err != nil {
		return err
	}
	currentTimestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	if len(timeHistory) > 0 {
		timeHistory[len(timeHistory)-1].Stopped = currentTimestamp
	}

	newHistoryBytes, err := json.Marshal(timeHistory)

	err = ioutil.WriteFile(h.historyFilePath, newHistoryBytes, 0666)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(h.currentKeyFilePath, []byte(""), 0666)
	if err != nil {
		return err
	}

	return nil
}

func (h history) AddHistoryItem(key string) error {
	timeHistory, err := h.openHistoryFile()
	if err != nil {
		return err
	}
	currentTimestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	if len(timeHistory) > 0 {
		timeHistory[len(timeHistory)-1].Stopped = currentTimestamp
	}

	timeHistory = append(timeHistory, tracker.Item{
		Key:     key,
		Started: currentTimestamp,
	})

	newHistoryBytes, err := json.Marshal(timeHistory)
	err = ioutil.WriteFile(h.historyFilePath, newHistoryBytes, 0666)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(h.currentKeyFilePath, []byte(key), 0666)
	if err != nil {
		return err
	}

	return nil
}

func (h history) GetHistoryItems() ([]tracker.Item, error) {
	timeHistory, err := h.openHistoryFile()
	if err != nil {
		return nil, err
	}
	return timeHistory, nil
}

func New() history {
	homeDirName, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	historyFilePath := homeDirName + "/.timekeeper.json"
	currentKeyFilePath := homeDirName + "/.timekeeper_current"

	if _, err := os.Stat(historyFilePath); errors.Is(err, os.ErrNotExist) {
		histFile, err := os.Create(historyFilePath)
		if err != nil {
			panic(err)
		}
		_, err = histFile.WriteString("[]")
		if err != nil {
			panic(err)
		}
		err = histFile.Close()
		if err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat(currentKeyFilePath); errors.Is(err, os.ErrNotExist) {
		emptyFile, err := os.Create(currentKeyFilePath)
		if err != nil {
			panic(err)
		}
		err = emptyFile.Close()
		if err != nil {
			panic(err)
		}
	}

	return history{
		historyFilePath:    historyFilePath,
		currentKeyFilePath: currentKeyFilePath,
	}
}
