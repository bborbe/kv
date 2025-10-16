// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package benchmark_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/kv/benchmark"
)

var _ = Describe("RandString", func() {
	It("returns string of correct length", func() {
		result := benchmark.RandString(10)
		Expect(result).To(HaveLen(10))
	})
	It("returns empty string for zero length", func() {
		result := benchmark.RandString(0)
		Expect(result).To(BeEmpty())
	})
	It("returns only letters", func() {
		result := benchmark.RandString(100)
		for _, r := range result {
			isLetter := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
			Expect(isLetter).To(BeTrue(), "character %c should be a letter", r)
		}
	})
	It("returns different strings on multiple calls", func() {
		// With high probability, two random strings should differ
		results := make(map[string]bool)
		for i := 0; i < 10; i++ {
			results[benchmark.RandString(20)] = true
		}
		Expect(results).To(HaveLen(10))
	})
})

var _ = Describe("ShuffleSlice", func() {
	It("preserves all elements", func() {
		original := []string{"a", "b", "c", "d", "e"}
		slice := make([]string, len(original))
		copy(slice, original)

		benchmark.ShuffleSlice(slice)

		Expect(slice).To(ConsistOf(original))
	})
	It("shuffles in place", func() {
		slice := []string{"a", "b", "c", "d", "e"}
		pointer := &slice[0]

		benchmark.ShuffleSlice(slice)

		// Verify it's the same slice (same underlying array)
		Expect(&slice[0]).To(BeIdenticalTo(pointer))
	})
	It("handles empty slice", func() {
		slice := []string{}
		benchmark.ShuffleSlice(slice)
		Expect(slice).To(BeEmpty())
	})
	It("handles single element slice", func() {
		slice := []string{"a"}
		benchmark.ShuffleSlice(slice)
		Expect(slice).To(Equal([]string{"a"}))
	})
	It("shuffles elements with high probability", func() {
		// Run multiple times to verify shuffling occurs
		original := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
		changed := false

		for i := 0; i < 20; i++ {
			slice := make([]string, len(original))
			copy(slice, original)
			benchmark.ShuffleSlice(slice)

			// Check if order changed
			orderChanged := false
			for j := range slice {
				if slice[j] != original[j] {
					orderChanged = true
					break
				}
			}
			if orderChanged {
				changed = true
				break
			}
		}

		Expect(
			changed,
		).To(BeTrue(), "Expected shuffle to change order in at least one of 20 attempts")
	})
})
