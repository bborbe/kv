// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/kv"
	"github.com/bborbe/kv/mocks"
)

var _ = Describe("Reset Handlers", func() {
	var db *mocks.DB
	var cancel context.CancelFunc
	var cancelCalled bool

	BeforeEach(func() {
		db = &mocks.DB{}
		_, cancel = context.WithCancel(context.Background())
		cancelCalled = false

		originalCancel := cancel
		cancel = func() {
			cancelCalled = true
			originalCancel()
		}
	})

	Describe("NewResetHandler", func() {
		var handler http.Handler
		var recorder *httptest.ResponseRecorder
		var request *http.Request

		BeforeEach(func() {
			handler = kv.NewResetHandler(db, cancel)
			recorder = httptest.NewRecorder()
			request = httptest.NewRequest("POST", "/reset", nil)
		})

		Context("when reset is successful", func() {
			BeforeEach(func() {
				db.CloseReturns(nil)
				db.RemoveReturns(nil)
			})

			It("closes and removes database successfully", func() {
				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(recorder.Body.String()).To(ContainSubstring("reset db successful"))
				Expect(db.CloseCallCount()).To(Equal(1))
				Expect(db.RemoveCallCount()).To(Equal(1))
				Expect(cancelCalled).To(BeTrue())
			})
		})

		Context("when close fails", func() {
			BeforeEach(func() {
				db.CloseReturns(errors.New("close failed"))
			})

			It("returns error and does not remove database", func() {
				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				Expect(recorder.Body.String()).To(ContainSubstring("reset db failed"))
				Expect(db.CloseCallCount()).To(Equal(1))
				Expect(db.RemoveCallCount()).To(Equal(0))
				Expect(cancelCalled).To(BeTrue())
			})
		})

		Context("when remove fails", func() {
			BeforeEach(func() {
				db.CloseReturns(nil)
				db.RemoveReturns(errors.New("remove failed"))
			})

			It("returns error after close", func() {
				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				Expect(recorder.Body.String()).To(ContainSubstring("remove db failed"))
				Expect(db.CloseCallCount()).To(Equal(1))
				Expect(db.RemoveCallCount()).To(Equal(1))
				Expect(cancelCalled).To(BeTrue())
			})
		})

		Context("when reset is already running", func() {
			It("prevents concurrent execution", func() {
				db.CloseStub = func() error {
					// Simulate a second request while first is running
					secondRecorder := httptest.NewRecorder()
					secondRequest := httptest.NewRequest("POST", "/reset", nil)
					handler.ServeHTTP(secondRecorder, secondRequest)

					Expect(secondRecorder.Code).To(Equal(http.StatusInternalServerError))
					Expect(secondRecorder.Body.String()).To(ContainSubstring("reset db already running"))
					return nil
				}
				db.RemoveReturns(nil)

				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(db.CloseCallCount()).To(Equal(1))
				Expect(db.RemoveCallCount()).To(Equal(1))
			})
		})
	})

	Describe("NewResetBucketHandler", func() {
		var handler http.Handler
		var recorder *httptest.ResponseRecorder
		var request *http.Request
		var tx *mocks.Tx

		BeforeEach(func() {
			handler = kv.NewResetBucketHandler(db, cancel)
			recorder = httptest.NewRecorder()
			tx = &mocks.Tx{}
		})

		Context("when bucket reset is successful", func() {
			BeforeEach(func() {
				request = httptest.NewRequest("DELETE", "/bucket/test-bucket", nil)
				request = mux.SetURLVars(request, map[string]string{"BucketName": "test-bucket"})

				db.UpdateStub = func(ctx context.Context, fn func(context.Context, kv.Tx) error) error {
					return fn(ctx, tx)
				}
				tx.DeleteBucketReturns(nil)
			})

			It("deletes bucket successfully", func() {
				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(recorder.Body.String()).To(ContainSubstring("reset bucket test-bucket successful"))
				Expect(db.UpdateCallCount()).To(Equal(1))
				Expect(tx.DeleteBucketCallCount()).To(Equal(1))
				_, actualBucketName := tx.DeleteBucketArgsForCall(0)
				Expect(string(actualBucketName)).To(Equal("test-bucket"))
				Expect(cancelCalled).To(BeTrue())
			})
		})

		Context("when bucket name is missing", func() {
			BeforeEach(func() {
				request = httptest.NewRequest("DELETE", "/bucket/", nil)
				request = mux.SetURLVars(request, map[string]string{})
			})

			It("returns bad request error", func() {
				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				Expect(recorder.Body.String()).To(ContainSubstring("parameter bucket missing"))
				Expect(db.UpdateCallCount()).To(Equal(0))
				Expect(cancelCalled).To(BeFalse())
			})
		})

		Context("when delete bucket fails", func() {
			BeforeEach(func() {
				request = httptest.NewRequest("DELETE", "/bucket/test-bucket", nil)
				request = mux.SetURLVars(request, map[string]string{"BucketName": "test-bucket"})

				db.UpdateStub = func(ctx context.Context, fn func(context.Context, kv.Tx) error) error {
					return fn(ctx, tx)
				}
				tx.DeleteBucketReturns(errors.New("delete failed"))
			})

			It("returns error", func() {
				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				Expect(recorder.Body.String()).To(ContainSubstring("remove bucket failed"))
				Expect(db.UpdateCallCount()).To(Equal(1))
				Expect(tx.DeleteBucketCallCount()).To(Equal(1))
				Expect(cancelCalled).To(BeFalse())
			})
		})

		Context("when db.Update fails", func() {
			BeforeEach(func() {
				request = httptest.NewRequest("DELETE", "/bucket/test-bucket", nil)
				request = mux.SetURLVars(request, map[string]string{"BucketName": "test-bucket"})

				db.UpdateReturns(errors.New("update failed"))
			})

			It("returns error", func() {
				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				Expect(recorder.Body.String()).To(ContainSubstring("remove bucket failed"))
				Expect(db.UpdateCallCount()).To(Equal(1))
				Expect(cancelCalled).To(BeFalse())
			})
		})

		Context("when reset is already running", func() {
			BeforeEach(func() {
				request = httptest.NewRequest("DELETE", "/bucket/test-bucket", nil)
				request = mux.SetURLVars(request, map[string]string{"BucketName": "test-bucket"})
			})

			It("prevents concurrent execution", func() {
				db.UpdateStub = func(ctx context.Context, fn func(context.Context, kv.Tx) error) error {
					// Simulate a second request while first is running
					secondRecorder := httptest.NewRecorder()
					secondRequest := httptest.NewRequest("DELETE", "/bucket/test-bucket", nil)
					secondRequest = mux.SetURLVars(secondRequest, map[string]string{"BucketName": "test-bucket"})
					handler.ServeHTTP(secondRecorder, secondRequest)

					Expect(secondRecorder.Code).To(Equal(http.StatusInternalServerError))
					Expect(secondRecorder.Body.String()).To(ContainSubstring("reset bucket test-bucket already running"))
					return fn(ctx, tx)
				}
				tx.DeleteBucketReturns(nil)

				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(db.UpdateCallCount()).To(Equal(1))
			})
		})
	})
})
