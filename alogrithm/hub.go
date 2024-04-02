package alogrithm

var Hub *AlgorithmHub = NewAlgorithmHub()

const (
	Stop        = 1
	StopAndSave = iota << 1
)

const (
	EEGFatigueLevel = 0
	CalculateError  = iota
)

const (
	EEGData = iota
)

type AlgorithmData struct {
	DataType   int
	Data       []byte
	identifier uint64
}

type AlgorithmResult struct {
	Type       uint8
	IntValue   int64
	IntSlice   []int64
	ByteArray  []byte
	FloatValue float64
	FloatSlice []float64
	Error      error
}

type AlgorithmHub struct {
	runners    map[uint64]*AlgorithmRunner
	Register   chan *AlgorithmRunner
	Unregister chan *AlgorithmRunner
	Broadcast  chan int
	Res        chan *AlgorithmResult
	Data       chan *AlgorithmData
}

func NewAlgorithmHub() *AlgorithmHub {
	return &AlgorithmHub{
		runners:    make(map[uint64]*AlgorithmRunner),
		Register:   make(chan *AlgorithmRunner),
		Unregister: make(chan *AlgorithmRunner),
		Broadcast:  make(chan int),
		Res:        make(chan *AlgorithmResult, 512),
		Data:       make(chan *AlgorithmData, 256),
	}
}

func (h *AlgorithmHub) Run() {
	for {
		select {
		case runner := <-h.Register:
			h.runners[runner.identifier] = runner
		case runner := <-h.Unregister:
			if _, ok := h.runners[runner.identifier]; ok {
				delete(h.runners, runner.identifier)
				close(runner.buf)
			}
		case result := <-h.Res:
			_ = result
		case message := <-h.Broadcast:
			_ = message
		case message := <-h.Data:
			id := message.identifier
			if runner, ok := h.runners[id]; ok {
				runner.buf <- message.Data
			}
		}
	}
}
