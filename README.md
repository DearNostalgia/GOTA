# GOTA

A Comprehensive Go Library for Technical and Financial Analysis.

In the context of financial technical analysis, Golang currently lacks the rich ecosystem found in programming languages like Java, C++, and Python. Therefore, we aim to develop a robust, well-documented, feature-rich, high-performance, and highly customizable and extensible technical analysis library. This will serve as a solid foundation for our team to develop more advanced technical analysis libraries, visualization tools, and a multi-language strategy platform in the future.
## Getting Started

```
go get -u github.com/dearnostalgia/gota
```

## Usage

### Bar

One possible implementation

```go
b := bar.NewBaseBar(
                time.UnixMilli(k.StartTime),
		openPrice,
		closePrice,
		highPrice,
		lowPrice,
		bar.WithEndTime(time.UnixMilli(k.EndTime)),
		bar.WithAmount(a),
		bar.WithVolume(v),
		true,
	)
b.GetBeginTime()
b.GetHighPrice()
b.GetVolume()
...
```

#### interface

Implement the basic interface: `Bar`, which allows you to freely customize your own candlestick (K-line).

```go
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
```

### BarSeries

`CircularBarSeries` is our primary implementation.

If WithMaxSize is used, the oldest data will be automatically deleted when the capacity limit is reached.

```go
var barSeries *bar_series.CircularBarSeries
// var barSeries bar_series.BarSeries
barSeries = bar_series.NewCircularBarSeries(
		bar_series.WithSymbol("BTCUSDT.P"),
		bar_series.WithInterval(kline.Interval3m.Duration),
		bar_series.WithMaxSize(512),
)
// All of the above parameters are optional. If WithMaxSize is not configured, please ensure sufficient memory is available.

barSeries.GetBar(9)
barSeries.GetFirstBar()
barSeries.GetLastBar()
barSeries.GetBarsCopy()
barSeries.Size()
barSeries.IsEmpty()
barSeries.GetSymbol()
barSeries.GetMaxSize()
l := barSeries.Subscribe()
for event := range l.Ch() {
    // handle event
}
	
```

#### interface

Implement the basic interface: `BarSeries`, which allows you to freely customize your own `K-line` series.

```go
type BarSeries interface {
	// Size returns the current number of bars in the series.
	// It indicates how many data points (bars) are currently stored in the series.
	Size() int

	// IsEmpty checks if the bar series is empty.
	// It returns true if the series contains no bars, and false otherwise.
	IsEmpty() bool

	// GetBar returns the Bar at the specified index.
	// The index is zero-based, meaning that the first bar in the series is at index 0.
	// If the index is out of range, it returns nil.
	GetBar(idx int) *bar.Bar

	// GetFirstBar returns the first Bar in the series.
	// If the series is empty, it returns nil.
	// The first bar is the one that was added to the series first, based on time.
	GetFirstBar() *bar.Bar

	// GetLastBar returns the last Bar in the series.
	// If the series is empty, it returns nil.
	// The last bar is the most recently added bar in the series.
	GetLastBar() *bar.Bar

	// AddBar adds a bar to the series.
	// If realTimeUpdateBar is true, the new data will be compared with the latest data in the bar series.
	// If the beginTime is the same, the existing data will be updated; if it is different, the new data will be added.
	// If realTimeUpdateBar is false, new data will be continuously added.
	AddBar(bar bar.Bar, realTimeUpdateBar bool) error

	// GetBarsCopy returns a deep copy of the bar.Bar slice.
	// This method ensures that the returned slice contains independent copies of
	// the elements from the original barSlice, so any modifications made to the
	// returned slice will not affect the original BarSeries structure or its data.
	GetBarsCopy() []bar.Bar

	// Subscribe registers a subscription to the barSeries data and returns a contact.Listener[BarSeriesEvent].
	// This method allows the subscriber to listen for broadcasted data events within the CircularBarSeries.
	// Usage:
	//     l := barSeries.Subscribe()
	//     for event := range l.Ch() {
	//         // handle event
	//     }
	// The CircularBarSeries broadcasts data to all subscribed goroutines.
	Subscribe() *contact.Listener[BarSeriesEvent]
}
```

### DataSource

DataSource serves as the calculation parameter for the Indicator. We have implemented the basic data source information, such as `HighPrice`, `ClosePrice`, etc.

```go
closeSource := indicator.NewAttributeSource[decimal.Decimal](barSeries, &indicator.ClosePriceStrategy{})
highSource := indicator.NewAttributeSource[decimal.Decimal](barSeries, &indicator.HighPriceStrategy{})
volumeSource := indicator.NewAttributeSource[decimal.Decimal](barSeries, &indicator.VolumeStrategy{})
```

If you want to customize your own data source, refer to the strategy pattern below, and ensure you design your custom Bar.

```go
type ClosePriceStrategy struct{}

func (s *ClosePriceStrategy) GetAttribute(bar *bar.Bar) decimal.Decimal {
	return (*bar).GetClosePrice()
}

type CustomeStrategy struct{}

func (s *CustomeStrategy) GetAttribute(bar *bar.Bar) decimal.Decimal {
	return (*bar).GetCustomeValue()
}

```

### Indicator

Using the factory pattern, `indicatorExecutor.Shoot` not only performs the initial calculation but also subscribes to `BarSeries` updates to maintain synchronized calculations for the EMAIndicator.

It is strongly recommended to use `IndicatorFactory.CreateIndicator` to generate any indicator. We have utilized the decorator pattern to encapsulate operations such as caching and concurrency protection.

If the barSeries is pre-populated (i.e., fixed before the initialization of the Executor), calculations can be performed directly by iteration, without the need to reserve additional time.

Otherwise, after invoking `indicatorExecutor.Shoot`, it is advisable to wait for approximately 50 milliseconds (depending on individual machine performance) to allow sufficient time for initialization and subscription.
During the `addBar` process, it is also recommended to reserve some time for broadcasting events and performing calculations.


```go
IndicatorFactory := new(factory.IndicatorFactory[decimal.Decimal])
ema, err := tech.NewEMAIndicator(20)
if err != nil {
	return err
}

closeSource := indicator.NewAttributeSource[decimal.Decimal](barSeries, &indicator.ClosePriceStrategy{})
EMAIndicator := IndicatorFactory.CreateIndicator(ema, closeSource)
indicatorExecutor := &schedule.Executor[decimal.Decimal]{}

err := indicatorExecutor.Shoot(EMAIndicator, 0)
	if err != nil {
		return err
}
time.Sleep(50 * time.Microsecond)

for _, bar := range bars {
	err := barSeries.AddBar(bar, true)
	if err != nil {
		return err
	}
	time.Sleep(30 * time.Microsecond)
}
res := EMAIndicator.GetResults()
```

If you want to retrieve the value of a specific calculation

```go
v, err := EMAIndicator.Calculate(idx)
```

#### interface

To implement your own Indicator, you need to implement the Indicator interface.

```go
type Indicator[T any] interface {
	Calculate(idx int) (T, error)
	SetSource(s *AttributeSource[T])
	GetSource() *AttributeSource[T]
	GetResults() []*Result[T]
}
```

If the indicator relies on previous values like `EMA`, it is recommended to reuse the `CacheIndicator` functionality. You can refer to the EMA implementation.

```go
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

	defer func() {
		if err != nil {
			common.IndicatorCalculateErrLog(idx, common.EMA, err)
		}
	}()

	prevEMA, err := e.Cache.GetValue(idx-1, e.Calculate)

	if err != nil {
		return decimal.Zero, err
	}

	res, err := e.calculateEMA(prevEMA, sv)
	if err != nil {
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
```

