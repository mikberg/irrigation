package cmd

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/stianeikeland/go-rpio/v4"
)

var (
	trigPin = rpio.Pin(23)
	echoPin = rpio.Pin(24)
)

var distanceCmd = &cobra.Command{
	Use:   "distance",
	Short: "measure distance",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("testing distance")

		if err := rpio.Open(); err != nil {
			log.Fatal().Err(err).Msg("failed to open gpio")
		}
		defer rpio.Close()

		// setup pins
		trigPin.Output()
		echoPin.Input()

		readTimeOfFlight := func() time.Duration {
			start := time.Now()
			end := time.Now()

			trigPin.High()
			time.Sleep(10 * time.Microsecond)
			trigPin.Low()

			for idx := 0; echoPin.Read() == rpio.Low; idx++ {
				if time.Since(start) > 100*time.Millisecond {
					return 0
				}
			}
			start = time.Now()

			for idx := 0; echoPin.Read() == rpio.High; idx++ {
				if time.Since(start) > 200*time.Millisecond {
					return 0
				}
			}
			end = time.Now()

			return end.Sub(start)
		}

		ticker := time.NewTicker(500 * time.Millisecond)
		for range ticker.C {
			timeOfFlight := readTimeOfFlight()
			distance := timeOfFlight.Seconds() * 34300 / 2
			log.Info().Msgf("Tof %d μs, dist %.2f cm", timeOfFlight.Microseconds(), distance)
		}
	},
}
