// /*
// * MIT License
// *
// * Copyright (c) 2024 DearNostalgia
// *
// * Permission is hereby granted, free of charge, to any person obtaining a copy
// * of this software and associated documentation files (the "Software"), to deal
// * in the Software without restriction, including without limitation the rights
// * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// * copies of the Software, and to permit persons to whom the Software is
// * furnished to do so, subject to the following conditions:
// *
// * The above copyright notice and this permission notice shall be included in all
// * copies or substantial portions of the Software.
// *
// * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// * SOFTWARE.
// *
// */
package bar_series

//
//import (
//	"math"
//	"sync"
//	"time"
//
//	clone "github.com/huandu/go-clone/generic"
//
//	"github.com/google/uuid"
//
//	"github.com/dearnostalgia/gota/pkg/base/bar"
//)
//
//const DefaultInitialSliceCapacity = 1024
//
//var DefaultMaxSize = math.MaxInt32
//
//var _ BarSeries = (*SliceBarSeries)(nil)
//var _ BarSeriesCreator = (*SliceBarSeriesCreator)(nil)
//
//type SliceBarSeries struct {
//	*BarSeriesMetaInfo
//	mu        sync.RWMutex
//	eventChan chan BarSeriesEvent
//	barSlice  *barSlice
//}
//
//type barSlice struct {
//	bars []bar.Bar
//}
//
//type SliceBarSeriesCreator struct{}
//
//func (*SliceBarSeriesCreator) Create(opts ...BaseSeriesOption) BarSeries {
//	sliceBarSeries := &SliceBarSeries{
//		BarSeriesMetaInfo: &BarSeriesMetaInfo{},
//		barSlice:          &barSlice{},
//		mu:                sync.RWMutex{},
//		eventChan:         make(chan BarSeriesEvent, 1),
//	}
//
//	for _, opt := range opts {
//		opt(sliceBarSeries.BarSeriesMetaInfo)
//	}
//
//	if sliceBarSeries.BarSeriesMetaInfo.symbol == nil {
//		symbol := uuid.New().String()
//		sliceBarSeries.BarSeriesMetaInfo.symbol = &symbol
//	}
//
//	maxSize := sliceBarSeries.BarSeriesMetaInfo.maxSize
//
//	if maxSize == nil {
//		sliceBarSeries.barSlice.bars = make([]bar.Bar, 0, DefaultInitialSliceCapacity)
//		sliceBarSeries.maxSize = &DefaultMaxSize
//		return sliceBarSeries
//	}
//
//	if *maxSize <= 0 {
//		sliceBarSeries.barSlice.bars = make([]bar.Bar, 0)
//		return sliceBarSeries
//	}
//
//	sliceBarSeries.barSlice.bars = make([]bar.Bar, 0, *maxSize)
//	return sliceBarSeries
//}
//
//// NewSliceBarSeries creates and initializes a new SliceBarSeries instance.
//// It accepts a variable number of options (BaseSeriesOption) to configure the instance.
//// If no options are provided, the symbol, maxSize, and other fields may remain uninitialized
//// until they are set by the logic within the function.
////
//// This function handles the following scenarios:
//// - If the symbol is not provided, it defaults to a string representation of a new UUID.
//// - If maxSize is not provided, the function sets a default initial slice capacity for the barSlice.
//// - If maxSize is provided and is less than or equal to 0, an empty slice is created.
//// - Otherwise, the slice is initialized with a capacity equal to maxSize.
////
//// Example usage:
////
////	sliceBarSeries := NewSliceBarSeries(
////	    WithSymbol("BTCUSDT.P"),              // Sets the symbol to "BTCUSDT.P"
////	    WithInterval(kline.Interval3m.Duration), // Sets the interval to 3 minutes
////	    WithMaxSize(3060),                    // Sets the maximum size to 3060 bars
////	)
////func NewSliceBarSeries(opts ...BaseSeriesOption) *SliceBarSeries {
////	sliceBarSeries := &SliceBarSeries{
////		BarSeriesMetaInfo: &BarSeriesMetaInfo{},
////		barSlice:          &barSlice{},
////		mu:                sync.RWMutex{},
////		eventChan:         make(chan BarSeriesEvent, 1),
////	}
////
////	for _, opt := range opts {
////		opt(sliceBarSeries.BarSeriesMetaInfo)
////	}
////
////	if sliceBarSeries.BarSeriesMetaInfo.symbol == nil {
////		symbol := uuid.New().String()
////		sliceBarSeries.BarSeriesMetaInfo.symbol = &symbol
////	}
////
////	maxSize := sliceBarSeries.BarSeriesMetaInfo.maxSize
////
////	if maxSize == nil {
////		sliceBarSeries.barSlice.bars = make([]bar.Bar, 0, DefaultInitialSliceCapacity)
////		sliceBarSeries.maxSize = &DefaultMaxSize
////		return sliceBarSeries
////	}
////
////	if *maxSize <= 0 {
////		sliceBarSeries.barSlice.bars = make([]bar.Bar, 0)
////		return sliceBarSeries
////	}
////
////	sliceBarSeries.barSlice.bars = make([]bar.Bar, 0, *maxSize)
////	return sliceBarSeries
////}
//
//func (s *SliceBarSeries) Publish(e *BarSeriesEvent) {
//	s.eventChan <- *e
//}
//
//func (s *SliceBarSeries) ReceiveBarEvent() BarSeriesEvent {
//	return <-s.eventChan
//}
//
//func (s *SliceBarSeries) GetBarSeriesCopy() []bar.Bar {
//	barsCopy := clone.Clone(s.barSlice.bars)
//	return barsCopy
//}
//
//func (s *SliceBarSeries) Size() int {
//	s.mu.RLock()
//	defer s.mu.RUnlock()
//
//	return len(s.barSlice.bars)
//}
//
//func (s *SliceBarSeries) IsEmpty() bool {
//	s.mu.RLock()
//	defer s.mu.RUnlock()
//
//	return len(s.barSlice.bars) == 0
//}
//
//func (s *SliceBarSeries) GetBar(idx int) *bar.Bar {
//	s.mu.RLock()
//	defer s.mu.RUnlock()
//
//	if idx >= len(s.barSlice.bars) || idx < 0 {
//		return nil
//	}
//
//	return &s.barSlice.bars[idx]
//}
//
//func (s *SliceBarSeries) GetFirstBar() *bar.Bar {
//	s.mu.RLock()
//	defer s.mu.RUnlock()
//
//	return &s.barSlice.bars[0]
//}
//
//func (s *SliceBarSeries) GetLastBar() *bar.Bar {
//	s.mu.RLock()
//	defer s.mu.RUnlock()
//
//	return &s.barSlice.bars[len(s.barSlice.bars)-1]
//}
//
//func (s *SliceBarSeries) AddBar(bar bar.Bar, realTimeUpdateBar bool) error {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	maxSize := *s.maxSize
//	if maxSize <= 0 {
//		return nil
//	}
//
//	size := len(s.barSlice.bars)
//	if size == 0 {
//		s.barSlice.bars = append(s.barSlice.bars, bar)
//		//s.Publish(&BarSeriesEvent{
//		//	BarSeriesEventType: AddNewBar,
//		//	idx:                size,
//		//	bar:                &bar,
//		//})
//		return nil
//	}
//	latestBar := s.barSlice.bars[size-1]
//	if realTimeUpdateBar && (latestBar.GetBeginTime().Equal(bar.GetBeginTime()) || !bar.IsEnd()) {
//		s.barSlice.bars[size-1] = bar
//		//s.Publish(&BarSeriesEvent{
//		//	BarSeriesEventType: UpdateLatestBar,
//		//	idx:                size - 1,
//		//	bar:                &bar,
//		//})
//		return nil
//	}
//	if size >= maxSize {
//		remove := size - maxSize + 1
//		s.barSlice.bars = s.barSlice.bars[remove:]
//		//s.Publish(&BarSeriesEvent{
//		//	BarSeriesEventType: RemoveInvalidBar,
//		//	idx:                remove,
//		//})
//	}
//	s.barSlice.bars = append(s.barSlice.bars, bar)
//	//s.Publish(&BarSeriesEvent{
//	//	BarSeriesEventType: AddNewBar,
//	//	idx:                size,
//	//	bar:                &bar,
//	//})
//	return nil
//}
//
//func (b barSlice) binarySearch(targetTime time.Time, size int) int {
//	l, r := 0, size-1
//	for l <= r {
//		mid := l + (r-l)/2
//
//		if (b.bars[mid]).GetBeginTime().Equal(targetTime) {
//			return mid
//		} else if (b.bars[mid]).GetBeginTime().Before(targetTime) {
//			l = mid + 1
//		} else {
//			r = mid - 1
//		}
//	}
//	return l
//}
//
//func (s *SliceBarSeries) Subscribe() *BarSeriesEvent {
//	//TODO implement me
//	panic("implement me")
//}
