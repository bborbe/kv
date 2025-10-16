// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package benchmark_test

import (
	"context"

	logmocks "github.com/bborbe/log/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/kv"
	"github.com/bborbe/kv/benchmark"
	"github.com/bborbe/kv/mocks"
)

var _ = Describe("Benchmark", func() {
	var (
		ctx               context.Context
		db                *mocks.DB
		samplerFactory    *logmocks.LogSamplerFactory
		sampler           *logmocks.LogSampler
		benchmarkInstance benchmark.Benchmark
	)

	BeforeEach(func() {
		ctx = context.Background()
		db = &mocks.DB{}
		samplerFactory = &logmocks.LogSamplerFactory{}
		sampler = &logmocks.LogSampler{}

		samplerFactory.SamplerReturns(sampler)
		sampler.IsSampleReturns(false)

		benchmarkInstance = benchmark.NewBenchmark(db, samplerFactory)
	})

	Describe("NewBenchmark", func() {
		It("returns benchmark instance", func() {
			Expect(benchmarkInstance).NotTo(BeNil())
		})
	})

	Describe("Benchmark", func() {
		var (
			tx     *mocks.Tx
			bucket *mocks.Bucket
			item   *mocks.Item
		)

		BeforeEach(func() {
			tx = &mocks.Tx{}
			bucket = &mocks.Bucket{}
			item = &mocks.Item{}

			// Setup DB to call the transaction function
			db.UpdateStub = func(ctx context.Context, fn func(context.Context, kv.Tx) error) error {
				return fn(ctx, tx)
			}
			db.ViewStub = func(ctx context.Context, fn func(context.Context, kv.Tx) error) error {
				return fn(ctx, tx)
			}

			// Setup tx to return bucket
			tx.CreateBucketIfNotExistsReturns(bucket, nil)
			tx.BucketReturns(bucket, nil)

			// Setup bucket operations
			bucket.PutReturns(nil)
			bucket.GetReturns(item, nil)

			// Setup item to call value function with non-empty data
			item.ValueStub = func(fn func([]byte) error) error {
				return fn([]byte("test-value"))
			}
		})

		It("returns result with correct amount and value size", func() {
			result, err := benchmarkInstance.Benchmark(ctx, 10, 100, 1)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
			Expect(result.Amount).To(Equal(10))
			Expect(result.ValueSize).To(Equal(100))
		})

		It("returns result with write and read durations", func() {
			result, err := benchmarkInstance.Benchmark(ctx, 5, 50, 1)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.WriteDuration).To(BeNumerically(">", 0))
			Expect(result.ReadDuration).To(BeNumerically(">", 0))
		})

		It("calls Update for writes", func() {
			_, err := benchmarkInstance.Benchmark(ctx, 10, 100, 1)
			Expect(err).NotTo(HaveOccurred())
			Expect(db.UpdateCallCount()).To(BeNumerically(">=", 10))
		})

		It("calls View for reads", func() {
			_, err := benchmarkInstance.Benchmark(ctx, 10, 100, 1)
			Expect(err).NotTo(HaveOccurred())
			Expect(db.ViewCallCount()).To(BeNumerically(">=", 10))
		})

		It("uses batch size for operations", func() {
			_, err := benchmarkInstance.Benchmark(ctx, 10, 100, 5)
			Expect(err).NotTo(HaveOccurred())
			// With batch size 5 and amount 10, we expect 2 transactions
			Expect(db.UpdateCallCount()).To(Equal(2))
			Expect(db.ViewCallCount()).To(Equal(2))
		})

		It("writes and reads all keys", func() {
			amount := 10
			_, err := benchmarkInstance.Benchmark(ctx, amount, 100, 1)
			Expect(err).NotTo(HaveOccurred())
			Expect(bucket.PutCallCount()).To(Equal(amount))
			Expect(bucket.GetCallCount()).To(Equal(amount))
		})

		Context("when context is cancelled", func() {
			BeforeEach(func() {
				// Cancel context after first call
				callCount := 0
				db.UpdateStub = func(ctx context.Context, fn func(context.Context, kv.Tx) error) error {
					callCount++
					if callCount > 1 {
						return context.Canceled
					}
					return fn(ctx, tx)
				}
			})

			It("returns error", func() {
				_, err := benchmarkInstance.Benchmark(ctx, 100, 100, 1)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
