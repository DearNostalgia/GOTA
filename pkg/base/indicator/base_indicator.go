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
	"time"
)

type Result[T any] struct {
	Value     T
	BeginTime time.Time
}

type Indicator[T any] interface {
	Calculate(idx int) (T, error)
	SetSource(s *AttributeSource[T])
	GetSource() *AttributeSource[T]
	GetResults() []*Result[T]
}

type BaseIndicator[T any] struct {
	Source *AttributeSource[T]
}

func NewBasicIndicator[T any]() *BaseIndicator[T] {
	return &BaseIndicator[T]{}
}

func (b *BaseIndicator[T]) SetSource(s *AttributeSource[T]) {
	b.Source = s
}

func (b *BaseIndicator[T]) GetSource() *AttributeSource[T] {
	return b.Source
}
