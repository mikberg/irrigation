package sensing

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/mikberg/irrigation/pkg/analog"
	"github.com/rs/zerolog/log"
)

type MoistureSensor struct {
	adc      analog.Single
	interval time.Duration
	name     string
}

func NewMoistureSensor(name string, adc analog.Single) Sensor {
	return &MoistureSensor{
		interval: 60 * time.Second,
		adc:      adc,
		name:     name,
	}
}

func (s *MoistureSensor) Start(ctx context.Context) (<-chan *write.Point, <-chan error, error) {
	log := log.With().Str("sensor", fmt.Sprintf("moisture(%s)", s.name)).Logger()
	datac := make(chan *write.Point)
	errc := make(chan error)

	go func() {
		defer close(datac)
		defer close(errc)

		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				voltage := s.adc.Read()
				p := influxdb2.NewPointWithMeasurement("moisture").
					AddTag("channel", s.name).
					AddField("voltage", voltage).
					SetTime(time.Now())
				datac <- p
			case <-ctx.Done():
				log.Info().Msg("stopping sensor")
				return
			}
		}
	}()

	return datac, errc, nil
}
