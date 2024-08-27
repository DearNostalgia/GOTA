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

package indicator

import (
	"github.com/dearnostalgia/gota/pkg/base/bar"
	"github.com/dearnostalgia/gota/pkg/base/bar_series"
	"github.com/govalues/decimal"
)

type AttributeStrategy[T any] interface {
	GetAttribute(bar bar.Bar) T
}

type AttributeSource[T any] struct {
	series   bar_series.BarSeries
	strategy AttributeStrategy[T]
}

// NewAttributeSource creates a new instance of AttributeSource, which is used to retrieve a specific attribute
// from a BarSeries. The generic parameter T represents the type of the attribute, such as decimal.Decimal or
// any other type you require.
//
// Parameters:
//   - series: bar_series.BarSeries is a dataset containing a series of bars, typically representing time-series data.
//   - strategy: AttributeStrategy[T] is a strategy interface that defines how to extract a specific type of
//     attribute from a given bar.
//
// Returns:
//   - A pointer to an AttributeSource[T] instance, which allows you to use the provided strategy to retrieve
//     the specified type of attribute from the BarSeries.
//
// Example usage:
// var sliceBarSeries bar_series.BarSeries
// closeSource := NewAttributeSource[decimal.Decimal](sliceBarSeries, &ClosePriceStrategy{})
// closeValue := closeSource.GetValue(0)
func NewAttributeSource[T any](series bar_series.BarSeries, strategy AttributeStrategy[T]) *AttributeSource[T] {
	return &AttributeSource[T]{
		series:   series,
		strategy: strategy,
	}
}

func (as *AttributeSource[T]) GetValue(idx int) T {
	return as.strategy.GetAttribute(as.series.GetBar(idx))
}

func (as *AttributeSource[T]) GetBarSeries() bar_series.BarSeries {
	return as.series
}

type ClosePriceStrategy struct{}

func (s *ClosePriceStrategy) GetAttribute(bar bar.Bar) decimal.Decimal {
	return bar.GetClosePrice()
}

type OpenPriceStrategy struct{}

func (s *OpenPriceStrategy) GetAttribute(bar bar.Bar) decimal.Decimal {
	return bar.GetOpenPrice()
}

type HighPriceStrategy struct{}

func (s *HighPriceStrategy) GetAttribute(bar bar.Bar) decimal.Decimal {
	return bar.GetHighPrice()
}

type LowPriceStrategy struct{}

func (s *LowPriceStrategy) GetAttribute(bar bar.Bar) decimal.Decimal {
	return bar.GetLowPrice()
}

type VolumeStrategy struct{}

func (s *VolumeStrategy) GetAttribute(bar bar.Bar) decimal.Decimal {
	return bar.GetVolume()
}
