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

package kline

import "time"

type Interval struct {
	Interval         string
	Seconds          int64
	DataPointsPerDay float32
	Duration         time.Duration
}

var (
	Interval1m  = Interval{"1m", 60, 1440, time.Minute}
	Interval3m  = Interval{"3m", 3 * 60, 480, 3 * time.Minute}
	Interval5m  = Interval{"5m", 5 * 60, 288, 5 * time.Minute}
	Interval15m = Interval{"15m", 15 * 60, 96, 15 * time.Minute}
	Interval30m = Interval{"30m", 30 * 60, 48, 30 * time.Minute}
	Interval1h  = Interval{"1h", 60 * 60, 24, time.Hour}
	Interval2h  = Interval{"2h", 2 * 60 * 60, 12, 2 * time.Hour}
	Interval3h  = Interval{"3h", 3 * 60 * 60, 8, 3 * time.Hour}
	Interval4h  = Interval{"4h", 4 * 60 * 60, 6, 4 * time.Hour}
	Interval6h  = Interval{"6h", 6 * 60 * 60, 4, 6 * time.Hour}
	Interval12h = Interval{"12h", 12 * 60 * 60, 2, 12 * time.Hour}
	Interval1d  = Interval{"1d", 24 * 60 * 60, 1, 24 * time.Hour}
)

func (k Interval) GetInterval() string {
	return k.Interval
}

func (k Interval) GetDuration() time.Duration {
	return k.Duration
}

func (k Interval) ToSeconds() int64 {
	return k.Seconds
}

func (k Interval) ToMillis() int64 {
	return k.ToSeconds() * 1000
}

func (k Interval) GetDataPointsPerDay() float32 {
	return k.DataPointsPerDay
}

func (k Interval) GetDataPointsForDays(days int) float32 {
	return float32(days) * k.GetDataPointsPerDay()
}
