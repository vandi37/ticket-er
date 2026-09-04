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

package finder_test

import (
	"reflect"
	"testing"

	"github.com/vandi37/ticket-er/pkg/finder"
)

type testT struct {
	Value int
}

type testU struct {
	Name string
}

type testContainer struct {
	T      testT
	TP     *testT
	U      testU
	UP     *testU
	Values []testT
	Items  map[string]*testU
}

func TestFindTypes(t *testing.T) {
	t1 := testT{Value: 1}
	t2 := testT{Value: 2}
	u1 := testU{Name: "one"}
	u2 := testU{Name: "two"}

	root := testContainer{
		T:      t1,
		TP:     &t2,
		U:      u1,
		UP:     &u2,
		Values: []testT{{Value: 3}, {Value: 4}},
		Items: map[string]*testU{
			"a": {Name: "three"},
			"b": {Name: "four"},
		},
	}

	gotT, gotU := finder.FindTypes[testT, testU](root)

	wantT := []testT{
		{Value: 1},
		{Value: 2},
		{Value: 3},
		{Value: 4},
	}
	wantU := []testU{
		{Name: "one"},
		{Name: "two"},
		{Name: "three"},
		{Name: "four"},
	}

	if !reflect.DeepEqual(gotT, wantT) {
		t.Errorf("T values = %#v, want %#v", gotT, wantT)
	}
	if !reflect.DeepEqual(gotU, wantU) {
		t.Errorf("U values = %#v, want %#v", gotU, wantU)
	}
}

func TestFindTypes_PointerRoot(t *testing.T) {
	root := &testContainer{
		T: testT{Value: 42},
		U: testU{Name: "answer"},
	}

	gotT, gotU := finder.FindTypes[testT, testU](root)

	wantT := []testT{{Value: 42}}
	wantU := []testU{{Name: "answer"}}

	if !reflect.DeepEqual(gotT, wantT) {
		t.Errorf("T values = %#v, want %#v", gotT, wantT)
	}
	if !reflect.DeepEqual(gotU, wantU) {
		t.Errorf("U values = %#v, want %#v", gotU, wantU)
	}
}

func TestFindTypes_NilPointers(t *testing.T) {
	var tp *testT
	var up *testU

	root := struct {
		T *testT
		U *testU
	}{
		T: tp,
		U: up,
	}

	gotT, gotU := finder.FindTypes[testT, testU](root)

	if len(gotT) != 0 {
		t.Errorf("T values = %#v, want empty", gotT)
	}
	if len(gotU) != 0 {
		t.Errorf("U values = %#v, want empty", gotU)
	}
}

func TestFindTypes_NilRoot(t *testing.T) {
	gotT, gotU := finder.FindTypes[testT, testU](nil)

	if len(gotT) != 0 {
		t.Errorf("T values = %#v, want empty", gotT)
	}
	if len(gotU) != 0 {
		t.Errorf("U values = %#v, want empty", gotU)
	}
}

func TestFindTypes_Interfaces(t *testing.T) {
	var tValue any = testT{Value: 10}
	var tPointer any = &testT{Value: 20}
	var uValue any = testU{Name: "value"}
	var uPointer any = &testU{Name: "pointer"}

	root := []any{
		tValue,
		tPointer,
		uValue,
		uPointer,
	}

	gotT, gotU := finder.FindTypes[testT, testU](root)

	wantT := []testT{
		{Value: 10},
		{Value: 20},
	}
	wantU := []testU{
		{Name: "value"},
		{Name: "pointer"},
	}

	if !reflect.DeepEqual(gotT, wantT) {
		t.Errorf("T values = %#v, want %#v", gotT, wantT)
	}
	if !reflect.DeepEqual(gotU, wantU) {
		t.Errorf("U values = %#v, want %#v", gotU, wantU)
	}
}

func TestFindTypes_LeafSemantics(t *testing.T) {
	type nestedT struct {
		Inner testT
	}

	root := nestedT{
		Inner: testT{Value: 99},
	}

	// testT is matched and treated as a leaf. There is no additional
	// testT inside it in this case, but this test establishes the intended
	// behavior for a matched struct.
	gotT, gotU := finder.FindTypes[testT, testU](root)

	wantT := []testT{{Value: 99}}

	if !reflect.DeepEqual(gotT, wantT) {
		t.Errorf("T values = %#v, want %#v", gotT, wantT)
	}
	if len(gotU) != 0 {
		t.Errorf("U values = %#v, want empty", gotU)
	}
}

func TestFindTypes_MatchedValueIsNotRecursedInto(t *testing.T) {
	type nested struct {
		Inner testT
	}

	// The outer value itself is T, so it must be returned as one T and
	// traversal must stop there.
	root := nested{Inner: testT{Value: 123}}

	gotT, gotU := finder.FindTypes[nested, testU](root)

	wantT := []nested{{Inner: testT{Value: 123}}}

	if !reflect.DeepEqual(gotT, wantT) {
		t.Errorf("T values = %#v, want %#v", gotT, wantT)
	}
	if len(gotU) != 0 {
		t.Errorf("U values = %#v, want empty", gotU)
	}
}

func TestFindTypes_Array(t *testing.T) {
	root := [3]testT{
		{Value: 1},
		{Value: 2},
		{Value: 3},
	}

	gotT, gotU := finder.FindTypes[testT, testU](root)

	wantT := []testT{
		{Value: 1},
		{Value: 2},
		{Value: 3},
	}

	if !reflect.DeepEqual(gotT, wantT) {
		t.Errorf("T values = %#v, want %#v", gotT, wantT)
	}
	if len(gotU) != 0 {
		t.Errorf("U values = %#v, want empty", gotU)
	}
}

func TestFindTypes_Map(t *testing.T) {
	root := map[string]testT{
		"a": {Value: 1},
		"b": {Value: 2},
	}

	gotT, gotU := finder.FindTypes[testT, testU](root)

	// Map iteration order is deliberately unspecified, so compare as sets.
	if len(gotT) != len(root) {
		t.Fatalf("got %d T values, want %d", len(gotT), len(root))
	}

	for _, want := range root {
		found := false
		for _, got := range gotT {
			if reflect.DeepEqual(got, want) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing T value %#v in %#v", want, gotT)
		}
	}

	if len(gotU) != 0 {
		t.Errorf("U values = %#v, want empty", gotU)
	}
}

func TestFindTypes_UnexportedFieldsAreSkipped(t *testing.T) {
	root := struct {
		Public     testT
		unexported testT
	}{
		Public:     testT{Value: 1},
		unexported: testT{Value: 2},
	}

	gotT, gotU := finder.FindTypes[testT, testU](root)

	wantT := []testT{{Value: 1}}

	if !reflect.DeepEqual(gotT, wantT) {
		t.Errorf("T values = %#v, want %#v", gotT, wantT)
	}
	if len(gotU) != 0 {
		t.Errorf("U values = %#v, want empty", gotU)
	}
}

func TestFindTypes_NoMatches(t *testing.T) {
	root := struct {
		Number int
		Text   string
	}{
		Number: 42,
		Text:   "hello",
	}

	gotT, gotU := finder.FindTypes[testT, testU](root)

	if len(gotT) != 0 {
		t.Errorf("T values = %#v, want empty", gotT)
	}
	if len(gotU) != 0 {
		t.Errorf("U values = %#v, want empty", gotU)
	}
}
