// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package benchmark_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"

	libhttp "github.com/bborbe/http"
	logmocks "github.com/bborbe/log/mocks"
	runmocks "github.com/bborbe/run/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/kv/benchmark"
	"github.com/bborbe/kv/mocks"
)

var _ = Describe("Handler", func() {
	var (
		provider         *mocks.Provider
		backgroundRunner *runmocks.FuncRunner
		samplerFactory   *logmocks.LogSamplerFactory
		sampler          *logmocks.LogSampler
		handler          libhttp.WithError
		recorder         *httptest.ResponseRecorder
		req              *http.Request
	)

	BeforeEach(func() {
		provider = &mocks.Provider{}
		backgroundRunner = &runmocks.FuncRunner{}
		samplerFactory = &logmocks.LogSamplerFactory{}
		sampler = &logmocks.LogSampler{}

		samplerFactory.SamplerReturns(sampler)
		sampler.IsSampleReturns(false)
		backgroundRunner.RunReturns(nil)

		handler = benchmark.NewBoltBenchmarkHandler(
			provider,
			backgroundRunner,
			samplerFactory,
		)

		recorder = httptest.NewRecorder()
	})

	Describe("NewBoltBenchmarkHandler", func() {
		It("returns handler", func() {
			Expect(handler).NotTo(BeNil())
		})

		Context("with default parameters", func() {
			BeforeEach(func() {
				var err error
				req, err = http.NewRequest("GET", "/benchmark", nil)
				Expect(err).NotTo(HaveOccurred())
			})

			It("calls background runner", func() {
				ctx := context.Background()
				err := handler.ServeHTTP(ctx, recorder, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(backgroundRunner.RunCallCount()).To(Equal(1))
			})

			It("uses default amount of 1000", func() {
				ctx := context.Background()
				err := handler.ServeHTTP(ctx, recorder, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(recorder.Body.String()).To(ContainSubstring("amount 1000"))
			})

			It("uses default valueLength of 10000", func() {
				ctx := context.Background()
				err := handler.ServeHTTP(ctx, recorder, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(recorder.Body.String()).To(ContainSubstring("valueLength 10000"))
			})

			It("uses default batchSize of 1", func() {
				ctx := context.Background()
				err := handler.ServeHTTP(ctx, recorder, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(recorder.Body.String()).To(ContainSubstring("batchSize 1"))
			})
		})

		Context("with custom parameters", func() {
			BeforeEach(func() {
				var err error
				params := url.Values{}
				params.Set("amount", "500")
				params.Set("valueLength", "5000")
				params.Set("batchSize", "10")
				req, err = http.NewRequest("GET", "/benchmark?"+params.Encode(), nil)
				Expect(err).NotTo(HaveOccurred())
			})

			It("uses custom amount", func() {
				ctx := context.Background()
				err := handler.ServeHTTP(ctx, recorder, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(recorder.Body.String()).To(ContainSubstring("amount 500"))
			})

			It("uses custom valueLength", func() {
				ctx := context.Background()
				err := handler.ServeHTTP(ctx, recorder, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(recorder.Body.String()).To(ContainSubstring("valueLength 5000"))
			})

			It("uses custom batchSize", func() {
				ctx := context.Background()
				err := handler.ServeHTTP(ctx, recorder, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(recorder.Body.String()).To(ContainSubstring("batchSize 10"))
			})
		})

		Context("with invalid parameters", func() {
			BeforeEach(func() {
				var err error
				params := url.Values{}
				params.Set("amount", "invalid")
				params.Set("valueLength", "-1")
				params.Set("batchSize", "0")
				req, err = http.NewRequest("GET", "/benchmark?"+params.Encode(), nil)
				Expect(err).NotTo(HaveOccurred())
			})

			It("falls back to defaults", func() {
				ctx := context.Background()
				err := handler.ServeHTTP(ctx, recorder, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(recorder.Body.String()).To(ContainSubstring("amount 1000"))
				Expect(recorder.Body.String()).To(ContainSubstring("valueLength 10000"))
				Expect(recorder.Body.String()).To(ContainSubstring("batchSize 1"))
			})
		})
	})
})
