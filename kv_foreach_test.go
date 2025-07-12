// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/kv"
	"github.com/bborbe/kv/mocks"
)

var _ = Describe("ForEach", func() {
	var ctx context.Context
	var bucket *mocks.Bucket
	var iterator *mocks.Iterator
	var item *mocks.Item

	BeforeEach(func() {
		ctx = context.Background()
		bucket = &mocks.Bucket{}
		iterator = &mocks.Iterator{}
		item = &mocks.Item{}
	})

	Context("with empty bucket", func() {
		BeforeEach(func() {
			bucket.IteratorReturns(iterator)
			iterator.ValidReturns(false)
		})

		It("executes successfully", func() {
			callCount := 0
			err := kv.ForEach(ctx, bucket, func(item kv.Item) error {
				callCount++
				return nil
			})
			Expect(err).To(BeNil())
			Expect(callCount).To(Equal(0))
			Expect(iterator.CloseCallCount()).To(Equal(1))
		})
	})

	Context("with items in bucket", func() {
		BeforeEach(func() {
			bucket.IteratorReturns(iterator)
			callCount := 0
			iterator.ValidStub = func() bool {
				callCount++
				return callCount <= 3
			}
			iterator.ItemReturns(item)
		})

		It("iterates over all items", func() {
			callCount := 0
			err := kv.ForEach(ctx, bucket, func(item kv.Item) error {
				callCount++
				return nil
			})
			Expect(err).To(BeNil())
			Expect(callCount).To(Equal(3))
			Expect(iterator.RewindCallCount()).To(Equal(1))
			Expect(iterator.NextCallCount()).To(Equal(3))
			Expect(iterator.CloseCallCount()).To(Equal(1))
		})

		It("returns error when function fails", func() {
			expectedErr := errors.New("test error")
			err := kv.ForEach(ctx, bucket, func(item kv.Item) error {
				return expectedErr
			})
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("fn failed"))
			Expect(iterator.CloseCallCount()).To(Equal(1))
		})
	})

	Context("with context cancellation", func() {
		BeforeEach(func() {
			bucket.IteratorReturns(iterator)
			iterator.ValidReturns(true)
			iterator.ItemReturns(item)
		})

		It("returns context error when cancelled", func() {
			cancelCtx, cancel := context.WithCancel(ctx)
			cancel()

			err := kv.ForEach(cancelCtx, bucket, func(item kv.Item) error {
				return nil
			})
			Expect(err).To(Equal(context.Canceled))
			Expect(iterator.CloseCallCount()).To(Equal(1))
		})
	})
})
