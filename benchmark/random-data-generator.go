// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package benchmark

import (
	"math/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandString returns random string of letters of size n
func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))] // #nosec G404 -- weak random is acceptable for benchmark data generation
	}
	return string(b)
}

// ShuffleSlice shuffles a slice in place
func ShuffleSlice(s []string) {
	rand.Shuffle(
		len(s),
		func(i, j int) { s[i], s[j] = s[j], s[i] },
	) // #nosec G404 -- weak random is acceptable for benchmark data generation
}
