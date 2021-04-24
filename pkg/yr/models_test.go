package yr

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseNowcast(t *testing.T) {
	f, err := os.Open("testdata/nowcast.json")
	require.Nil(t, err)
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	require.Nil(t, err)

	nowcast := METJSONForecast{}
	err = json.Unmarshal(data, &nowcast)
	require.Nil(t, err)

	assert.Equal(t, 8.8, nowcast.Properties.Timeseries[0].Data.Instant.Details["air_temperature"])
	assert.Equal(t, time.Date(2021, 04, 24, 16, 25, 0, 0, time.UTC), nowcast.Properties.Timeseries[0].Time)
	assert.Equal(t, "celsius", nowcast.Properties.Meta.Units["air_temperature"])
}
