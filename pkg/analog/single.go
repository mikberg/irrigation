package analog

type Single interface {
	Read() float64
}

func NewSingle(a ADC, ch Channel) Single {
	return &single{a, ch}
}

type single struct {
	a  ADC
	ch Channel
}

func (s *single) Read() float64 {
	return s.a.Read(s.ch)
}
