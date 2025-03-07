// Copyright (c) 2022-present, DiceDB contributors
// All rights reserved. Licensed under the BSD 3-Clause License. See LICENSE file in the project root for full license information.

package ironhawk

import (
	"testing"
)

func TestDECR(t *testing.T) {
	client := getLocalConnection()
	defer client.Close()

	testCases := []TestCase{
		{
			name:     "DECR",
			commands: []string{"SET key1 2", "DECR key1", "DECR key1", "DECR key1"},
			expected: []interface{}{"OK", 1, 0, -1},
		},
	}
	runTestcases(t, client, testCases)
}
