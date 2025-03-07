// Copyright (c) 2022-present, DiceDB contributors
// All rights reserved. Licensed under the BSD 3-Clause License. See LICENSE file in the project root for full license information.

package ironhawk

import (
	"errors"
	"strconv"
	"testing"
	"time"
)

func TestSET(t *testing.T) {
	client := getLocalConnection()
	defer client.Close()

	expiryTime := strconv.FormatInt(time.Now().Add(1*time.Minute).UnixMilli(), 10)

	testCases := []TestCase{
		{
			name:     "Set and Get Simple Value",
			commands: []string{"SET k v", "GET k"},
			expected: []interface{}{"OK", "v"},
		},
		{
			name:     "Set and Get Integer Value",
			commands: []string{"SET k 123456789", "GET k"},
			expected: []interface{}{"OK", 123456789},
		},
		{
			name:     "Overwrite Existing Key",
			commands: []string{"SET k v1", "SET k 5", "GET k"},
			expected: []interface{}{"OK", "OK", int64(5)},
		},
		{
			name:     "Set with EX option",
			commands: []string{"SET k v EX 2", "GET k", "GET k"},
			expected: []interface{}{"OK", "v", nil},
			delay:    []time.Duration{0, 0, 4 * time.Second},
		},
		{
			name:     "Set with PX option",
			commands: []string{"SET k v PX 2000", "GET k", "GET k"},
			expected: []interface{}{"OK", "v", nil},
			delay:    []time.Duration{0, 0, 4 * time.Second},
		},
		{
			name:     "Set with EX and PX option",
			commands: []string{"SET k v EX 2 PX 2000"},
			expected: []interface{}{
				errors.New("invalid syntax for 'SET' command"),
			},
		},
		{
			name:     "XX on non-existing key",
			commands: []string{"SET k99 v XX", "GET k99"},
			expected: []interface{}{nil, nil},
		},
		{
			name:     "NX on non-existing key",
			commands: []string{"SET k1729 v NX", "GET k1729"},
			expected: []interface{}{"OK", "v"},
		},
		{
			name:     "NX on existing key",
			commands: []string{"SET k1730 v NX", "GET k1730", "SET k1730 v2 NX", "GET k1730"},
			expected: []interface{}{"OK", "v", nil, "v"},
		},
		{
			name:     "PXAT option",
			commands: []string{"SET k v PXAT " + expiryTime, "GET k"},
			expected: []interface{}{"OK", "v"},
		},
		{
			name:     "PXAT option with delete",
			commands: []string{"SET k1 v1 PXAT " + expiryTime, "GET k1", "DEL k1"},
			expected: []interface{}{"OK", "v1", 1},
			delay:    []time.Duration{0, 0, 3 * time.Second},
		},
		{
			name:     "PXAT option with invalid unix time ms",
			commands: []string{"SET k2 v2 PXAT 123123", "GET k2"},
			expected: []interface{}{
				errors.New("invalid value for a parameter in 'SET' command for PXAT parameter"),
				nil,
			},
		},
		{
			name:     "XX on existing key",
			commands: []string{"SET k v1", "SET k v2 XX", "GET k"},
			expected: []interface{}{"OK", "OK", "v2"},
		},
		{
			name:     "Multiple XX operations",
			commands: []string{"SET k v1", "SET k v2 XX", "SET k v3 XX", "GET k"},
			expected: []interface{}{"OK", "OK", "OK", "v3"},
		},
		{
			name:     "EX option",
			commands: []string{"SET k v EX 1", "GET k", "GET k"},
			expected: []interface{}{"OK", "v", nil},
			delay:    []time.Duration{0, 0, 3 * time.Second},
		},
		{
			name:     "XX option",
			commands: []string{"SET k9 v9 XX", "GET k9", "SET k9 v9", "GET k9", "SET k9 v10 XX", "GET k9"},
			expected: []interface{}{nil, nil, "OK", "v9", "OK", "v10"},
		},
		{
			name:     "GET with Existing Value",
			commands: []string{"SET k v", "SET k vv GET"},
			expected: []interface{}{"OK", "v"},
		},
		{
			name:     "GET with Non-Existing Value",
			commands: []string{"SET k10 vv GET"},
			expected: []interface{}{nil},
		},
		{
			commands: []string{"SET k v EX 2", "SET k vv KEEPTTL", "GET k", "SET kk vv", "SET kk vvv KEEPTTL",
				"GET kk", "SET K V EX 2 KEEPTTL",
				"SET K1 vv PX 2000 KEEPTTL",
				"SET K2 vv EXAT " + expiryTime + " KEEPTTL"},
			expected: []interface{}{"OK", "OK", "vv", "OK", "OK", "vvv",
				errors.New("invalid syntax for 'SET' command"),
				errors.New("invalid syntax for 'SET' command"),
				errors.New("invalid syntax for 'SET' command"),
			},
		},
		{
			name:     "SET with no keys or arguments",
			commands: []string{"SET"},
			expected: []interface{}{
				errors.New("wrong number of arguments for 'SET' command"),
			},
		},
	}
	runTestcases(t, client, testCases)
}
