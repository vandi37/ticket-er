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

package finder

import (
	"reflect"
)

// FindTypes walks root (a struct, a pointer to one, or anything containing
// one) and returns two slices:
//   - ts: every value found of type T or *T (pointers are dereferenced)
//   - us: every value found of type U or *U (pointers are dereferenced)
//
// A matched T/U value is treated as a leaf: FindTypes does not recurse
// further into it, even if T or U is itself a struct. If you need to also
// search *inside* matched values, remove the early `return` in the match
// cases below.
//
// Written with Claude
func FindTypes[T any, U any](root any) (ts []T, us []U) {
	tType := reflect.TypeOf((*T)(nil)).Elem()
	uType := reflect.TypeOf((*U)(nil)).Elem()
	tPtrType := reflect.PointerTo(tType)
	uPtrType := reflect.PointerTo(uType)

	walk(reflect.ValueOf(root), tType, uType, tPtrType, uPtrType, &ts, &us)
	return ts, us
}

func walk[T any, U any](
	v reflect.Value,
	tType, uType, tPtrType, uPtrType reflect.Type,
	ts *[]T, us *[]U,
) {
	if !v.IsValid() {
		return
	}

	// Unwrap interface values to get at the concrete value inside.
	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	vt := v.Type()
	switch {
	case vt == tType:
		*ts = append(*ts, v.Interface().(T))
		return
	case vt == uType:
		*us = append(*us, v.Interface().(U))
		return
	case vt == tPtrType:
		if !v.IsNil() {
			*ts = append(*ts, v.Elem().Interface().(T))
		}
		return
	case vt == uPtrType:
		if !v.IsNil() {
			*us = append(*us, v.Elem().Interface().(U))
		}
		return
	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return
		}
		walk(v.Elem(), tType, uType, tPtrType, uPtrType, ts, us)

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			if !f.CanInterface() {
				// Unexported field; skip (can't be read via reflection
				// without unsafe).
				continue
			}
			walk(f, tType, uType, tPtrType, uPtrType, ts, us)
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			walk(v.Index(i), tType, uType, tPtrType, uPtrType, ts, us)
		}

	case reflect.Map:
		for _, k := range v.MapKeys() {
			walk(v.MapIndex(k), tType, uType, tPtrType, uPtrType, ts, us)
		}
	}
	// Other kinds (int, string, chan, func, ...) that didn't match T/U are
	// leaves we can't/don't need to descend into.
}
