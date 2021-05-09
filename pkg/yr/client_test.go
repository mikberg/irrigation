package yr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetNowcast(t *testing.T) {
	// Things to validate:
	// * using user agent
	// * caches responses
	// * returns from cache if not yet expired
	// * sends with not-before

	reqCount := 0
	serverExpires := time.Now().Add(1 * time.Minute)
	serverModified := time.Now().Add(-1 * time.Minute)

	ctx := context.Background()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCount++

		// our own user agent must be set
		assert.Equal(t, userAgent, r.UserAgent())

		fmt.Println("if-modified-since header", r.Header.Get("If-Modified-Since"))

		// correct formatting of query
		assert.Equal(t, "1.0000", r.URL.Query().Get("lat"))
		assert.Equal(t, "1.2346", r.URL.Query().Get("lon"))
		assert.Equal(t, "", r.URL.Query().Get("alt"))

		// check the path
		assert.Equal(t, r.URL.Path, "/weatherapi/nowcast/2.0/complete")

		// set expiry header
		lastModified := time.Now().Add(-1 * time.Minute)
		w.Header().Set("Last-Modified", serverModified.Format(http.TimeFormat))
		w.Header().Set("Expires", serverExpires.Format(http.TimeFormat))

		// send a response
		nc := METJSONForecast{
			Type: "forecast",
			Properties: Forecast{
				Meta: ForecastMeta{
					UpdatedAt: lastModified,
				},
			},
		}
		d, err := json.Marshal(nc)
		require.Nil(t, err)
		_, err = w.Write(d)
		require.Nil(t, err)

	}))
	defer ts.Close()

	c := NewClient()

	c.(*client).baseURL = ts.URL

	nowcast, expires, err := c.GetNowcast(ctx, 1.0, 1.23456789, 0.0)
	require.Nil(t, err)
	assert.Equal(t, serverExpires.Unix(), expires.Unix())
	assert.NotNil(t, nowcast)
	assert.Equal(t, 1, reqCount, "should have done 1 request")

	nowcast, expires, err = c.GetNowcast(ctx, 1.0, 1.23456789, 0.0)
	require.Nil(t, err)
	assert.Equal(t, serverExpires.Unix(), expires.Unix())
	assert.NotNil(t, nowcast)
	assert.Equal(t, 1, reqCount, "should still have done only 1 request")
}
