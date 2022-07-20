package directusapi

import (
	"fmt"
	"time"
)

const datetimeFormat = "2006-01-02 15:04:05"

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", t.Format(datetimeFormat))
	return []byte(stamp), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	parsedT, err := time.Parse(fmt.Sprintf("\"%s\"", datetimeFormat), string(data))
	if err != nil {
		return err
	}
	t.Time = parsedT
	return nil
}
