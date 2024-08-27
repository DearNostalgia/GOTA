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

package tech

import (
	"errors"
	"github.com/dearnostalgia/gota/pkg/base/common"
	"github.com/dearnostalgia/gota/pkg/base/indicator"
	"github.com/govalues/decimal"
)

var _ indicator.CacheIndicator[decimal.Decimal] = (*EMAIndicator)(nil)

type EMAIndicator struct {
	*indicator.EnhancedCacheIndicator[decimal.Decimal]
	period     int
	multiplier decimal.Decimal
}

func NewEMAIndicator(period int) (indicator.CacheIndicator[decimal.Decimal], error) {
	multiplier, err := decimal.NewFromFloat64(2.0 / float64(period+1))
	if err != nil {
		common.IndicatorCalculateErrLog(0, common.EMA, errors.New("period:"+err.Error()))
		return nil, err
	}
	emaIndicator := &EMAIndicator{
		EnhancedCacheIndicator: indicator.NewEnhancedCacheIndicator[decimal.Decimal](),
		period:                 period,
		multiplier:             multiplier,
	}
	return emaIndicator, nil
}

func (e *EMAIndicator) Calculate(idx int) (decimal.Decimal, error) {
	var (
		sv  = e.Source.GetValue(idx)
		err error
	)

	if idx == 0 {
		return sv, nil
	}

	prevEMA, err := e.Cache.GetValue(idx-1, e.Calculate)

	if err != nil {
		common.IndicatorCalculateErrLog(idx, common.EMA, err)
		return decimal.Zero, err
	}

	res, err := e.calculateEMA(prevEMA, sv)
	if err != nil {
		common.IndicatorCalculateErrLog(idx, common.EMA, err)
		return decimal.Zero, err
	}
	return res, nil
}

func (e *EMAIndicator) calculateEMA(prevEMA decimal.Decimal, sourceVal decimal.Decimal) (decimal.Decimal, error) {
	s, err := sourceVal.Sub(prevEMA)
	if err != nil {
		return decimal.Zero, err
	}
	m, err := s.Mul(e.multiplier)
	if err != nil {
		return decimal.Zero, err
	}

	res, err := m.Add(prevEMA)
	if err != nil {
		return decimal.Zero, err
	}
	return res, nil
}
