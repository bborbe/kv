// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/kv"
)

var _ = Describe("ByteItem", func() {
	Context("NewByteItem", func() {
		It("creates new byte item", func() {
			key := []byte("test-key")
			value := []byte("test-value")
			item := kv.NewByteItem(key, value)
			Expect(item).NotTo(BeNil())
		})
	})

	Context("Exists", func() {
		It("returns true for non-empty value", func() {
			item := kv.NewByteItem([]byte("key"), []byte("value"))
			Expect(item.Exists()).To(BeTrue())
		})

		It("returns false for empty value", func() {
			item := kv.NewByteItem([]byte("key"), []byte(""))
			Expect(item.Exists()).To(BeFalse())
		})

		It("returns false for nil value", func() {
			item := kv.NewByteItem([]byte("key"), nil)
			Expect(item.Exists()).To(BeFalse())
		})
	})

	Context("Key", func() {
		It("returns the key", func() {
			key := []byte("test-key")
			item := kv.NewByteItem(key, []byte("value"))
			Expect(item.Key()).To(Equal(key))
		})

		It("handles empty key", func() {
			key := []byte("")
			item := kv.NewByteItem(key, []byte("value"))
			Expect(item.Key()).To(Equal(key))
		})

		It("handles nil key", func() {
			item := kv.NewByteItem(nil, []byte("value"))
			Expect(item.Key()).To(BeNil())
		})
	})

	Context("Value", func() {
		It("calls function with value", func() {
			value := []byte("test-value")
			item := kv.NewByteItem([]byte("key"), value)

			var receivedValue []byte
			err := item.Value(func(val []byte) error {
				receivedValue = val
				return nil
			})

			Expect(err).To(BeNil())
			Expect(receivedValue).To(Equal(value))
		})

		It("handles empty value", func() {
			item := kv.NewByteItem([]byte("key"), []byte(""))

			var receivedValue []byte
			err := item.Value(func(val []byte) error {
				receivedValue = val
				return nil
			})

			Expect(err).To(BeNil())
			Expect(receivedValue).To(Equal([]byte("")))
		})

		It("handles nil value", func() {
			item := kv.NewByteItem([]byte("key"), nil)

			var receivedValue []byte
			err := item.Value(func(val []byte) error {
				receivedValue = val
				return nil
			})

			Expect(err).To(BeNil())
			Expect(receivedValue).To(BeNil())
		})

		It("propagates function error", func() {
			item := kv.NewByteItem([]byte("key"), []byte("value"))
			testError := errors.New("test error")

			err := item.Value(func(val []byte) error {
				return testError
			})

			Expect(err).To(Equal(testError))
		})

		It("allows function to modify value without affecting item", func() {
			originalValue := []byte("test-value")
			item := kv.NewByteItem([]byte("key"), originalValue)

			err := item.Value(func(val []byte) error {
				if len(val) > 0 {
					val[0] = 'X'
				}
				return nil
			})

			Expect(err).To(BeNil())

			// Verify original value was modified (since it's the same slice)
			var receivedValue []byte
			err = item.Value(func(val []byte) error {
				receivedValue = val
				return nil
			})

			Expect(err).To(BeNil())
			Expect(receivedValue[0]).To(Equal(byte('X')))
		})
	})
})
