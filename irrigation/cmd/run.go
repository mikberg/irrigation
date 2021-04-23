package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the irrigator",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("starting server")

		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("hey ho"))
		})

		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println(err)
		}

		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			log.Debug().Msg("still running")
		}
	},
}
