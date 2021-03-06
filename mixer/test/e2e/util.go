// Copyright 2016 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"istio.io/istio/mixer/pkg/adapter"
	"istio.io/istio/mixer/test/spyAdapter"
)

// ConstructAdapterInfos constructs spyAdapters for each of the adptBehavior. It returns
// the constructed spyAdapters along with the adapters Info functions.
func ConstructAdapterInfos(adptBehaviors []spyAdapter.AdapterBehavior) ([]adapter.InfoFn, []*spyAdapter.Adapter) {
	var adapterInfos []adapter.InfoFn = make([]adapter.InfoFn, 0)
	spyAdapters := make([]*spyAdapter.Adapter, 0)
	for _, b := range adptBehaviors {
		sa := spyAdapter.NewSpyAdapter(b)
		spyAdapters = append(spyAdapters, sa)
		adapterInfos = append(adapterInfos, sa.GetAdptInfoFn())
	}
	return adapterInfos, spyAdapters
}

// CmpSliceAndErr compares two slices
func CmpSliceAndErr(t *testing.T, msg string, act, exp interface{}) {
	a := interfaceSlice(exp)
	b := interfaceSlice(act)
	if len(a) != len(b) {
		t.Errorf(fmt.Sprintf("Not equal -> %s.\nActual :\n%s\n\nExpected :\n%s", msg, spew.Sdump(act), spew.Sdump(exp)))
		return
	}

	for _, x1 := range a {
		f := false
		for _, x2 := range b {
			if reflect.DeepEqual(x1, x2) {
				f = true
			}
		}
		if !f {
			t.Errorf(fmt.Sprintf("Not equal -> %s.\nActual :\n%s\n\nExpected :\n%s", msg, spew.Sdump(act), spew.Sdump(exp)))
			return
		}
	}
}

// CmpMapAndErr compares two maps
func CmpMapAndErr(t *testing.T, msg string, act, exp interface{}) {
	want := interfaceMap(exp)
	got := interfaceMap(act)
	if len(want) != len(got) {
		t.Errorf(fmt.Sprintf("Not equal -> %s.\nActual :\n%s\n\nExpected :\n%s", msg, spew.Sdump(act), spew.Sdump(exp)))
		return
	}

	for wk, wv := range want {
		if v, found := got[wk]; !found {
			t.Errorf(fmt.Sprintf("Not equal -> %s.\nActual :\n%s\n\nExpected :\n%s", msg, spew.Sdump(act), spew.Sdump(exp)))
		} else {
			if !reflect.DeepEqual(wv, v) {
				t.Errorf(fmt.Sprintf("Not equal -> %s.\nActual :\n%s\n\nExpected :\n%s", msg, spew.Sdump(act), spew.Sdump(exp)))
			}
		}
	}
}

func interfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)

	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func interfaceMap(m interface{}) map[interface{}]interface{} {
	s := reflect.ValueOf(m)

	ret := make(map[interface{}]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		k := s.MapKeys()[i]
		v := s.MapIndex(k)
		ret[k.Interface()] = v.Interface()
	}

	return ret
}
