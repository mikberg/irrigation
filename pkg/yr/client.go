package yr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	nowcastURL = "https://api.met.no/weatherapi/nowcast/2.0/complete"
	userAgent  = "irrigationGo github.com/mikbe/irrigation"
)

type Client interface {
	Nowcast(lat, lon, alt float64) (METJSONForecast, error)
}

type client struct {
	httpClient *http.Client
}

func NewClient() Client {
	httpClient := &http.Client{
		Transport: &RoundTripper{
			http.DefaultTransport,
		},
	}

	return &client{
		httpClient: httpClient,
	}
}

func (c *client) Nowcast(lat, lon, alt float64) (METJSONForecast, error) {
	forecast := METJSONForecast{}

	// u, err := url.Parse(nowcastURL)
	// if err != nil {
	// 	return forecast, fmt.Errorf("failed to parse URL: %w", err)
	// }

	values := url.Values{}
	values.Add("lat", fmt.Sprintf("%f", lat))
	values.Add("lon", fmt.Sprintf("%f", lon))
	if alt != 0.0 {
		values.Add("alt", fmt.Sprintf("%f", alt))
	}

	u := fmt.Sprintf("%s?%s", nowcastURL, values.Encode())
	ret, err := c.httpClient.Get(u)
	if err != nil {
		return forecast, fmt.Errorf("failed to get forecast: %w", err)
	}
	if ret.StatusCode != 200 {
		return forecast, fmt.Errorf("unexpected status code: %d", ret.StatusCode)
	}

	defer ret.Body.Close()
	data, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		return forecast, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(data, &forecast); err != nil {
		return forecast, fmt.Errorf("failed to parse response: %w", err)
	}

	return forecast, nil
}

type RoundTripper struct {
	http.RoundTripper
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", userAgent)
	return r.RoundTripper.RoundTrip(req)
}
