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

var _ = Describe("ProviderFunc", func() {
	var ctx context.Context
	var db *mocks.DB
	var err error
	var result kv.DB

	BeforeEach(func() {
		ctx = context.Background()
		db = &mocks.DB{}
	})

	Context("successful provider", func() {
		var provider kv.ProviderFunc

		BeforeEach(func() {
			provider = func(ctx context.Context) (kv.DB, error) {
				return db, nil
			}
		})

		JustBeforeEach(func() {
			result, err = provider.Get(ctx)
		})

		It("returns the database", func() {
			Expect(err).To(BeNil())
			Expect(result).To(Equal(db))
		})
	})

	Context("provider with error", func() {
		var provider kv.ProviderFunc
		var testError error

		BeforeEach(func() {
			testError = errors.New("provider error")
			provider = func(ctx context.Context) (kv.DB, error) {
				return nil, testError
			}
		})

		JustBeforeEach(func() {
			result, err = provider.Get(ctx)
		})

		It("returns the error", func() {
			Expect(err).To(Equal(testError))
			Expect(result).To(BeNil())
		})
	})

	Context("provider that uses context", func() {
		var provider kv.ProviderFunc
		var receivedCtx context.Context

		BeforeEach(func() {
			provider = func(ctx context.Context) (kv.DB, error) {
				receivedCtx = ctx
				return db, nil
			}
		})

		JustBeforeEach(func() {
			result, err = provider.Get(ctx)
		})

		It("passes context correctly", func() {
			Expect(err).To(BeNil())
			Expect(receivedCtx).To(Equal(ctx))
		})
	})

	Context("provider that respects context cancellation", func() {
		var provider kv.ProviderFunc
		var cancelCtx context.Context
		var cancel context.CancelFunc

		BeforeEach(func() {
			cancelCtx, cancel = context.WithCancel(context.Background())
			ctx = cancelCtx

			provider = func(ctx context.Context) (kv.DB, error) {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				default:
					return db, nil
				}
			}

			cancel() // Cancel immediately
		})

		JustBeforeEach(func() {
			result, err = provider.Get(ctx)
		})

		It("returns context cancellation error", func() {
			Expect(err).To(Equal(context.Canceled))
			Expect(result).To(BeNil())
		})
	})
})
