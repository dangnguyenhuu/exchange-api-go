package exchange

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type CurrentData struct {
	Prices []Price `json:"prices"`
}

type Price struct {
	Instrument string    `json:"instrument"`
	Time       time.Time `json:"time"`
	Bid        float64   `json:"bid"`
	Ask        float64   `json:"ask"`
	Status     string    `json:"status"`
}

type OANDACurrentData struct {
	url         string
	instruments []string
	layout      string
	since       string
}

func (cd *OANDACurrentData) SetData(instruments []string, layout, since string) {
	cd.url = currentURL
	cd.instruments = instruments
	cd.layout = layout
	cd.since = since
}

func (cd *OANDACurrentData) GetData() (*CurrentData, error) {
	resp, err := cd.GetResponse()
	if err != nil {
		return nil, errors.Wrap(err, "Error1 at CurrentData")
	}
	defer resp.Body.Close()

	var data CurrentData
	err = GetUnmarshal(resp.Body, &data)
	if err != nil {
		return nil, errors.Wrap(err, "Error2 at CurrentData")
	}

	return &data, nil
}

func (cd *OANDACurrentData) GetResponse() (*http.Response, error) {
	values := url.Values{}
	values.Set("accountId", userID)
	values.Add("instruments", strings.Join(cd.instruments, ","))

	if cd.since != "" {
		s, err := time.Parse(cd.layout, cd.since)
		if err != nil {
			return nil, &ParseTimeError{}
		}
		values.Add("since", fmt.Sprint(s.Format(time.RFC3339)))
	}

	req, err := http.NewRequest("GET", cd.url, nil)
	if err != nil {
		return nil, &CreateReqError{}
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.URL.RawQuery = values.Encode()

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, &GetRespError{}
	}

	return resp, nil
}
