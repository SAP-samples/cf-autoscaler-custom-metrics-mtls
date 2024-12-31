package server

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type CPU interface {
	IncreaseCPU(utilization int64, duration time.Duration) error
	StopCPU() error
	isCPURunning() bool
}

type CPUWaster struct {
	isRunning bool
	mutex     sync.RWMutex
}

func (c *CPUWaster) IncreaseCPU(utilization int64, duration time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.isRunning {
		return errors.New("CPU test is already running")
	}
	c.isRunning = true

	go func() {
		defer func() {
			c.mutex.Lock()
			defer c.mutex.Unlock()
			c.isRunning = false
		}()
		endTime := time.Now().Add(duration)
		for time.Now().Before(endTime) {
			if !c.isRunning {
				break
			}
			busyLoop(utilization)
		}
	}()
	return nil
}

func (c *CPUWaster) StopCPU() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if !c.isRunning {
		return errors.New("CPU test is not running")
	}
	c.isRunning = false

	return nil
}

func busyLoop(utilization int64) {
	start := time.Now()
	desiredUtilization := float64(utilization) / 100.0
	fmt.Println("putting load on CPU")
	for time.Since(start).Seconds() < desiredUtilization {
		// Busy loop to simulate CPU load
		fmt.Printf(".")
	}
	time.Sleep(time.Duration(100-utilization) * time.Millisecond)
}
