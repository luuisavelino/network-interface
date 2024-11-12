package entities

type Connection struct {
	ErrorRate float64
	Latency   float64
}

func (dw Connection) GetErrorRate() float64 {
	return dw.ErrorRate
}

func (dw Connection) GetLatency() float64 {
	return dw.Latency
}
