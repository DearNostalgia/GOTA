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

package contact

import (
	"context"
	"sync"
)

// Relay is the struct in charge of handling the listeners and dispatching the notifications.
type Relay[T any] struct {
	mu      sync.RWMutex
	n       uint32
	clients map[uint32]*Listener[T]
}

// Listener is a Relay listener.
type Listener[T any] struct {
	ch    chan T
	id    uint32
	relay *Relay[T]
	once  sync.Once
}

// NewRelay is the factory to create a Relay.
func NewRelay[T any]() *Relay[T] {
	return &Relay[T]{
		clients: make(map[uint32]*Listener[T]),
	}
}

// Notify sends a notification to all the listeners.
// It guarantees that all the listeners will receive the notification.
func (r *Relay[T]) Notify(v T) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, client := range r.clients {
		client.ch <- v
	}
}

// NotifyCtx tries sending a notification to all the listeners until the context times out or is canceled.
func (r *Relay[T]) NotifyCtx(ctx context.Context, v T) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, client := range r.clients {
		select {
		case client.ch <- v:
		case <-ctx.Done():
			return
		}
	}
}

// Broadcast broadcasts a notification to all the listeners.
// The notification is sent in a non-blocking manner, so there's no guarantee that a listener receives it.
func (r *Relay[T]) Broadcast(v T) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, client := range r.clients {
		select {
		case client.ch <- v:
		default:
		}
	}
}

// Listener creates a new listener given a channel capacity.
func (r *Relay[T]) Listener(capacity int) *Listener[T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	listener := &Listener[T]{
		ch:    make(chan T, capacity),
		id:    r.n,
		relay: r,
	}
	r.clients[r.n] = listener
	r.n++
	return listener
}

// Close closes a relay.
// This operation can be safely called in the meantime as Listener.Close()
func (r *Relay[T]) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, client := range r.clients {
		r.closeRelay(client)
	}
	r.clients = nil
}

func (r *Relay[T]) closeRelay(l *Listener[T]) {
	l.once.Do(func() {
		close(l.ch)
		delete(r.clients, l.id)
	})
}

func (r *Relay[T]) closeListener(l *Listener[T]) {
	r.mu.Lock()
	defer r.mu.Unlock()

	close(l.ch)
	delete(r.clients, l.id)
}

// Ch returns the Listener channel.
func (l *Listener[T]) Ch() <-chan T {
	return l.ch
}

// Close closes a listener.
// This operation can be safely called in the meantime as Relay.Close()
func (l *Listener[T]) Close() {
	l.once.Do(func() {
		l.relay.closeListener(l)
	})
}
