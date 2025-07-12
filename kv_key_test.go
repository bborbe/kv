// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/kv"
)

var _ = Describe("Key", func() {
	Context("String", func() {
		It("returns string representation", func() {
			key := kv.Key([]byte("test-key"))
			Expect(key.String()).To(Equal("test-key"))
		})

		It("handles empty key", func() {
			key := kv.Key([]byte(""))
			Expect(key.String()).To(Equal(""))
		})

		It("handles binary data", func() {
			key := kv.Key([]byte{0x01, 0x02, 0x03})
			Expect(key.String()).To(Equal("\x01\x02\x03"))
		})
	})

	Context("Bytes", func() {
		It("returns byte representation", func() {
			originalBytes := []byte("test-key")
			key := kv.Key(originalBytes)
			Expect(key.Bytes()).To(Equal(originalBytes))
		})

		It("handles empty key", func() {
			key := kv.Key([]byte(""))
			Expect(key.Bytes()).To(Equal([]byte("")))
		})

		It("handles binary data", func() {
			originalBytes := []byte{0x01, 0x02, 0x03}
			key := kv.Key(originalBytes)
			Expect(key.Bytes()).To(Equal(originalBytes))
		})

		It("returns copy of original data", func() {
			originalBytes := []byte("test")
			key := kv.Key(originalBytes)
			returnedBytes := key.Bytes()

			// Modify original slice
			originalBytes[0] = 'X'

			// Returned bytes should reflect the change since Key is just a type alias
			Expect(returnedBytes).To(Equal([]byte("Xest")))
		})
	})
})
