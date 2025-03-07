// Copyright (c) 2022-present, DiceDB contributors
// All rights reserved. Licensed under the BSD 3-Clause License. See LICENSE file in the project root for full license information.

package ironhawk

import (
	"testing"
)

func TestDEL(t *testing.T) {
	client := getLocalConnection()
	defer client.Close()

	testCases := []TestCase{
		{
			name:     "DEL with set key",
			commands: []string{"SET k1 v1", "DEL k1", "GET k1"},
			expected: []interface{}{"OK", 1, nil},
		},
		{
			name:     "DEL multiple keys",
			commands: []string{"SET k1 v1", "SET k2 v2", "DEL k1 k2", "GET k1", "GET k2"},
			expected: []interface{}{"OK", "OK", 2, nil, nil},
		},
		{
			name:     "DEL with key not set",
			commands: []string{"GET k3", "DEL k3"},
			expected: []interface{}{nil, 0},
		},
		{
			name:     "DEL with no keys or arguments",
			commands: []string{"DEL"},
			expected: []interface{}{0},
		},
	}
	runTestcases(t, client, testCases)
}
