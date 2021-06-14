package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/mikberg/irrigation/pkg/analog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {

}

var analogCmd = &cobra.Command{
	Use:   "analog",
	Short: "test analog",
	Run: func(cmd *cobra.Command, args []string) {
		adc := analog.NewADC()
		if err := adc.Start(); err != nil {
			log.Fatal().Err(err).Msg("failed to start adc")
		}
		defer adc.Close()

		ticker := time.NewTicker(500 * time.Millisecond)
		w := tabwriter.NewWriter(os.Stderr, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)

		for range ticker.C {
			for analogCh := uint8(0); analogCh < 8; analogCh++ {
				voltage := adc.Read(analog.Channel(analogCh))

				switch analogCh {
				case 1:
					// theta_v := 2.820/voltage - 1.014
					// theta_v := 2.48/voltage - 0.72  // from article
					theta_v := 2.11/voltage - 0.76
					fmt.Fprintf(w, "%.2f 𝞱_v / %.2fv\t", theta_v, voltage)
				case 3:
					temp := 100*voltage - 50
					// temp := 10 * voltage
					fmt.Fprintf(w, "%.2f°C / %.2fv\t", temp, voltage)
				default:
					fmt.Fprintf(w, "%.2fv\t", voltage)
				}
			}
			fmt.Fprint(w, "\n")
			w.Flush()
		}
	},
}
