package yr

import (
	"errors"
	"time"
)

var (
	ErrNoInstant = errors.New("no instant data")
)

type METJSONForecast struct {
	Type       string        `json:"type"`
	Geometry   PointGeometry `json:"geometry"`
	Properties Forecast      `json:"properties"`
}

func (f METJSONForecast) GetInstant() (*ForecastTimeStep, error) {
	for _, d := range f.Properties.Timeseries {
		if len(d.Data.Instant.Details) > 0 {
			return &d, nil
		}
	}

	return nil, ErrNoInstant
}

type PointGeometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Forecast struct {
	Timeseries []ForecastTimeStep `json:"timeseries"`
	Meta       ForecastMeta       `json:"meta"`
}

type ForecastMeta struct {
	RadarCoverage string        `json:"radar_coverage"`
	Units         ForecastUnits `json:"units"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type ForecastTimeStep struct {
	// The time these forecast values are valid for
	Time time.Time            `json:"time"`
	Data ForecastTimeStepData `json:"data"`
}

type ForecastTimeStepData struct {
	Instant     ForecastTimeStepDataValue `json:"instant"`
	Next1Hours  ForecastTimeStepDataValue `json:"next_1_hours"`
	Next6Hours  ForecastTimeStepDataValue `json:"next_6_hours"`
	Next12Hours ForecastTimeStepDataValue `json:"next_12_hours"`
}

type ForecastTimeStepDataValue struct {
	Summary ForecastSummary `json:"summary"`
	Details ForecastData    `json:"details"`
}

type ForecastSummary struct {
	SymbolCode string `json:"symbol_code"`
}

type ForecastUnits map[string]string
type ForecastData map[string]float64
type WeatherSymbol string
