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

	testCmdDuration time.Duration
	testCmdChannel  uint
)

func init() {
	testCmd.Flags().DurationVar(&testCmdDuration, "duration", 1*time.Second, "duration to test for")
	testCmd.Flags().UintVar(&testCmdChannel, "channel", 1, "channel to test on")
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "test things",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msgf("testing for %.2f seconds", testCmdDuration.Seconds())

		if err := rpio.Open(); err != nil {
			log.Fatal().Err(err).Msg("failed to open gpio")
		}
		defer rpio.Close()

		waterer := water.NewWaterer(relay1, []rpio.Pin{relay2, relay3, relay4})

		channel := water.Channel(testCmdChannel)
		waterer.Water(channel, testCmdDuration)
	},
}
