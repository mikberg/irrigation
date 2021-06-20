package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/mikberg/irrigation/pkg/analog"
	"github.com/mikberg/irrigation/pkg/sensing"
	"github.com/mikberg/irrigation/pkg/server"
	"github.com/mikberg/irrigation/pkg/water"
	"github.com/mikberg/irrigation/pkg/yr"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/stianeikeland/go-rpio/v4"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the irrigator",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		log.Info().Msg("starting server")

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Let's irrigate!"))
		})
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
		writeAPI := influxClient.WriteAPIBlocking("", "irrigation")

		// Sensors
		yrClient := yr.NewClient()
		yrSensor := sensing.NewYrNowcastSensor(yrClient, 59.9084, 10.7785, 0.0)

		// @TODO: ugly
		go func() {
			if yrdatac, errc, err := yrSensor.Start(ctx); err == nil {
				go func() {
					for {
						select {
						case err := <-errc:
							log.Error().Err(err).Msg("error from weather sensor")
						case point := <-yrdatac:
							writeAPI.WritePoint(ctx, point)
						}
					}

				}()
			} else {
				log.Fatal().Err(err).Msg("failed to start yr sensor")
			}
		}()

		piTempSensor := sensing.NewPiTemperatureSensor()
		go func() {
			if datac, errc, err := piTempSensor.Start(ctx); err == nil {
				go func() {
					for {
						select {
						case err := <-errc:
							log.Error().Err(err).Msg("error from pi sensor")
						case point := <-datac:
							writeAPI.WritePoint(ctx, point)
						}
					}
				}()
			}
		}()

		adc := analog.NewADC()
		if err := adc.Start(); err != nil {
			log.Fatal().Err(err).Msg("failed to start adc")
		}
		defer adc.Close()

		moistureSensor0 := sensing.NewMoistureSensor("0", analog.NewSingle(adc, analog.Ch0))
		go func() {
			if datac, errc, err := moistureSensor0.Start(ctx); err == nil {
				go func() {
					for {
						select {
						case err := <-errc:
							log.Error().Err(err).Msg("error from moisture sensor")
						case point := <-datac:
							writeAPI.WritePoint(ctx, point)
						}
					}
				}()
			}
		}()

		moistureSensor1 := sensing.NewMoistureSensor("1", analog.NewSingle(adc, analog.Ch1))
		go func() {
			if datac, errc, err := moistureSensor1.Start(ctx); err == nil {
				go func() {
					for {
						select {
						case err := <-errc:
							log.Error().Err(err).Msg("error from moisture sensor")
						case point := <-datac:
							writeAPI.WritePoint(ctx, point)
						}
					}
				}()
			}
		}()

		moistureSensor2 := sensing.NewMoistureSensor("2", analog.NewSingle(adc, analog.Ch2))
		go func() {
			if datac, errc, err := moistureSensor2.Start(ctx); err == nil {
				go func() {
					for {
						select {
						case err := <-errc:
							log.Error().Err(err).Msg("error from moisture sensor")
						case point := <-datac:
							writeAPI.WritePoint(ctx, point)
						}
					}
				}()
			}
		}()

		// Watering
		waterer := water.NewWaterer(relay1, []rpio.Pin{relay2, relay3, relay4})

		// Water level
		waterLevelSensor := sensing.NewWaterLevelSensor()
		go func() {
			if datac, errc, err := waterLevelSensor.Start(ctx); err == nil {
				go func() {
					for {
						select {
						case err := <-errc:
							log.Error().Err(err).Msg("error from water level sensor")
						case point := <-datac:
							writeAPI.WritePoint(ctx, point)
						}
					}
				}()
			}
		}()

		// gRPC
		serverConfig := &server.ServerConfig{
			MoistureSensors: map[uint32]*sensing.MoistureSensor{
				0: moistureSensor0.(*sensing.MoistureSensor),
				1: moistureSensor1.(*sensing.MoistureSensor),
				2: moistureSensor2.(*sensing.MoistureSensor),
			},
			Waterer:          waterer,
			WaterLevelSensor: waterLevelSensor.(*sensing.WaterLevelSensor),
		}

		grpcServer := server.NewServer(serverConfig)
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatal().Err(err).Msg("failed to listen")
		}
		go func() {
			log.Info().Msg("starting gRPC server")
			if err := grpcServer.Serve(lis); err != nil {
				log.Fatal().Err(err).Msg("failed to start gRPC server")
			}
		}()

		<-ctx.Done()
	},
}
