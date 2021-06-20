package sensing

import (
	"context"
	"math"
	"sort"
	"sync"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/stianeikeland/go-rpio/v4"
)

type WaterLevelSensor struct {
	mu       *sync.Mutex
	trigPin  rpio.Pin
	echoPin  rpio.Pin
	interval time.Duration
}

func NewWaterLevelSensor() Sensor {
	return &WaterLevelSensor{
		mu:       &sync.Mutex{},
		trigPin:  rpio.Pin(23),
		echoPin:  rpio.Pin(24),
		interval: 60 * time.Second,
	}
}

func (s *WaterLevelSensor) Start(ctx context.Context) (<-chan *write.Point, <-chan error, error) {
	datac := make(chan *write.Point)
	errc := make(chan error)

	go func() {
		defer close(datac)
		defer close(errc)

		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				volume, err := s.Read()
				if err != nil {
					errc <- err
					continue
				}
				p := influxdb2.NewPointWithMeasurement("waterlevel").
					AddField("volume", volume).
					SetTime(time.Now())
				datac <- p
			case <-ctx.Done():
				return
			}
		}
	}()

	return datac, errc, nil
}

func (s *WaterLevelSensor) Read() (float64, error) {
	// return s.readTimeOfFlight().Seconds() * 34300.0 / 2, nil

	// var sum float64
	// for i := 0; i < 10; i++ {
	// 	sum += s.readTimeOfFlight().Seconds() * 34300 / 2
	// 	time.Sleep(100 * time.Millisecond)
	// }

	values := make([]float64, 9)
	for i := 0; i < 9; i++ {
		values[i] = s.readTimeOfFlight().Seconds() * 34300 / 2
		time.Sleep(100 * time.Millisecond)
	}

	sort.Float64s(values)

	// Empty reading: 25.50
	// Approx 1.496 liters per cm
	return math.Max(0, (25.50-values[5])*1.496), nil

	// return sum / 10, nil
}

func (s *WaterLevelSensor) readTimeOfFlight() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()

	start := time.Now()
	end := time.Now()

	s.trigPin.Low()
	time.Sleep(5 * time.Microsecond)
	s.trigPin.High()
	time.Sleep(10 * time.Microsecond)
	s.trigPin.Low()

	for idx := 0; s.echoPin.Read() == rpio.Low; idx++ {
		if time.Since(start) > 100*time.Millisecond {
			return 0
		}
	}
	start = time.Now()

	for idx := 0; s.echoPin.Read() == rpio.High; idx++ {
		if time.Since(start) > 200*time.Millisecond {
			return 0
		}
	}
	end = time.Now()

	return end.Sub(start)
}
