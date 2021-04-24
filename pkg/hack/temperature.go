package hack

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/rs/zerolog/log"
)

func LogTemperatures(ctx context.Context, client influxdb2.Client) error {
	writeAPI := client.WriteAPIBlocking("", "irrigation")
	defer client.Close()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			go func() {
				temp, err := readTemperature()
				if err != nil {
					log.Error().Err(err).Msg("failed to read temperature")
				}

				p := influxdb2.NewPointWithMeasurement("temperature").
					AddTag("unit", "celsius").
					AddTag("asset", "pi").
					AddField("value", temp).
					SetTime(time.Now())
				writeAPI.WritePoint(ctx, p)

				log.Info().Float64("temperature", temp).Msg("wrote temperature to database")
			}()
		case <-ctx.Done():
			log.Info().Msg("stopping reading temperatures")
		}
	}
}

func readTemperature() (float64, error) {
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
