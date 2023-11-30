package util

import "time"

func ParseTimestamp(timestamp string) (*time.Time, error) {
	timestamp = RemoveNonPrintableCharacters(timestamp)
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
