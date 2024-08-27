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

package bar

import (
	"time"

	"github.com/govalues/decimal"
)

// Bar represents a financial market data bar, typically used in time-series analysis.
type Bar interface {
	// GetBeginTime returns the starting time of the bar.
	GetBeginTime() time.Time

	// GetEndTime returns the ending time of the bar.
	GetEndTime() time.Time

	// GetOpenPrice returns the opening price of the bar.
	GetOpenPrice() decimal.Decimal

	// GetHighPrice returns the highest price during the bar's time period.
	GetHighPrice() decimal.Decimal

	// GetLowPrice returns the lowest price during the bar's time period.
	GetLowPrice() decimal.Decimal

	// GetClosePrice returns the closing price of the bar.
	GetClosePrice() decimal.Decimal

	// GetVolume returns the trading volume during the bar's time period.
	GetVolume() decimal.Decimal

	// IsEnd indicates whether the bar represents the end of a time period,
	// such as the last bar in a trading day or the final bar in a series.
	IsEnd() bool
}
