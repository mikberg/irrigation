package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/stianeikeland/go-rpio/v4"
)

func init() {

}

var analogCmd = &cobra.Command{
	Use:   "analog",
	Short: "test analog",
	Run: func(cmd *cobra.Command, args []string) {
		if err := rpio.Open(); err != nil {
			log.Fatal().Err(err).Msg("failed to open rpio")
		}
		defer rpio.Close()

		if err := rpio.SpiBegin(rpio.Spi0); err != nil {
			log.Fatal().Err(err).Msg("failed to open spi0")
		}

		rpio.SpiSpeed(1000000)
		rpio.SpiChipSelect(0)
		ticker := time.NewTicker(500 * time.Millisecond)

		w := tabwriter.NewWriter(os.Stderr, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)

		for range ticker.C {
			data := []byte{0, 0, 0}
			for analogCh := uint(0); analogCh < 8; analogCh++ {
				data[0] = 1
				data[1] = (8 + byte(analogCh)) << 4
				data[2] = 0

				rpio.SpiExchange(data)

				code := int(data[1]&3)<<8 + int(data[2])

				voltage := (float32(code) * 3.3) / 1024

				if analogCh == 0 {
					theta_v := 2.820/voltage - 1.014
					fmt.Fprintf(w, "%.2f 𝞱 / %.2fv\t", theta_v, voltage)
				} else {
					fmt.Fprintf(w, "%.2fv\t", voltage)
				}
			}
			fmt.Fprint(w, "\n")
			w.Flush()
		}

		defer rpio.SpiEnd(rpio.Spi0)
	},
}
