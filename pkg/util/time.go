package util

import "time"

func TimeFromMs(ms int64) *time.Time {
	ts := time.Unix(0, ms*int64(time.Millisecond))
	return &ts
}

func TimeFromRFC3339(timestamp string) (*time.Time, error) {
	timestamp = RemoveNonPrintableCharacters(timestamp)
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
