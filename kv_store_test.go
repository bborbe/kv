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

type TestObject struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var _ = Describe("Store", func() {
	var ctx context.Context
	var db *mocks.DB
	var tx *mocks.Tx
	var bucket *mocks.Bucket
	var item *mocks.Item
	var storeTx kv.StoreTx[string, TestObject]
	var store kv.Store[string, TestObject]
	var bucketName kv.BucketName

	BeforeEach(func() {
		ctx = context.Background()
		db = &mocks.DB{}
		tx = &mocks.Tx{}
		bucket = &mocks.Bucket{}
		item = &mocks.Item{}
		bucketName = kv.NewBucketName("test-bucket")
		storeTx = kv.NewStoreTx[string, TestObject](bucketName)
		store = kv.NewStoreFromTx(db, storeTx)
	})

	Describe("NewStore", func() {
		It("creates a store with given db and bucket name", func() {
			store := kv.NewStore[string, TestObject](db, bucketName)
			Expect(store).NotTo(BeNil())
		})
	})

	Describe("NewStoreFromTx", func() {
		It("creates a store from existing StoreTx", func() {
			store := kv.NewStoreFromTx(db, storeTx)
			Expect(store).NotTo(BeNil())
		})
	})

	Describe("Add", func() {
		It("calls db.Update and storeTx.Add", func() {
			testObj := TestObject{Name: "John", Age: 30}

			tx.CreateBucketIfNotExistsReturns(bucket, nil)
			bucket.PutReturns(nil)

			db.UpdateStub = func(ctx context.Context, fn func(context.Context, kv.Tx) error) error {
				return fn(ctx, tx)
			}

			err := store.Add(ctx, "key1", testObj)
			Expect(err).To(BeNil())
			Expect(db.UpdateCallCount()).To(Equal(1))
		})

		It("returns error when db.Update fails", func() {
			testObj := TestObject{Name: "John", Age: 30}
			expectedErr := errors.New("update failed")

			db.UpdateReturns(expectedErr)

			err := store.Add(ctx, "key1", testObj)
			Expect(err).To(Equal(expectedErr))
		})
	})

	Describe("Remove", func() {
		It("calls db.Update and storeTx.Remove", func() {
			tx.CreateBucketIfNotExistsReturns(bucket, nil)
			bucket.DeleteReturns(nil)

			db.UpdateStub = func(ctx context.Context, fn func(context.Context, kv.Tx) error) error {
				return fn(ctx, tx)
			}

			err := store.Remove(ctx, "key1")
			Expect(err).To(BeNil())
			Expect(db.UpdateCallCount()).To(Equal(1))
		})

		It("returns error when db.Update fails", func() {
			expectedErr := errors.New("update failed")

			db.UpdateReturns(expectedErr)

			err := store.Remove(ctx, "key1")
			Expect(err).To(Equal(expectedErr))
		})
	})

	Describe("Get", func() {
		It("calls db.View and storeTx.Get", func() {
			tx.BucketReturns(bucket, nil)
			bucket.GetReturns(item, nil)
			item.ValueStub = func(fn func([]byte) error) error {
				return fn([]byte{})
			}

			db.ViewStub = func(ctx context.Context, fn func(context.Context, kv.Tx) error) error {
				return fn(ctx, tx)
			}

			result, err := store.Get(ctx, "key1")
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
			Expect(db.ViewCallCount()).To(Equal(1))
		})

		It("returns error when db.View fails", func() {
			expectedErr := errors.New("view failed")

			db.ViewReturns(expectedErr)

			result, err := store.Get(ctx, "key1")
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("view failed"))
			Expect(result).To(BeNil())
		})
	})

	Describe("Exists", func() {
		It("calls db.View and storeTx.Exists", func() {
			tx.BucketReturns(bucket, nil)
			bucket.GetReturns(item, nil)
			item.ExistsReturns(false)

			db.ViewStub = func(ctx context.Context, fn func(context.Context, kv.Tx) error) error {
				return fn(ctx, tx)
			}

			exists, err := store.Exists(ctx, "key1")
			Expect(err).To(BeNil())
			Expect(exists).To(BeFalse())
			Expect(db.ViewCallCount()).To(Equal(1))
		})

		It("returns error when db.View fails", func() {
			expectedErr := errors.New("view failed")

			db.ViewReturns(expectedErr)

			exists, err := store.Exists(ctx, "key1")
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("view failed"))
			Expect(exists).To(BeFalse())
		})
	})

	Describe("Map", func() {
		It("calls db.View and storeTx.Map", func() {
			tx.BucketReturns(bucket, nil)
			iterator := &mocks.Iterator{}
			bucket.IteratorReturns(iterator)
			iterator.ValidReturns(false)

			db.ViewStub = func(ctx context.Context, fn func(context.Context, kv.Tx) error) error {
				return fn(ctx, tx)
			}

			callCount := 0
			err := store.Map(ctx, func(ctx context.Context, key string, object TestObject) error {
				callCount++
				return nil
			})
			Expect(err).To(BeNil())
			Expect(db.ViewCallCount()).To(Equal(1))
		})

		It("returns error when db.View fails", func() {
			expectedErr := errors.New("view failed")

			db.ViewReturns(expectedErr)

			err := store.Map(ctx, func(ctx context.Context, key string, object TestObject) error {
				return nil
			})
			Expect(err).To(Equal(expectedErr))
		})
	})

	Describe("Stream", func() {
		It("streams objects through channel", func() {
			ch := make(chan TestObject, 10)
			defer close(ch)

			tx.BucketReturns(bucket, nil)
			iterator := &mocks.Iterator{}
			bucket.IteratorReturns(iterator)
			iterator.ValidReturns(false)

			db.ViewStub = func(ctx context.Context, fn func(context.Context, kv.Tx) error) error {
				return fn(ctx, tx)
			}

			err := store.Stream(ctx, ch)
			Expect(err).To(BeNil())
			Expect(db.ViewCallCount()).To(Equal(1))
		})

		It("handles context cancellation", func() {
			ch := make(chan TestObject, 10)
			defer close(ch)

			cancelCtx, cancel := context.WithCancel(ctx)
			cancel()

			tx.BucketReturns(bucket, nil)
			iterator := &mocks.Iterator{}
			bucket.IteratorReturns(iterator)
			iterator.ValidReturns(false)

			db.ViewStub = func(ctx context.Context, fn func(context.Context, kv.Tx) error) error {
				// Simulate the Map function being called which checks context
				return fn(ctx, tx)
			}

			err := store.Stream(cancelCtx, ch)
			Expect(err).To(BeNil())
		})
	})
})
