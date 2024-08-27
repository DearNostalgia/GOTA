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
 */

package okx

import (
	"github.com/dearnostalgia/gota/pkg/base/bar"
	"github.com/govalues/decimal"
	"time"
)

type Option func(*KlineTicker)

type KlineTicker struct {
	bar.BaseBar
	Vol         *decimal.Decimal
	VolCcy      *decimal.Decimal
	VolCcyQuote *decimal.Decimal
}

func NewOkxKline(beginTime time.Time, openPrice, closePrice, highPrice, lowPrice decimal.Decimal) KlineTicker {
	return KlineTicker{
		BaseBar: *bar.NewBaseBar(beginTime, openPrice, closePrice, highPrice, lowPrice),
	}
}

func (k KlineTicker) WithVol(vol decimal.Decimal) KlineTicker {
	k.Vol = &vol
	return k
}

func (k KlineTicker) WithVolCcy(volCcy decimal.Decimal) KlineTicker {
	k.VolCcyQuote = &volCcy
	return k
}

func (k KlineTicker) WithVolCcyQuote(volCcyQuote decimal.Decimal) KlineTicker {
	k.VolCcyQuote = &volCcyQuote
	return k
}
