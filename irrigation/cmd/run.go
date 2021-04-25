package cmd

import (
	"context"
	"fmt"
	"net/http"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/mikberg/irrigation/pkg/hack"
	"github.com/mikberg/irrigation/pkg/yr"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the irrigator",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		log.Info().Msg("starting server")

		yrClient := yr.NewClient()
		forecast, err := yrClient.Nowcast(59.9084295, 10.7785315, 0.0)
		if err != nil {
			log.Error().Err(err).Msg("failed to query")
		}
		instant, err := forecast.GetInstant()
		if err != nil {
			log.Error().Err(err).Msg("no instant data")
		}
		fmt.Println(instant.Time)
		fmt.Println(instant.Data.Instant.Details)

		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("hey??"))
		})
		go func() {
			log.Info().Msg("starting http server")
			if err := http.ListenAndServe(":8080", nil); err != nil {
				fmt.Println(err)
			}
		}()

		influxClient := influxdb2.NewClient("http://localhost:8086", "irrigation:bluppface")
		defer influxClient.Close()
		go func() {
			if err := hack.LogTemperatures(ctx, influxClient); err != nil {
				log.Error().Err(err).Msg("error logging temperatures")
			}
		}()

		go func() {
			if err := hack.LogNowcasts(ctx, influxClient, yrClient); err != nil {
				log.Error().Err(err).Msg("error logging nowcast")
			}
		}()

		<-ctx.Done()
	},
}
