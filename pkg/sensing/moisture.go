package sensing

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/rs/zerolog/log"
	"github.com/stianeikeland/go-rpio/v4"
)

type MoistureSensor struct {
	interval time.Duration
	analogCh uint
}

func NewMoistureSensor(c uint) Sensor {
	return &MoistureSensor{
		interval: 60 * time.Second,
		analogCh: c,
	}
}

func (s *MoistureSensor) Start(ctx context.Context) (<-chan *write.Point, <-chan error, error) {
	log := log.With().Str("sensor", fmt.Sprintf("moisture(ch%d)", s.analogCh)).Logger()
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
				voltage, err := s.readVoltage()
				if err != nil {
					errc <- err
					continue
				}
				p := influxdb2.NewPointWithMeasurement("moisture").
					AddTag("channel", fmt.Sprintf("%d", s.analogCh)).
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

func (s *MoistureSensor) readVoltage() (float64, error) {
	// @TODO: should be centralized
	if err := rpio.Open(); err != nil {
		return 0.0, fmt.Errorf("failed to open rpio: %w", err)
	}
	defer rpio.Close()

	if err := rpio.SpiBegin(rpio.Spi0); err != nil {
		return 0.0, fmt.Errorf("faield to open spi0: %w", err)
	}

	rpio.SpiSpeed(1000000)
	rpio.SpiChipSelect(0)

	data := []byte{0, 0, 0}
	data[0] = 1
	data[1] = (8 + byte(s.analogCh)) << 4
	data[2] = 0

	rpio.SpiExchange(data)

	code := int(data[1]&3)<<8 + int(data[2])
	voltage := (float64(code) * 3.3) / 1024

	return voltage, nil
}
