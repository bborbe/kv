// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/kv"
	"github.com/bborbe/kv/mocks"
)

var _ = Describe("Count", func() {
	var ctx context.Context
	var bucket *mocks.Bucket
	var iterator *mocks.Iterator

	BeforeEach(func() {
		ctx = context.Background()
		bucket = &mocks.Bucket{}
		iterator = &mocks.Iterator{}
	})

	Context("with empty bucket", func() {
		BeforeEach(func() {
			bucket.IteratorReturns(iterator)
			iterator.ValidReturns(false)
		})

		It("returns count of 0", func() {
			count, err := kv.Count(ctx, bucket)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(int64(0)))
			Expect(iterator.CloseCallCount()).To(Equal(1))
		})
	})

	Context("with items in bucket", func() {
		BeforeEach(func() {
			bucket.IteratorReturns(iterator)
			callCount := 0
			iterator.ValidStub = func() bool {
				callCount++
				return callCount <= 5
			}
		})

		It("returns correct count", func() {
			count, err := kv.Count(ctx, bucket)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(int64(5)))
			Expect(iterator.RewindCallCount()).To(Equal(1))
			Expect(iterator.NextCallCount()).To(Equal(5))
			Expect(iterator.CloseCallCount()).To(Equal(1))
		})
	})

	Context("with context cancellation", func() {
		BeforeEach(func() {
			bucket.IteratorReturns(iterator)
			iterator.ValidReturns(true)
		})

		It("returns -1 and context error when cancelled", func() {
			cancelCtx, cancel := context.WithCancel(ctx)
			cancel()

			count, err := kv.Count(cancelCtx, bucket)
			Expect(err).To(Equal(context.Canceled))
			Expect(count).To(Equal(int64(-1)))
			Expect(iterator.CloseCallCount()).To(Equal(1))
		})
	})
})
