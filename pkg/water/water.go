package water

import (
	"errors"
	"sync"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

var (
	ErrNoValve = errors.New("no such valve")
)

type Channel uint8

type Waterer interface {
	Water(c Channel, dur time.Duration) error
}

type waterer struct {
	mut       sync.Mutex
	pumpPin   rpio.Pin
	valvePins []rpio.Pin
}

func NewWaterer(pumpPin rpio.Pin, valvePins []rpio.Pin) Waterer {
	// set pins to be output
	pumpPin.Output()
	for _, p := range valvePins {
		p.Output()
	}

	// reset pins
	pumpPin.Low()
	for _, p := range valvePins {
		p.Low()
	}

	return &waterer{
		pumpPin:   pumpPin,
		valvePins: valvePins,
	}
}

func (w *waterer) Water(ch Channel, dur time.Duration) error {
	w.mut.Lock()
	defer w.mut.Unlock()

	if int(ch) > len(w.valvePins) {
		return ErrNoValve
	}
	valvePin := w.valvePins[ch]

	// open the valve
	valvePin.High()
	defer valvePin.Low()
	time.Sleep(100 * time.Millisecond)

	// start pumping
	w.pumpPin.High()
	defer w.pumpPin.Low()

	time.Sleep(dur)

	return nil
}
