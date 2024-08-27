/*
 * MIT License
 *
 * Copyright (c) 2024 DearNostalgia
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package schedule

import (
	"github.com/dearnostalgia/gota/pkg/base/bar_series"
	"github.com/dearnostalgia/gota/pkg/base/indicator"
	"sync"
	"sync/atomic"
)

const (
	oldLock = iota
	newLock
)

type EventProcessor[T any] struct {
	indicator indicator.Indicator[T]
	locked    int32
	mu        sync.Mutex
}

func NewEventProcessor[T any](indicator indicator.Indicator[T]) *EventProcessor[T] {
	return &EventProcessor[T]{
		indicator: indicator,
		locked:    oldLock,
		mu:        sync.Mutex{},
	}
}

func (e *EventProcessor[T]) ProcessEvent(event bar_series.BarSeriesEvent) error {
	var bar = *event.Bar
	if bar.IsEnd() {
		return e.ProcessFinalBar(event.Idx)
	}
	return e.ProcessUnstableBar(event.Idx)
}

func (e *EventProcessor[T]) ProcessFinalBar(idx int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, err := e.indicator.Calculate(idx)
	if err != nil {
		return err
	}
	return nil
}

func (e *EventProcessor[T]) ProcessUnstableBar(idx int) error {
	if atomic.CompareAndSwapInt32(&e.locked, oldLock, newLock) {
		defer atomic.StoreInt32(&e.locked, oldLock)
		_, err := e.indicator.Calculate(idx)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
