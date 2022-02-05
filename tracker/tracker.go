package tracker

import (
	"encoding/json"
	"fmt"
)

type HistoryStorage interface {
	GetHistoryItems() ([]Item, error)
	AddHistoryItem(key string) error
	StopLastHistoryItem() error
}

type timeTracker struct {
	HistoryStorage
}

func New(storage HistoryStorage) timeTracker {
	return timeTracker{
		storage,
	}
}

func (t timeTracker) StopTracking() (string, error) {
	err := t.StopLastHistoryItem()
	if err != nil {
		return "", err
	}
	return "stopped", nil
}

func (t timeTracker) StartTracking(key string) (string, error) {
	err := t.AddHistoryItem(key)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	return "start tracking " + key, nil
}

func (t timeTracker) Export(format string) (string, error) {
	switch format {
	case "json":
		items, err := t.GetHistoryItems()
		if err != nil {
			return "", err
		}
		jsonBytes, err := json.Marshal(items)
		if err != nil {
			return "", err
		}
		return string(jsonBytes), nil
	case "csv":
		return ",,,", nil
	}
	return "", fmt.Errorf("unknown format %s", format)
}
