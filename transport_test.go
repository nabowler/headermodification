package headermodification_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/nabowler/headermodification"
)

func TestModify(t *testing.T) {
	// func() http.Header used so that keys will be canonicalized
	type testcase struct {
		name     string
		initial  func() http.Header
		add      func() http.Header
		set      func() http.Header
		expected func() http.Header
	}

	for _, tc := range []testcase{
		{
			name: "add to empty",
			add: func() http.Header {
				h := http.Header{}
				h.Add("one_key", "one_val")
				h.Add("two_key", "two_val")
				return h
			},
			expected: func() http.Header {
				h := http.Header{}
				h.Add("one_key", "one_val")
				h.Add("two_key", "two_val")
				return h
			},
		},
		{
			name: "add to initial",
			initial: func() http.Header {
				h := http.Header{}
				h.Add("zero_key", "zero_val")
				return h
			},
			add: func() http.Header {
				h := http.Header{}
				h.Add("one_key", "one_val")
				h.Add("two_key", "two_val")
				return h
			},
			expected: func() http.Header {
				h := http.Header{}
				h.Add("zero_key", "zero_val")
				h.Add("one_key", "one_val")
				h.Add("two_key", "two_val")
				return h
			},
		},
		{
			// I'm undecided if the this _should_ be the expected behavior,
			// that we should add duplicate values for the same key,
			// but since this is the current behavior,
			// this test verifies this doesn't change without intent
			name: "add to initial with duplicate",
			initial: func() http.Header {
				h := http.Header{}
				h.Add("zero_key", "zero_val")
				return h
			},
			add: func() http.Header {
				h := http.Header{}
				h.Add("zero_key", "zero_val")
				h.Add("zero_key", "zero_val")
				h.Add("one_key", "one_val")
				h.Add("two_key", "two_val")
				return h
			},
			expected: func() http.Header {
				h := http.Header{}
				h.Add("zero_key", "zero_val")
				h.Add("zero_key", "zero_val")
				h.Add("zero_key", "zero_val")
				h.Add("one_key", "one_val")
				h.Add("two_key", "two_val")

				if len(h.Values("zero_key")) != 3 {
					// sanity check for myself that header.Add will continue to add duplicate values
					t.Errorf("test case does not test what I would have expected")
				}
				return h
			},
		},
		{
			name: "set to empty",
			set: func() http.Header {
				h := http.Header{}
				h.Add("one_key", "one_val")
				h.Add("two_key", "two_val")
				return h
			},
			expected: func() http.Header {
				h := http.Header{}
				h.Add("one_key", "one_val")
				h.Add("two_key", "two_val")
				return h
			},
		},
		{
			name: "set to initial",
			initial: func() http.Header {
				h := http.Header{}
				h.Add("zero_key", "zero_val")
				return h
			},
			set: func() http.Header {
				h := http.Header{}
				h.Add("one_key", "one_val")
				h.Add("two_key", "two_val")
				return h
			},
			expected: func() http.Header {
				h := http.Header{}
				h.Add("zero_key", "zero_val")
				h.Add("one_key", "one_val")
				h.Add("two_key", "two_val")
				return h
			},
		},
		{
			name: "set to initial with overwrite",
			initial: func() http.Header {
				h := http.Header{}
				h.Add("zero_key", "zero_val")
				return h
			},
			set: func() http.Header {
				h := http.Header{}
				h.Add("zero_key", "zero_val_overwrite")
				h.Add("one_key", "one_val")
				h.Add("two_key", "two_val")
				return h
			},
			expected: func() http.Header {
				h := http.Header{}
				h.Add("zero_key", "zero_val_overwrite")
				h.Add("one_key", "one_val")
				h.Add("two_key", "two_val")
				return h
			},
		},
		{
			name: "set and add to empty",
			set: func() http.Header {
				h := http.Header{}
				h.Add("one_key", "one_val")
				h.Add("two_key", "two_val")
				return h
			},
			add: func() http.Header {
				h := http.Header{}
				h.Add("three_key", "three_val")
				h.Add("four_key", "four_val")
				return h
			},
			expected: func() http.Header {
				h := http.Header{}
				h.Add("one_key", "one_val")
				h.Add("two_key", "two_val")
				h.Add("three_key", "three_val")
				h.Add("four_key", "four_val")
				return h
			},
		},
		{
			name: "set and add to initial",
			initial: func() http.Header {
				h := http.Header{}
				h.Add("zero_key", "zero_val")
				return h
			},
			set: func() http.Header {
				h := http.Header{}
				h.Add("one_key", "one_val")
				h.Add("two_key", "two_val")
				return h
			},
			add: func() http.Header {
				h := http.Header{}
				h.Add("three_key", "three_val")
				h.Add("four_key", "four_val")
				return h
			},
			expected: func() http.Header {
				h := http.Header{}
				h.Add("zero_key", "zero_val")
				h.Add("one_key", "one_val")
				h.Add("two_key", "two_val")
				h.Add("three_key", "three_val")
				h.Add("four_key", "four_val")
				return h
			},
		},
		{
			name: "set and to initial with overwrite",
			initial: func() http.Header {
				h := http.Header{}
				h.Add("zero_key", "zero_val")
				return h
			},
			set: func() http.Header {
				h := http.Header{}
				h.Add("zero_key", "zero_val_overwrite")
				h.Add("one_key", "one_val")
				h.Add("two_key", "two_val")
				return h
			},
			add: func() http.Header {
				h := http.Header{}
				h.Add("zero_key", "zero_val_added")
				h.Add("three_key", "three_val")
				h.Add("four_key", "four_val")
				return h
			},
			expected: func() http.Header {
				h := http.Header{}
				h.Add("zero_key", "zero_val_overwrite")
				h.Add("zero_key", "zero_val_added")
				h.Add("one_key", "one_val")
				h.Add("two_key", "two_val")
				h.Add("three_key", "three_val")
				h.Add("four_key", "four_val")
				return h
			},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
			if err != nil {
				t.Fatalf(err.Error())
			}
			if tc.initial != nil {
				req.Header = tc.initial()
			}

			add, set, expected := http.Header{}, http.Header{}, http.Header{}
			if tc.add != nil {
				add = tc.add()
			}
			if tc.set != nil {
				set = tc.set()
			}
			if tc.expected != nil {
				expected = tc.expected()
			}

			headermodification.ModifyHeaders(req, set, add)
			actual := req.Header

			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("Expected: %+v Actual: %+v", expected, actual)
			}
		})
	}
}
