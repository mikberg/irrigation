package sensing

import (
	"context"
	"math/rand"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/mikberg/irrigation/pkg/yr"
	"github.com/rs/zerolog/log"
)

type YrNowcastSensor struct {
	client yr.Client
	lat    float64
	lon    float64
	alt    float64
}

// NewYrNowcastSensor creates a new sensor measuring the current weather at some
// location using an Yr client.
func NewYrNowcastSensor(client yr.Client, lat, lon, alt float64) Sensor {
	return &YrNowcastSensor{
		client: client,
		lat:    lat,
		lon:    lon,
		alt:    alt,
	}
}

func (s *YrNowcastSensor) Start(ctx context.Context) (<-chan *write.Point, <-chan error, error) {
	log.With().Str("sensor", "yrnowcast").Logger()
	datac := make(chan *write.Point)
	errc := make(chan error)

	getPoint := func() (point *write.Point, expires time.Time, err error) {
		nowcast, expires, err := s.client.GetNowcast(ctx, s.lat, s.lon, s.alt)
		if err != nil {
			return
		}

		instant, err := nowcast.GetInstant()
		if err != nil {
			return
		}

		point = influxdb2.NewPointWithMeasurement("weather").
			AddTag("source", "yr").
			SetTime(instant.Time)
		for field, value := range instant.Data.Instant.Details {
			point = point.AddField(field, value)
		}

		return
	}

	go func() {
		defer close(datac)
		defer close(errc)

		for {
			point, expires, err := getPoint()
			if err != nil {
				log.Error().Err(err).Msg("")
				errc <- err
			} else if point != nil {
				log.Debug().
					Str("sensor", "yrnowcast").
					Time("expires", expires).
					Time("validFrom", point.Time()).
					Dur("freshness", time.Since(point.Time())).
					Msg("produced data")
				datac <- point
			}

			// add some jitter; wait at least 30 secs
			wait := time.Until(expires.Add(time.Duration(10 * rand.Float64() * float64(time.Second))))
			if wait < 30*time.Second {
				wait = 30 * time.Second
			}

			// wait until it expires or context is done
			select {
			case <-ctx.Done():
				return
			case <-time.After(wait):
			}
		}
	}()

	return datac, errc, nil
}
