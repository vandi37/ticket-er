/**
 * The Ticket-er Telegram Bot Source Code
 * Copyright (C) 2026 Lev (Leo) Kondukov (aka DiceBarrel, Barrel, Vandi)
 * 
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License.
 * 
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 * 
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package waiting

import (
	"cmp"
	"sync"
)

type Waiter[K cmp.Ordered, V any, C comparable] struct {
	mu    sync.RWMutex
	queue map[K]struct {
		ch      chan V
		control *C
		cancel  Cancel
	}
}

func New[K cmp.Ordered, V any, C comparable]() *Waiter[K, V, C] {
	return &Waiter[K, V, C]{
		queue: make(map[K]struct {
			ch      chan V
			control *C
			cancel  Cancel
		}),
	}
}

func (w *Waiter[K, V, C]) add(key K, control *C) (chan V, Cancel) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if old, ok := w.queue[key]; ok {
		old.cancel.Cancel()
		close(old.ch)
	}

	ch := make(chan V)
	cancel := NewCancel()

	w.queue[key] = struct {
		ch      chan V
		control *C
		cancel  Cancel
	}{
		ch:      ch,
		control: control,
		cancel:  cancel,
	}

	return ch, cancel
}

func (w *Waiter[K, V, C]) Add(key K) (chan V, Cancel) {
	return w.add(key, nil)
}

func (w *Waiter[K, V, C]) AddControl(key K, control C) (chan V, Cancel) {
	return w.add(key, &control)
}

func (w *Waiter[K, V, C]) check(key K, val V, control *C) bool {
	w.mu.RLock()
	item, ok := w.queue[key]
	w.mu.RUnlock()

	if !ok || (item.control == nil && control != nil) ||
		(item.control != nil && control == nil) ||
		(item.control != nil && control != nil && *item.control != *control) {
		return false
	}

	item.ch <- val
	return true
}

func (w *Waiter[K, V, C]) Check(key K, val V) bool {
	return w.check(key, val, nil)
}

func (w *Waiter[K, V, C]) CheckControl(key K, val V, control C) bool {
	return w.check(key, val, &control)
}

func (w *Waiter[K, V, C]) Remove(key K) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	item, ok := w.queue[key]
	if !ok {
		return false
	}

	item.cancel.Cancel()
	close(item.ch)
	delete(w.queue, key)

	return true
}
