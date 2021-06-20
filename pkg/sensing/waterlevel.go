package sensing

import (
	"context"
	"sync"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/stianeikeland/go-rpio/v4"
)

type WaterLevelSensor struct {
	mu      *sync.Mutex
	trigPin rpio.Pin
	echoPin rpio.Pin
}

func NewWaterLevelSensor() Sensor {
	return &WaterLevelSensor{
		mu:      &sync.Mutex{},
		trigPin: rpio.Pin(23),
		echoPin: rpio.Pin(24),
	}
}

func (s *WaterLevelSensor) Start(ctx context.Context) (<-chan *write.Point, <-chan error, error) {
	return nil, nil, nil
}

func (s *WaterLevelSensor) Read() (float64, error) {
	return s.readTimeOfFlight().Seconds() * 34300.0/ 2, nil
}

func (s *WaterLevelSensor) readTimeOfFlight() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()

	start := time.Now()
	end := time.Now()

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
