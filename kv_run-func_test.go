// Copyright (c) 2025 Benjamin Borbe All rights reserved.
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

var _ = Describe("FuncTx", func() {
	var ctx context.Context
	var tx *mocks.Tx
	var err error
	var callCount int
	var receivedCtx context.Context
	var receivedTx kv.Tx

	BeforeEach(func() {
		ctx = context.Background()
		tx = &mocks.Tx{}
		callCount = 0
		receivedCtx = nil
		receivedTx = nil
	})

	Context("successful execution", func() {
		var funcTx kv.FuncTx

		BeforeEach(func() {
			funcTx = func(ctx context.Context, tx kv.Tx) error {
				callCount++
				receivedCtx = ctx
				receivedTx = tx
				return nil
			}
		})

		JustBeforeEach(func() {
			err = funcTx.Run(ctx, tx)
		})

		It("returns no error", func() {
			Expect(err).To(BeNil())
		})

		It("calls the function once", func() {
			Expect(callCount).To(Equal(1))
		})

		It("passes context and tx correctly", func() {
			Expect(receivedCtx).To(Equal(ctx))
			Expect(receivedTx).To(Equal(tx))
		})
	})

	Context("function with error", func() {
		var funcTx kv.FuncTx
		var testError error

		BeforeEach(func() {
			testError = errors.New("function error")
			funcTx = func(ctx context.Context, tx kv.Tx) error {
				callCount++
				return testError
			}
		})

		JustBeforeEach(func() {
			err = funcTx.Run(ctx, tx)
		})

		It("returns the error", func() {
			Expect(err).To(Equal(testError))
		})

		It("calls the function once", func() {
			Expect(callCount).To(Equal(1))
		})
	})

	Context("function that uses context", func() {
		var funcTx kv.FuncTx
		var cancelCtx context.Context
		var cancel context.CancelFunc

		BeforeEach(func() {
			cancelCtx, cancel = context.WithCancel(context.Background())
			ctx = cancelCtx

			funcTx = func(ctx context.Context, tx kv.Tx) error {
				callCount++
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					return nil
				}
			}

			cancel() // Cancel immediately
		})

		JustBeforeEach(func() {
			err = funcTx.Run(ctx, tx)
		})

		It("returns context cancellation error", func() {
			Expect(err).To(Equal(context.Canceled))
		})

		It("calls the function once", func() {
			Expect(callCount).To(Equal(1))
		})
	})

	Context("function that modifies transaction", func() {
		var funcTx kv.FuncTx
		var bucket *mocks.Bucket

		BeforeEach(func() {
			bucket = &mocks.Bucket{}
			tx.BucketReturns(bucket, nil)

			funcTx = func(ctx context.Context, tx kv.Tx) error {
				callCount++
				_, err := tx.Bucket(ctx, kv.NewBucketName("test"))
				return err
			}
		})

		JustBeforeEach(func() {
			err = funcTx.Run(ctx, tx)
		})

		It("returns no error", func() {
			Expect(err).To(BeNil())
		})

		It("calls the function once", func() {
			Expect(callCount).To(Equal(1))
		})

		It("interacts with transaction", func() {
			Expect(tx.BucketCallCount()).To(Equal(1))
			_, bucketName := tx.BucketArgsForCall(0)
			Expect(bucketName.String()).To(Equal("test"))
		})
	})

	Context("type compatibility", func() {
		It("can be used as RunnableTx", func() {
			var funcTx kv.FuncTx = func(ctx context.Context, tx kv.Tx) error {
				return nil
			}

			var runnable kv.RunnableTx = funcTx
			Expect(runnable).NotTo(BeNil())

			err := runnable.Run(ctx, tx)
			Expect(err).To(BeNil())
		})
	})
})
