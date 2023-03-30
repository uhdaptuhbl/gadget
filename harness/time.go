package harness

import (
	"encoding/json"
	"time"
)

const SecondsInDay = time.Second * 60 * 60 * 24

func ParseRFC3339(timestr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timestr)
}

func RFC3339(t time.Time) string {
	return t.Format(time.RFC3339)
}
func Date(t time.Time) string {
	return t.Format("2006-01-02")
}
func UTCnow() time.Time {
	var utc = time.Now().UTC()

	return utc
}

func UTCnowRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func DayStart(t time.Time) time.Time {
	return t.Truncate(SecondsInDay)
}

func DayEnd(t time.Time) time.Time {
	return t.Truncate(SecondsInDay).Add(SecondsInDay - (time.Second * 1))
}

func DayStartRFC3339(t time.Time) string {
	return t.Truncate(SecondsInDay).Format(time.RFC3339)
}

func DayEndRFC3339(t time.Time) string {
	return t.Truncate(SecondsInDay).Add(SecondsInDay - (time.Second * 1)).Format(time.RFC3339)
}

// MarshalTimePtr casts a time.Time value to *MarshalTime
func MarshalTimePtr(t time.Time) *MarshalTime {
	var val = MarshalTime(t)
	return &val
}

type MarshalTime time.Time

func (t MarshalTime) String() string {
	return time.Time(t).Truncate(time.Second * 1).Format(time.RFC3339)
}
func (t MarshalTime) MarshalJSON() ([]byte, error) {
	var err error
	var data []byte
	if data, err = json.Marshal(t.String()); err != nil {
		return nil, err
	}
	return data, nil
}
func (t *MarshalTime) UnmarshalJSON(data []byte) error {
	var err error
	var raw string
	var tnew time.Time
	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if tnew, err = time.Parse(time.RFC3339, raw); err != nil {
		return err
	}
	*t = MarshalTime(tnew)
	return nil
}
func (t MarshalTime) Time() time.Time {
	return time.Time(t)
}
