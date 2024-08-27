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

type BarOption func(*BaseBar)

var _ Bar = (*BaseBar)(nil)

type BaseBar struct {
	BeginTime   time.Time        `json:"begin_time"`
	EndTime     time.Time        `json:"end_time"`
	OpenPrice   decimal.Decimal  `json:"open_price"`
	ClosePrice  decimal.Decimal  `json:"close_price"`
	HighPrice   decimal.Decimal  `json:"high_price"`
	LowPrice    decimal.Decimal  `json:"low_price"`
	Volume      *decimal.Decimal `json:"volume"`
	Amount      *decimal.Decimal `json:"amount,omitempty"`
	BuyVolume   *decimal.Decimal `json:"buy_volume,omitempty"`
	BuyAmount   *decimal.Decimal `json:"buy_amount,omitempty"`
	TradesCount *decimal.Decimal `json:"trades_count,omitempty"`
	End         bool             `json:"is_end,omitempty"`
}

func NewBaseBar(beginTime time.Time, openPrice, closePrice, highPrice, lowPrice decimal.Decimal, opts ...BarOption) *BaseBar {
	bar := &BaseBar{
		BeginTime:  beginTime,
		OpenPrice:  openPrice,
		ClosePrice: closePrice,
		HighPrice:  highPrice,
		LowPrice:   lowPrice,
	}
	for _, opt := range opts {
		opt(bar)
	}
	return bar
}

func WithEndTime(time time.Time) BarOption {
	return func(b *BaseBar) {
		b.EndTime = time
	}
}

func WithVolume(volume decimal.Decimal) BarOption {
	return func(b *BaseBar) {
		b.Volume = &volume
	}
}

func WithAmount(amount decimal.Decimal) BarOption {
	return func(b *BaseBar) {
		b.Amount = &amount
	}
}

func WithIsEnd(end bool) BarOption {
	return func(b *BaseBar) {
		b.End = end
	}
}

func WithBuyVolume(buyVolume decimal.Decimal) BarOption {
	return func(b *BaseBar) {
		b.BuyVolume = &buyVolume
	}
}

func WithBuyAmount(buyAmount decimal.Decimal) BarOption {
	return func(b *BaseBar) {
		b.BuyAmount = &buyAmount
	}
}

func WithTradesCount(tradesCount decimal.Decimal) BarOption {
	return func(b *BaseBar) {
		b.TradesCount = &tradesCount
	}
}

func (b *BaseBar) IsEnd() bool {
	return b.End
}

func (b *BaseBar) GetBeginTime() time.Time {
	return b.BeginTime
}

func (b *BaseBar) GetEndTime() time.Time {
	return b.EndTime
}

func (b *BaseBar) GetOpenPrice() decimal.Decimal {
	return b.OpenPrice
}

func (b *BaseBar) GetClosePrice() decimal.Decimal {
	return b.ClosePrice
}

func (b *BaseBar) GetHighPrice() decimal.Decimal {
	return b.HighPrice
}

func (b *BaseBar) GetLowPrice() decimal.Decimal {
	return b.LowPrice
}

func (b *BaseBar) GetVolume() decimal.Decimal {
	return *b.Volume
}

func (b *BaseBar) GetAmount() decimal.Decimal {
	return *b.Amount
}

func (b *BaseBar) GetBuyVolume() decimal.Decimal {
	return *b.BuyVolume
}

func (b *BaseBar) GetBuyAmount() decimal.Decimal {
	return *b.BuyAmount
}

func (b *BaseBar) GetTradesCount() decimal.Decimal {
	return *b.TradesCount
}

func (b *BaseBar) IsBullish() bool {
	return b.ClosePrice.Cmp(b.OpenPrice) == 1
}

func (b *BaseBar) IsBearish() bool {
	return b.OpenPrice.Cmp(b.ClosePrice) == 1
}
