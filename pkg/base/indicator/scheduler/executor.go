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

package schedule

import (
	"github.com/dearnostalgia/gota/pkg/base/indicator"
	"github.com/dearnostalgia/gota/pkg/base/logger"
	clone "github.com/huandu/go-clone/generic"
	"go.uber.org/zap"
)

type Executor[T any] struct{}

func (e *Executor[T]) Shoot(indicator indicator.Indicator[T], startCalIdx int) error {
	source := indicator.GetSource()
	barSeries := source.GetBarSeries()
	eventManager := NewEventProcessor[T](indicator)

	done := make(chan struct{})

	go func() {
		<-done
		l := (*barSeries).Subscribe()
		for event := range l.Ch() {
			err := (*eventManager).ProcessEvent(event)
			if err != nil {
				l.Close()
				logger.Logger.Error("indicator calculator failed to process event, close describe, ",
					zap.Int("event_idx", event.Idx),
					zap.Any("event_type", event.BarSeriesEventType),
					zap.Any("bar_data", *event.Bar),
					zap.Any("indicator", indicator),
					zap.Error(err),
				)
				return
			}
		}
	}()

	if startCalIdx >= 0 {
		copSource := clone.Wrap(source)
		indicator.SetSource(copSource)
		copSize := (*copSource.GetBarSeries()).Size()

		for i := startCalIdx; i < copSize; i++ {
			_, err := indicator.Calculate(i)
			if err != nil {
				return err
			}
		}
		indicator.SetSource(source)
	}

	close(done)
	return nil
}
