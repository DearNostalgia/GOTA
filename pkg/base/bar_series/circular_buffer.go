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

package bar_series

import (
	"errors"
	"github.com/dearnostalgia/gota/pkg/base/bar"
	"github.com/dearnostalgia/gota/pkg/base/contact"
	"github.com/google/uuid"
	clone "github.com/huandu/go-clone/generic"
	"sync"
)

var _ BarSeries = (*CircularBarSeries)(nil)

// CircularBarSeries implements the BarSeries interface and manages a circular buffer
// of bar data, ensuring thread-safe operations and efficient storage.
type CircularBarSeries struct {
	*BarSeriesMetaInfo
	relay          *contact.Relay[BarSeriesEvent]
	circularBuffer *circularBuffer
	mu             sync.RWMutex
}

// circularBuffer defines the structure for storing bars in a circular manner,
// which is efficient for operations where old bars are overwritten by new ones.
type circularBuffer struct {
	recoverIdx int
	cnt        int
	inCircular bool
	bars       []bar.Bar
}

func NewCircularBarSeries(opts ...BaseSeriesOption) *CircularBarSeries {
	circularBarSeries := &CircularBarSeries{
		BarSeriesMetaInfo: &BarSeriesMetaInfo{},
		circularBuffer:    &circularBuffer{},
		relay:             contact.NewRelay[BarSeriesEvent](),
		mu:                sync.RWMutex{},
	}
	for _, opt := range opts {
		opt(circularBarSeries.BarSeriesMetaInfo)
	}

	if circularBarSeries.BarSeriesMetaInfo.symbol == nil {
		symbol := uuid.New().String()
		circularBarSeries.BarSeriesMetaInfo.symbol = &symbol
	}

	maxSize := circularBarSeries.BarSeriesMetaInfo.maxSize
	if maxSize == nil || *maxSize <= 0 {
		circularBarSeries.circularBuffer.bars = make([]bar.Bar, 0)
		return circularBarSeries
	}

	circularBarSeries.circularBuffer.bars = make([]bar.Bar, *maxSize)
	return circularBarSeries
}

func (c *CircularBarSeries) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.circularBuffer.cnt
}

func (c *CircularBarSeries) IsEmpty() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.circularBuffer.cnt == 0
}

func (c *CircularBarSeries) GetBar(idx int) bar.Bar {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.maxSize == nil {
		if idx >= 0 && idx < len(c.circularBuffer.bars) {
			return c.circularBuffer.bars[idx]
		}
		return nil
	}

	if idx >= c.circularBuffer.cnt || idx < 0 {
		return nil
	}
	return c.circularBuffer.bars[(c.circularBuffer.recoverIdx+idx)%(c.circularBuffer.cnt)]
}

func (c *CircularBarSeries) GetFirstBar() bar.Bar {
	return c.GetBar(0)
}

func (c *CircularBarSeries) GetLastBar() bar.Bar {
	return c.GetBar(c.circularBuffer.cnt - 1)
}

func (c *CircularBarSeries) AddBar(realTimeUpdateBar bool, bar ...bar.Bar) error {
	if bar == nil {
		return errors.New("bar cannot be nil")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.maxSize != nil && *c.maxSize <= 0 {
		return nil
	}

	for _, b := range bar {
		size := c.circularBuffer.cnt
		if size == 0 {
			c.addBar(b)
			return nil
		}

		var realLastIdx int
		if c.maxSize != nil {
			realLastIdx = (c.circularBuffer.recoverIdx + size) % size
		} else {
			realLastIdx = size - 1
		}

		lastBar := c.circularBuffer.bars[realLastIdx]

		if realTimeUpdateBar && (lastBar.GetBeginTime().Equal(b.GetBeginTime()) || !b.IsEnd()) {
			c.updateBar(b, realLastIdx)
			return nil
		}
		c.addBar(b)
	}
	return nil
}

func (c *CircularBarSeries) updateBar(bar bar.Bar, idx int) {
	c.circularBuffer.bars[idx] = bar
	c.publish(&BarSeriesEvent{
		Idx:                c.circularBuffer.cnt - 1,
		Bar:                bar,
		BarSeriesEventType: UpdateLatestBar,
	})
}

func (c *CircularBarSeries) addBar(bar bar.Bar) {
	defer func() {
		c.publish(&BarSeriesEvent{
			Idx:                c.circularBuffer.cnt - 1,
			Bar:                bar,
			BarSeriesEventType: AddBarToFront,
		})
	}()

	if c.maxSize == nil {
		c.circularBuffer.bars = append(c.circularBuffer.bars, bar)
		c.circularBuffer.cnt++
		return
	}
	c.circularBuffer.recoverIdx = c.circularBuffer.recoverIdx % (*c.maxSize)
	c.circularBuffer.bars[c.circularBuffer.recoverIdx] = bar
	c.circularBuffer.recoverIdx++
	if !c.circularBuffer.inCircular {
		c.circularBuffer.cnt++
	}
	if c.circularBuffer.recoverIdx == *c.maxSize {
		c.circularBuffer.inCircular = true
	}
}

func (c *CircularBarSeries) Subscribe(opts ...SubscribeOption) *contact.Listener[BarSeriesEvent] {
	c.mu.Lock()
	defer c.mu.Unlock()

	options := &SubscribeOptions{
		Size: DefaultDescribeCapacity,
	}
	for _, opt := range opts {
		opt(options)
	}
	l := c.relay.Listener(options.Size)
	return l
}

func (c *CircularBarSeries) publish(e *BarSeriesEvent) {
	if e.Bar.IsEnd() {
		c.relay.Notify(*e)
		return
	}
	c.relay.Broadcast(*e)
}

func (c *CircularBarSeries) GetBarsCopy() []bar.Bar {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return clone.Wrap(c.circularBuffer.bars)
}
