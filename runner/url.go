package harness

import (
	"encoding/json"
	"net/url"
)

type MarshalURL url.URL

func (u MarshalURL) String() string {
	var uc = url.URL(u)
	return (&uc).String()
}
func (u MarshalURL) MarshalJSON() ([]byte, error) {
	var err error
	var data []byte
	if data, err = json.Marshal(u.String()); err != nil {
		return nil, err
	}
	return data, nil
}
func (u *MarshalURL) UnmarshalJSON(data []byte) error {
	var err error
	var raw string
	var unew *url.URL
	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if unew, err = url.Parse(raw); err != nil {
		return err
	}
	*u = MarshalURL(*unew)
	return nil
}
