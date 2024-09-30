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

package binance

import (
	"github.com/dearnostalgia/gota/pkg/base/bar"
	"github.com/govalues/decimal"
	"time"
)

type Option func(*KlineTicker)

type KlineTicker struct {
	bar.BaseBar
	Amount      *decimal.Decimal `json:"amount,omitempty"`
	BuyVolume   *decimal.Decimal `json:"buy_volume,omitempty"`
	BuyAmount   *decimal.Decimal `json:"buy_amount,omitempty"`
	TradesCount *int             `json:"trades_count,omitempty"`
}

func NewBinanceKline(beginTime time.Time, openPrice, closePrice, highPrice, lowPrice decimal.Decimal, opts ...bar.Option) KlineTicker {
	return KlineTicker{
		BaseBar: *bar.NewBaseBar(beginTime, openPrice, closePrice, highPrice, lowPrice, opts...),
	}
}

func (k KlineTicker) WithAmount(amount decimal.Decimal) KlineTicker {
	k.Amount = &amount
	return k
}

func (k KlineTicker) WithBuyVolume(buyVolume decimal.Decimal) KlineTicker {
	k.BuyVolume = &buyVolume
	return k
}

func (k KlineTicker) WithBuyAmount(buyAmount decimal.Decimal) KlineTicker {
	k.BuyAmount = &buyAmount
	return k
}

func (k KlineTicker) WithTradesCount(tradesCount int) KlineTicker {
	k.TradesCount = &tradesCount
	return k
}
