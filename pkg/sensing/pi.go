package sensing

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/rs/zerolog/log"
)

// PiTemperatureSensor measures the temperature on the Raspberry Pi.
type PiTemperatureSensor struct {
	interval time.Duration
}

func NewPiTemperatureSensor() Sensor {
	return &PiTemperatureSensor{
		interval: 60 * time.Second,
	}
}

func (s *PiTemperatureSensor) Start(ctx context.Context) (<-chan *write.Point, <-chan error, error) {
	log.With().Str("sensor", "pitemperature").Logger()
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
				temp, err := s.readTemperature()
				if err != nil {
					errc <- err
					continue
				}
				p := influxdb2.NewPointWithMeasurement("temperature").
					AddTag("unit", "celsius").
					AddTag("asset", "pi").
					AddField("value", temp).
					SetTime(time.Now())
				datac <- p
			case <-ctx.Done():
				log.Info().Msg("stopping reading temperatures")
				return
			}
		}

	}()

	return datac, errc, nil
}

func (s *PiTemperatureSensor) readTemperature() (float64, error) {
	f, err := os.Open("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		return 0, fmt.Errorf("failed to open temperature file: %w", err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return 0, fmt.Errorf("failed to read temperature file: %w", err)
	}

	tempMilliC, err := strconv.ParseInt(strings.TrimSpace(string(b)), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse temperature: %w", err)
	}

	return float64(tempMilliC) / 1000, nil
}
