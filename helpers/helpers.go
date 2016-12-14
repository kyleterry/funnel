package helpers

import "net/http"

var FunnelUserAgent = "linux:FunnelReader:v1.0 (by Kyle Terry <funnel@kyleterry.com>)"

func Fetch(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", FunnelUserAgent)

	return client.Do(req)
}
