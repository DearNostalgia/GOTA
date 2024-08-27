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
	"errors"
	"fmt"
	"github.com/dearnostalgia/gota/pkg/base/bar_series"
	"sync"
)

var (
	errorCacheResultsOutRange       = errors.New("index out of range")
	errorBarSeriesNotMatchIndicator = errors.New("bar series not match")
)

type CacheIndicator[T any] interface {
	Indicator[T]
	GetCache() *Cache[T]
}

type EnhancedCacheIndicator[T any] struct {
	*BaseIndicator[T]
	Cache *Cache[T]
}

func NewEnhancedCacheIndicator[T any]() *EnhancedCacheIndicator[T] {
	return &EnhancedCacheIndicator[T]{
		BaseIndicator: NewBasicIndicator[T](),
	}
}

func (b *EnhancedCacheIndicator[T]) GetCache() *Cache[T] {
	return b.Cache
}

func (b *EnhancedCacheIndicator[T]) GetResults() []*Result[T] {
	return b.Cache.GetResults()
}

func (b *EnhancedCacheIndicator[T]) SetSource(s *AttributeSource[T]) {
	b.Source = s
	if s != nil && s.GetBarSeries() != nil && b.Cache == nil {
		b.Cache = NewCache[T](s.GetBarSeries())
	}
}

type Cache[T any] struct {
	barSeries     bar_series.BarSeries
	results       []*Result[T]
	unstableValue T
	mu            sync.RWMutex
}

func NewCache[T any](barSeries bar_series.BarSeries) *Cache[T] {
	cacheIndicator := &Cache[T]{
		barSeries: barSeries,
		results:   make([]*Result[T], 0),
		mu:        sync.RWMutex{},
	}
	return cacheIndicator
}

func (cs *Cache[T]) GetResults() []*Result[T] {
	return cs.results
}

func (cs *Cache[T]) checkPosition() bool {

	barStartTime := cs.barSeries.GetFirstBar().GetBeginTime()

	if len(cs.results) == 0 || cs.results[0].BeginTime.Equal(barStartTime) {
		return true
	}

	if cs.results[0].BeginTime.Before(barStartTime) {
		removeIdx := 0
		for removeIdx < len(cs.results) && cs.results[removeIdx].BeginTime.Before(barStartTime) {
			removeIdx++
		}

		cs.results = cs.results[removeIdx:]
		if len(cs.results) > 0 && cs.results[0].BeginTime.Equal(barStartTime) {
			return true
		}
	}
	return false
}

func (cs *Cache[T]) GetValue(idx int, calcFunc func(int) (T, error)) (T, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	var (
		zero          T
		err           error
		b             = cs.barSeries.GetBar(idx)
		barSeriesSize = cs.barSeries.Size()
	)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("cache get value panic occurred: %v", r)
		}
	}()

	if !cs.checkPosition() {
		return zero, errorBarSeriesNotMatchIndicator
	}

	// Ensure index is within bounds
	if idx < 0 || idx >= barSeriesSize {
		return zero, errorCacheResultsOutRange
	}

	if idx < len(cs.results) {
		tRes := cs.results[idx]
		if tRes.BeginTime != b.GetBeginTime() {
			return zero, errorBarSeriesNotMatchIndicator
		}
		return tRes.Value, nil
	}

	// Calculate the result using the provided function
	result, err := calcFunc(idx)
	if err != nil {
		return zero, err
	}
	return result, nil
}

func (cs *Cache[T]) store(idx int, result T) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if !cs.checkPosition() {
		return errorBarSeriesNotMatchIndicator
	}
	var b = cs.barSeries.GetBar(idx)
	if b.IsEnd() && idx >= len(cs.results) {
		cs.results = append(cs.results, &Result[T]{
			BeginTime: b.GetBeginTime(),
			Value:     result,
		})
		return nil
	}
	cs.unstableValue = result
	return nil
}
