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

package indicator

import (
	clone "github.com/huandu/go-clone/generic"
	"sync"
)

type BaseIndicatorDecorator[T any] struct {
	Indicator Indicator[T]
	Res       []*Result[T]
	mu        sync.RWMutex
}

func (b *BaseIndicatorDecorator[T]) SetSource(s *AttributeSource[T]) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Indicator.SetSource(s)
}

func (b *BaseIndicatorDecorator[T]) GetSource() *AttributeSource[T] {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Indicator.GetSource()
}

func (b *BaseIndicatorDecorator[T]) Calculate(idx int) (T, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	var zero T
	res, err := b.Indicator.Calculate(idx)
	if err != nil {
		return zero, err
	}
	return res, nil
}

func (b *BaseIndicatorDecorator[T]) GetResults() []*Result[T] {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return clone.Wrap(b.Res)
}

type CacheIndicatorDecorator[T any] struct {
	CacheIndicator CacheIndicator[T]
	Cache          *Cache[T]
	mu             sync.RWMutex
}

func (c *CacheIndicatorDecorator[T]) SetSource(s *AttributeSource[T]) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.CacheIndicator.SetSource(s)
	c.Cache = c.CacheIndicator.GetCache()
}

func (c *CacheIndicatorDecorator[T]) GetSource() *AttributeSource[T] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.CacheIndicator.GetSource()
}

func (c *CacheIndicatorDecorator[T]) Calculate(idx int) (T, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var zero T
	res, err := c.CacheIndicator.Calculate(idx)
	if err != nil {
		return zero, err
	}
	err = c.Cache.store(idx, res)
	if err != nil {
		return zero, err
	}
	return res, nil
}

func (c *CacheIndicatorDecorator[T]) GetResults() []*Result[T] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return clone.Wrap(c.Cache.GetResults())
}
