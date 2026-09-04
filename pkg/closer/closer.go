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

package closer

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/vandi37/ticket-er/pkg/logger"
	"github.com/vandi37/vanerrors"
)

const (
	ShutdownCancelled = "shutdown cancelled"
	GotSomeErrors     = "got some errors"
)

type Fn func(ctx context.Context) error

type Closer struct {
	mu  sync.Mutex
	fns []Fn
}

func (c *Closer) Add(fn Fn) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.fns = append(c.fns, fn)
}

func (c *Closer) AddNoCtx(fn func() error) {
	c.Add(func(context.Context) error { return fn() })
}
func (c *Closer) AddNoError(fn func()) {
	c.Add(func(context.Context) error { fn(); return nil })
}

func (c *Closer) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var (
		errs = make([]error, 0, len(c.fns))
		wg   sync.WaitGroup
	)

	for _, f := range c.fns {
		wg.Add(1)
		go func(f Fn) {
			defer wg.Done()

			if err := f(ctx); err != nil {
				errs = append(errs, err)
			}

		}(f)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		break
	case <-ctx.Done():
		return vanerrors.Simple(ShutdownCancelled)
	}

	if len(errs) > 0 {
		for _, err := range errs {
			logger.Error(ctx, "Close error", zap.Error(err))
		}
		return vanerrors.New(GotSomeErrors, fmt.Sprint(len(errs)))
	}

	return nil
}

func New() *Closer {
	return &Closer{}
}
