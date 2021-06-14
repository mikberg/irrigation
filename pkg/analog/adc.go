package analog

import (
	"sync"

	"github.com/stianeikeland/go-rpio/v4"
)

type Channel uint8

const (
	Ch0 = iota
	Ch1
	Ch2
	Ch3
	Ch4
	Ch5
	Ch6
	Ch7
)

type ADC interface {
	Start() error
	Close() error
	Read(Channel) float64
}

func NewADC() ADC {
	return &adc{
		mu:   &sync.Mutex{},
		data: make([]byte, 3),
	}
}

type adc struct {
	mu   *sync.Mutex
	data []byte
}

func (a *adc) Start() error {
	if err := rpio.Open(); err != nil {
		return err
	}

	if err := rpio.SpiBegin(rpio.Spi0); err != nil {
		return err
	}

	rpio.SpiSpeed(1000000)
	rpio.SpiChipSelect(0)

	return nil
}

func (a *adc) Close() error {
	rpio.SpiEnd(rpio.Spi0)
	return rpio.Close()
}

func (a *adc) Read(ch Channel) float64 {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.read(ch)
}

func (a *adc) read(ch Channel) float64 {
	a.data[0], a.data[1], a.data[2] = 1, (8+byte(ch))<<4, 0
	rpio.SpiExchange(a.data)
	code := int(a.data[1]&3)<<8 + int(a.data[2])
	return float64(code) * 3.3 / 1024
}
