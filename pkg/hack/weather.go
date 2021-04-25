package hack

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/mikberg/irrigation/pkg/yr"
	"github.com/rs/zerolog/log"
)

func LogNowcasts(ctx context.Context, influxClient influxdb2.Client, yrClient yr.Client) error {
	writeAPI := influxClient.WriteAPIBlocking("", "irrigation")

	go func() {
		for {
			nowcast, err := yrClient.Nowcast(59.9084295, 10.7785315, 0.0)
			if err != nil {
				log.Error().Err(err).Msg("failed to get nowcast")
			}

			instant, err := nowcast.GetInstant()
			if err != nil {
				log.Error().Err(err).Msg("failed to get instant from nowcast")
			}

			p := influxdb2.NewPointWithMeasurement("weather").
				AddTag("source", "yr").
				SetTime(instant.Time)
			for field, value := range instant.Data.Instant.Details {
				p = p.AddField(field, value)
			}
			writeAPI.WritePoint(ctx, p)
			log.Info().Time("instant", instant.Time).Time("updated_at", nowcast.Properties.Meta.UpdatedAt).Msg("wrote nowcast from yr")

			// sleep until next:
			// there's room for being smarter here. For now, just wait 5 mins
			next := time.Now().Add(3 * time.Minute)

			timer := time.NewTimer(time.Until(next))
			log.Info().Time("next", next).Msg("next measurement expected")

			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				timer.Stop()
				continue
			}
		}
	}()

	return nil
}
