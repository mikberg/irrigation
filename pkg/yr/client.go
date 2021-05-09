package yr

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gregjones/httpcache"
	"github.com/rs/zerolog/log"
)

var (
	userAgent = "irrigationGo github.com/mikbe/irrigation"
)

type Client interface {
	GetNowcast(ctx context.Context, lat, lon, alt float64) (*METJSONForecast, time.Time, error)
}

// client is an YR client, obseriving YR API best practices.
//
// TL;DR:
//  * Use a good browser agent
//  * The "Expires" header
//  * Use the "If-Modified-Since" header on requests
//  * Coordinates: don't use more than 4 decimals
type client struct {
	baseURL    string
	httpClient *http.Client
	cache      httpcache.Cache
}

func NewClient() Client {
	cache := httpcache.NewMemoryCache()
	cacheTransport := httpcache.NewTransport(cache)

	httpClient := &http.Client{
		Transport: cacheTransport,
	}

	return &client{
		httpClient: httpClient,
		cache:      cache,
		baseURL:    "https://api.met.no",
	}
}

func (c *client) GetNowcast(ctx context.Context, lat, lon, alt float64) (*METJSONForecast, time.Time, error) {
	values := url.Values{}
	values.Add("lat", fmt.Sprintf("%.4f", lat))
	values.Add("lon", fmt.Sprintf("%.4f", lon))
	if alt != 0.0 {
		values.Add("alt", fmt.Sprintf("%.4f", alt))
	}

	res, err := c.request(ctx, "/weatherapi/nowcast/2.0/complete", values)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to request nowcast: %w", err)
	}

	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to read body: %w", err)
	}

	nowcast := &METJSONForecast{}
	if err := json.Unmarshal(data, nowcast); err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to pase nowcast: %w", err)
	}

	expiry, _ := time.Parse(http.TimeFormat, res.Header.Get("Expires"))

	return nowcast, expiry, nil
}

func (c *client) request(ctx context.Context, path string, values url.Values) (*http.Response, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}
	u.Path = path
	u.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	// set the user agent
	req.Header.Set("User-Agent", userAgent)

	log := log.With().
		Str("method", req.Method).
		Stringer("url", req.URL).
		Logger()
	log.Debug().Msg("request")

	// cached response?
	var cachedRes *http.Response
	if res, err := httpcache.CachedResponse(c.cache, req); err == nil && res != nil {
		req.Header.Set("If-Modified-Since", res.Header.Get("Last-Modified"))
		cachedRes = res
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	log.Debug().Str("x-from-cache", res.Header.Get(httpcache.XFromCache)).Msg(res.Status)

	switch res.StatusCode {
	case http.StatusOK:
		return res, nil
	case http.StatusNotModified:
		return cachedRes, nil
	case http.StatusNonAuthoritativeInfo:
		log.Warn().Msg("product is in beta or deprecated")
		return res, nil
	case http.StatusForbidden:
		return nil, fmt.Errorf("403 forbidden: check user agent")
	case http.StatusUnprocessableEntity:
		return nil, fmt.Errorf("422 unprocessable entity: location not covered?")
	default:
		return nil, fmt.Errorf("unexpected error from api: %d", res.StatusCode)
	}
}
