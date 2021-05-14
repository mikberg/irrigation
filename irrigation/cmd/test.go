package cmd

import (
	"time"

	"github.com/mikberg/irrigation/pkg/water"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/stianeikeland/go-rpio/v4"
)

var (
	relay1 = rpio.Pin(14)
	relay2 = rpio.Pin(15)
	relay3 = rpio.Pin(18)
	relay4 = rpio.Pin(17)
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "test things",
	Run: func(cmd *cobra.Command, args []string) {
		if err := rpio.Open(); err != nil {
			log.Fatal().Err(err).Msg("failed to open gpio")
		}
		defer rpio.Close()

		waterer := water.NewWaterer(relay1, []rpio.Pin{relay2, relay3, relay4})
		channels := []water.Channel{water.Channel(0), water.Channel(1), water.Channel(2)}

		ticker := time.NewTicker(1 * time.Second)
		idx := 0
		for range ticker.C {
			waterer.Water(channels[idx], 500*time.Millisecond)
			idx = (idx + 1) % len(channels)
		}
	},
}
