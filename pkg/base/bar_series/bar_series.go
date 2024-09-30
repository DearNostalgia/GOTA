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
	"github.com/dearnostalgia/gota/pkg/base/contact"
	"time"

	"github.com/dearnostalgia/gota/pkg/base/bar"
)

type BarSeriesEventType int

const (
	AddBarToFront BarSeriesEventType = iota + 1
	AddBarToEnd
	UpdateLatestBar
	RemoveInvalidBar
	InsertInMiddle
)

type BarSeriesEvent struct {
	Bar bar.Bar
	BarSeriesEventType
	Idx int
}

type SubscribeOptions struct {
	Size int
}

type SubscribeOption func(*SubscribeOptions)

func WithSubscribeBufferSize(size int) SubscribeOption {
	return func(o *SubscribeOptions) {
		if size > 0 {
			o.Size = size
		}
	}
}

type BarSeries interface {
	// Size returns the current number of bars in the series.
	// It indicates how many data points (bars) are currently stored in the series.
	Size() int

	// IsEmpty checks if the bar series is empty.
	// It returns true if the series contains no bars, and false otherwise.
	IsEmpty() bool

	// GetBar returns the Bar at the specified index.
	// The index is zero-based, meaning that the first bar in the series is at index 0.
	// If the index is out of range, it returns nil.
	GetBar(idx int) bar.Bar

	// GetFirstBar returns the first Bar in the series.
	// If the series is empty, it returns nil.
	// The first bar is the one that was added to the series first, based on time.
	GetFirstBar() bar.Bar

	// GetLastBar returns the last Bar in the series.
	// If the series is empty, it returns nil.
	// The last bar is the most recently added bar in the series.
	GetLastBar() bar.Bar

	// AddBar adds a bar to the series.
	// If realTimeUpdateBar is true, the new data will be compared with the latest data in the bar series.
	// If the beginTime is the same, the existing data will be updated; if it is different, the new data will be added.
	// If realTimeUpdateBar is false, new data will be continuously added.
	AddBar(realTimeUpdateBar bool, bar ...bar.Bar) error

	// GetBarsCopy returns a deep copy of the bar.Bar slice.
	// This method ensures that the returned slice contains independent copies of
	// the elements from the original barSlice, so any modifications made to the
	// returned slice will not affect the original BarSeries structure or its data.
	GetBarsCopy() []bar.Bar

	// Subscribe registers a subscription to the barSeries data and returns a contact.Listener[BarSeriesEvent].
	// This method allows the subscriber to listen for broadcasted data events within the CircularBarSeries.
	// The `Subscribe` method accepts optional parameters to configure the subscription, such as the buffer size for the listener's channel.
	// Usage:
	//     l := barSeries.Subscribe(bar_series.WithSubscribeBufferSize(10))
	//     for event := range l.Ch() {
	//         // handle event
	//     }
	// The CircularBarSeries broadcasts data to all subscribed goroutines.
	Subscribe(opts ...SubscribeOption) *contact.Listener[BarSeriesEvent]
}

type BaseSeriesOption func(*BarSeriesMetaInfo)

type BarSeriesMetaInfo struct {
	symbol   *string
	interval *int64
	maxSize  *int
}

// WithMaxSize sets the maximum size of the bar series.
// If the number of bars exceeds maxSize, the oldest bars will be discarded.
// WithMaxSize(size int) BarSeries
func WithMaxSize(size int) BaseSeriesOption {
	return func(bs *BarSeriesMetaInfo) {
		bs.maxSize = &size
	}
}

// WithSymbol sets the symbol of the bar series.
// Example:
//
//	WithSymbol("BTCUSDT.P")
//	WithSymbol("EURUSD")
func WithSymbol(symbol string) BaseSeriesOption {
	return func(bs *BarSeriesMetaInfo) {
		bs.symbol = &symbol
	}
}

// WithInterval sets the interval of the bar series.
// Example:
//
//	15m:WithInterval(15 * time.Minute)
//	1h:WithInterval(1 * time.Hour)
//	2h:WithInterval(2 * time.Hour)
//	4h:WithInterval(4 * time.Hour)
//	1d:WithInterval(1 * time.Day)
func WithInterval(interval time.Duration) BaseSeriesOption {
	return func(bs *BarSeriesMetaInfo) {
		intervalNano := interval.Nanoseconds()
		bs.interval = &intervalNano
	}
}

// GetSymbol returns the symbol of the bar series
func (b *BarSeriesMetaInfo) GetSymbol() string {
	return *b.symbol
}

// GetInterval returns the interval of the bar series.
func (b *BarSeriesMetaInfo) GetInterval() time.Duration {
	return time.Duration(*b.interval)
}

// GetMaxSize returns the max size of the bar series.
func (b *BarSeriesMetaInfo) GetMaxSize() int {
	return *b.maxSize
}
