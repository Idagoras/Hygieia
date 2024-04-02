package alogrithm

import (
	"Hygieia/util"
)

type RunnerConfig struct {
}

type AlgorithmRunner struct {
	identifier uint64
	hub        *AlgorithmHub
	buf        chan []byte
}

func (r *AlgorithmRunner) Run() {
	defer func() {
		r.hub.Unregister <- r
	}()
	for {
		select {
		case data := <-r.buf:
			_ = data
			level := util.RandomInt(1, 255)
			result := &AlgorithmResult{
				Type:       1,
				IntValue:   level,
				IntSlice:   nil,
				ByteArray:  nil,
				FloatValue: 0,
				FloatSlice: nil,
				Error:      nil,
			}
			r.hub.Res <- result
		}
	}
}
