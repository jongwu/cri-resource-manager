package topologyaware

import (
	"testing"
)

func TestResolveRequest(t *testing.T) {
	type req_type struct {
		fraction int
	}
	type expect struct {
		req0 req_type
		req1 req_type
	}
	tcs := []struct {
		req  *request
		expected expect
	}{
	{
		req : &request{
			isolate:   false,
			full:      20,
			container: &mockContainer{},
		},
		expected : expect {
			req0 : req_type {
				fraction:	16000,
			},
			req1 : req_type {
				fraction:	4000,
			},
		},
	},
	{
		req : &request{
			isolate:	false,
			full:		24,
			container:	&mockContainer{},
		},
		expected : expect {
			req0 : req_type {
				fraction:	16000,
			},
			req1 : req_type {
				fraction:	8000,
			},
		},
	},
	{
		req : &request{
			isolate:	false,
			full:		32,
			container:	&mockContainer{},
		},
		expected : expect {
			req0 : req_type {
				fraction:	16000,
			},
			req1 : req_type {
				fraction:	16000,
			},
		},
	},
	}
	for _, tc := range tcs {
		req0, req1 := resolveRequest(tc.req, 8, 16)
		var got expect
		if req1 != nil {
			got = expect {
				req0 : req_type {
					fraction:	req0.(*request).fraction,
				},
				req1 : req_type {
					fraction:	req1.(*request).fraction,
				},
			}
		} else {
			got = expect {
				req0 : req_type {
					fraction:	req0.(*request).fraction,
				},
				req1: req_type{},
			}
		}
	
		if req0.(*request).fraction != tc.expected.req0.fraction || tc.expected.req1.fraction != -1 && req1 == nil || req1 != nil && tc.expected.req1.fraction == -1 || req1.(*request).fraction != tc.expected.req1.fraction {
			t.Errorf("expected: req0.fraction: %d, req1.fraction: %d, got: req0.fraction: %d, req1.fraction: %d", tc.expected.req0.fraction, tc.expected.req1.fraction, got.req0.fraction, got.req1.fraction)
		}
	}
}
