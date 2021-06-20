package cmd

import (
	"time"

	"github.com/mikberg/irrigation/pkg/sensing"
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

		waterLevelSensor := sensing.NewWaterLevelSensor().(*sensing.WaterLevelSensor)

		ticker := time.NewTicker(500 * time.Millisecond)
		for range ticker.C {
			distance, err := waterLevelSensor.Read()
			if err != nil {
				log.Error().Err(err).Msg("failed to read distance")
			}
			log.Info().Msgf("Distance %.2f cm", distance)
		}
	},
}
