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
	"github.com/dearnostalgia/gota/pkg/base/common"
	"github.com/dearnostalgia/gota/pkg/base/indicator"
	"github.com/dearnostalgia/gota/pkg/helper"
	"github.com/govalues/decimal"
)

var _ indicator.CacheIndicator[decimal.Decimal] = (*SMAIndicator)(nil)

type SMAIndicator struct {
	*indicator.EnhancedCacheIndicator[decimal.Decimal]
	period int
}

func NewSMAIndicator(period int) (indicator.CacheIndicator[decimal.Decimal], error) {
	smaIndicator := &SMAIndicator{
		EnhancedCacheIndicator: indicator.NewEnhancedCacheIndicator[decimal.Decimal](),
		period:                 period,
	}
	return smaIndicator, nil
}

func (e *SMAIndicator) Calculate(idx int) (decimal.Decimal, error) {
	var err error
	sum := decimal.Zero
	for i := helper.Max(0, idx-e.period+1); i <= idx; i++ {
		sum, err = sum.Add(e.Source.GetValue(i))
		if err != nil {
			common.IndicatorCalculateErrLog(idx, common.SMA, err)
			return decimal.Zero, err
		}
	}

	realBarCnt, _ := decimal.NewFromFloat64(float64(helper.Min(e.period, idx+1)))
	return sum.Quo(realBarCnt)
}
