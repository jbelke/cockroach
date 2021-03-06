// Copyright 2015 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License. See the AUTHORS file
// for names of contributors.
//
// Author: Spencer Kimball (spencer.kimball@gmail.com)

package storage

import (
	"math"
	"testing"

	"github.com/cockroachdb/cockroach/proto"
	"github.com/cockroachdb/cockroach/storage/engine"
	"github.com/cockroachdb/cockroach/util/leaktest"
)

// TestVerifyQueueShouldQueue verifies shouldQueue method correctly
// indicates that a range should be queued for verification if the
// time since last verification exceeds the threshold limit.
func TestVerifyQueueShouldQueue(t *testing.T) {
	defer leaktest.AfterTest(t)
	tc := testContext{}
	tc.Start(t)
	defer tc.Stop()

	// Put empty verification timestamp
	key := engine.RangeLastVerificationTimestampKey(tc.rng.Desc().RaftID)
	if err := engine.MVCCPutProto(tc.rng.rm.Engine(), nil, key, proto.ZeroTimestamp, nil, &proto.Timestamp{}); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		now      proto.Timestamp
		shouldQ  bool
		priority float64
	}{
		// No GC'able bytes, no intent bytes, verification interval elapsed.
		{makeTS(verificationInterval.Nanoseconds(), 0), false, 0},
		// No GC'able bytes, no intent bytes, verification interval * 2 elapsed.
		{makeTS(verificationInterval.Nanoseconds()*2, 0), true, 2},
	}

	verifyQ := newVerifyQueue(nil)

	for i, test := range testCases {
		shouldQ, priority := verifyQ.shouldQueue(test.now, tc.rng)
		if shouldQ != test.shouldQ {
			t.Errorf("%d: should queue expected %t; got %t", i, test.shouldQ, shouldQ)
		}
		if math.Abs(priority-test.priority) > 0.00001 {
			t.Errorf("%d: priority expected %f; got %f", i, test.priority, priority)
		}
	}
}
