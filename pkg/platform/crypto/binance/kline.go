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

package binance

import (
	"github.com/dearnostalgia/gota/pkg/base/bar"
	"github.com/govalues/decimal"
	"time"
)

type KlineHttpResp struct {
	//todo fix websocket
	//for example
	//
	//
	//"k":{
	//    "t":1607443020000,      // Kline start time
	//    "T":1607443079999,      // Kline close time
	//    "i":"1m",               // Interval
	//    "f":116467658886,       // First updateId
	//    "L":116468012423,       // Last updateId
	//    "o":"18787.00",         // Open price
	//    "c":"18804.04",         // Close price
	//    "h":"18804.04",         // High price
	//    "l":"18786.54",         // Low price
	//    "v":"197.664",          // volume
	//    "n": 543,               // Number of trades
	//    "x":false,              // Is this kline closed?
	//    "q":"3715253.19494",    // Quote asset volume
	//    "V":"184.769",          // Taker buy volume
	//    "Q":"3472925.84746",    //Taker buy quote asset volume
	//    "B":"0"                 // Ignore
	//  }
	StartTime  int64   `json:"startTime"`
	EndTime    int64   `json:"endTime"`
	OpenPrice  float64 `json:"openPrice"`
	HighPrice  float64 `json:"highPrice"`
	LowPrice   float64 `json:"lowPrice"`
	ClosePrice float64 `json:"closePrice"`
	Volume     float64 `json:"volume,omitempty"`
	TurnOver   float64 `json:"turnOver,omitempty"`
	IsEnd      bool    `json:"isEnd,omitempty"`
}

func (k *KlineHttpResp) Convert2BaseBar() (KlineTicker, error) {
	var (
		err   error
		kline KlineTicker
	)

	o, err := decimal.NewFromFloat64(k.OpenPrice)
	if err != nil {
		return kline, err
	}
	h, err := decimal.NewFromFloat64(k.HighPrice)
	if err != nil {
		return kline, err
	}
	l, err := decimal.NewFromFloat64(k.LowPrice)
	if err != nil {
		return kline, err
	}
	c, err := decimal.NewFromFloat64(k.ClosePrice)
	if err != nil {
		return kline, err
	}
	v, err := decimal.NewFromFloat64(k.Volume)
	if err != nil {
		return kline, err
	}
	a, err := decimal.NewFromFloat64(k.TurnOver)
	if err != nil {
		return kline, err
	}

	return NewBinanceKline(
		time.UnixMilli(k.StartTime),
		o,
		c,
		h,
		l,
		bar.WithEndTime(time.UnixMilli(k.EndTime)),
		bar.WithVolume(v),
		bar.WithIsEnd(k.IsEnd),
	).WithAmount(a), nil
}
