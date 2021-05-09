// Package sensing contains sensors and tools for storing the sensor data
package sensing

import (
	"context"

	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

// DataPoint represents a datapoint from some sensor. Reusing the Point object
// from influxdb for now.
// type DataPoint write.Point

// Sensor models objects that can continuously output sensor data, until the
// context is cancelled.
type Sensor interface {
	Start(context.Context) (<-chan *write.Point, <-chan error, error)
}
