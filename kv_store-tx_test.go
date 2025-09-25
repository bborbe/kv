// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv_test

import (
	"context"
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/kv"
	"github.com/bborbe/kv/mocks"
)

var _ = Describe("StoreTx", func() {
	var ctx context.Context
	var tx *mocks.Tx
	var bucket *mocks.Bucket
	var item *mocks.Item
	var iterator *mocks.Iterator
	var storeTx kv.StoreTx[string, TestObject]
	var bucketName kv.BucketName

	BeforeEach(func() {
		ctx = context.Background()
		tx = &mocks.Tx{}
		bucket = &mocks.Bucket{}
		item = &mocks.Item{}
		iterator = &mocks.Iterator{}
		bucketName = kv.NewBucketName("test-bucket")
		storeTx = kv.NewStoreTx[string, TestObject](bucketName)
	})

	Describe("NewStoreTx", func() {
		It("creates a new StoreTx with bucket name", func() {
			storeTx := kv.NewStoreTx[string, TestObject](bucketName)
			Expect(storeTx).NotTo(BeNil())
		})
	})

	Describe("Add", func() {
		BeforeEach(func() {
			tx.CreateBucketIfNotExistsReturns(bucket, nil)
			bucket.PutReturns(nil)
		})

		It("marshals object and stores in bucket", func() {
			testObj := TestObject{Name: "John", Age: 30}

			err := storeTx.Add(ctx, tx, "key1", testObj)
			Expect(err).To(BeNil())

			Expect(tx.CreateBucketIfNotExistsCallCount()).To(Equal(1))
			_, actualBucketName := tx.CreateBucketIfNotExistsArgsForCall(0)
			Expect(actualBucketName).To(Equal(bucketName))

			Expect(bucket.PutCallCount()).To(Equal(1))
			_, actualKey, actualValue := bucket.PutArgsForCall(0)
			Expect(actualKey).To(Equal([]byte("key1")))

			var unmarshaled TestObject
			err = json.Unmarshal(actualValue, &unmarshaled)
			Expect(err).To(BeNil())
			Expect(unmarshaled).To(Equal(testObj))
		})

		It("returns error when CreateBucketIfNotExists fails", func() {
			expectedErr := errors.New("bucket creation failed")
			tx.CreateBucketIfNotExistsReturns(nil, expectedErr)

			err := storeTx.Add(ctx, tx, "key1", TestObject{})
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("get bucket failed"))
		})

		It("returns error when Put fails", func() {
			expectedErr := errors.New("put failed")
			bucket.PutReturns(expectedErr)

			err := storeTx.Add(ctx, tx, "key1", TestObject{})
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("set failed"))
		})
	})

	Describe("Remove", func() {
		BeforeEach(func() {
			tx.CreateBucketIfNotExistsReturns(bucket, nil)
			bucket.DeleteReturns(nil)
		})

		It("removes key from bucket", func() {
			err := storeTx.Remove(ctx, tx, "key1")
			Expect(err).To(BeNil())

			Expect(bucket.DeleteCallCount()).To(Equal(1))
			_, actualKey := bucket.DeleteArgsForCall(0)
			Expect(actualKey).To(Equal([]byte("key1")))
		})

		It("handles bucket not found gracefully", func() {
			tx.CreateBucketIfNotExistsReturns(nil, kv.BucketNotFoundError)

			err := storeTx.Remove(ctx, tx, "key1")
			Expect(err).To(BeNil())
		})

		It("returns error when CreateBucketIfNotExists fails with other error", func() {
			expectedErr := errors.New("other error")
			tx.CreateBucketIfNotExistsReturns(nil, expectedErr)

			err := storeTx.Remove(ctx, tx, "key1")
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("get bucket failed"))
		})

		It("returns error when Delete fails", func() {
			expectedErr := errors.New("delete failed")
			bucket.DeleteReturns(expectedErr)

			err := storeTx.Remove(ctx, tx, "key1")
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("remove key1 failed"))
		})
	})

	Describe("Get", func() {
		BeforeEach(func() {
			tx.BucketReturns(bucket, nil)
			bucket.GetReturns(item, nil)
		})

		It("retrieves and unmarshals object", func() {
			testObj := TestObject{Name: "John", Age: 30}
			data, _ := json.Marshal(testObj)

			item.ValueStub = func(fn func([]byte) error) error {
				return fn(data)
			}

			result, err := storeTx.Get(ctx, tx, "key1")
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(*result).To(Equal(testObj))

			Expect(tx.BucketCallCount()).To(Equal(1))
			_, actualBucketName := tx.BucketArgsForCall(0)
			Expect(actualBucketName).To(Equal(bucketName))
		})

		It("returns error when key not found", func() {
			item.ValueStub = func(fn func([]byte) error) error {
				return fn([]byte{})
			}

			result, err := storeTx.Get(ctx, tx, "key1")
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("key(key1) not found"))
			Expect(result).To(BeNil())
		})

		It("returns error when bucket not found", func() {
			tx.BucketReturns(nil, kv.BucketNotFoundError)

			result, err := storeTx.Get(ctx, tx, "key1")
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("get bucket failed"))
			Expect(result).To(BeNil())
		})

		It("returns error when Get fails", func() {
			expectedErr := errors.New("get failed")
			bucket.GetReturns(nil, expectedErr)

			result, err := storeTx.Get(ctx, tx, "key1")
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("get key1 failed"))
			Expect(result).To(BeNil())
		})

		It("returns error when unmarshal fails", func() {
			item.ValueStub = func(fn func([]byte) error) error {
				return fn([]byte("invalid json"))
			}

			result, err := storeTx.Get(ctx, tx, "key1")
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("handel value failed"))
			Expect(result).To(BeNil())
		})
	})

	Describe("Exists", func() {
		BeforeEach(func() {
			tx.BucketReturns(bucket, nil)
			bucket.GetReturns(item, nil)
		})

		It("returns true when item exists", func() {
			item.ExistsReturns(true)

			exists, err := storeTx.Exists(ctx, tx, "key1")
			Expect(err).To(BeNil())
			Expect(exists).To(BeTrue())
		})

		It("returns false when item does not exist", func() {
			item.ExistsReturns(false)

			exists, err := storeTx.Exists(ctx, tx, "key1")
			Expect(err).To(BeNil())
			Expect(exists).To(BeFalse())
		})

		It("returns false when bucket not found", func() {
			tx.BucketReturns(nil, kv.BucketNotFoundError)

			exists, err := storeTx.Exists(ctx, tx, "key1")
			Expect(err).To(BeNil())
			Expect(exists).To(BeFalse())
		})

		It("returns error when bucket retrieval fails with other error", func() {
			expectedErr := errors.New("other error")
			tx.BucketReturns(nil, expectedErr)

			exists, err := storeTx.Exists(ctx, tx, "key1")
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("get bucket failed"))
			Expect(exists).To(BeFalse())
		})

		It("returns error when Get fails", func() {
			expectedErr := errors.New("get failed")
			bucket.GetReturns(nil, expectedErr)

			exists, err := storeTx.Exists(ctx, tx, "key1")
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("get key1 failed"))
			Expect(exists).To(BeFalse())
		})
	})

	Describe("Map", func() {
		BeforeEach(func() {
			tx.BucketReturns(bucket, nil)
			bucket.IteratorReturns(iterator)
		})

		It("iterates over all items in bucket", func() {
			testObj1 := TestObject{Name: "John", Age: 30}
			testObj2 := TestObject{Name: "Jane", Age: 25}

			data1, _ := json.Marshal(testObj1)
			data2, _ := json.Marshal(testObj2)

			callCount := 0
			iterator.ValidStub = func() bool {
				callCount++
				return callCount <= 2
			}

			itemCallCount := 0
			iterator.ItemStub = func() kv.Item {
				itemCallCount++
				mockItem := &mocks.Item{}
				if itemCallCount == 1 {
					mockItem.KeyReturns([]byte("key1"))
					mockItem.ValueStub = func(fn func([]byte) error) error {
						return fn(data1)
					}
				} else {
					mockItem.KeyReturns([]byte("key2"))
					mockItem.ValueStub = func(fn func([]byte) error) error {
						return fn(data2)
					}
				}
				return mockItem
			}

			var results []TestObject
			err := storeTx.Map(ctx, tx, func(ctx context.Context, key string, object TestObject) error {
				results = append(results, object)
				return nil
			})

			Expect(err).To(BeNil())
			Expect(len(results)).To(Equal(2))
			Expect(results[0]).To(Equal(testObj1))
			Expect(results[1]).To(Equal(testObj2))
		})

		It("handles bucket not found gracefully", func() {
			tx.BucketReturns(nil, kv.BucketNotFoundError)

			callCount := 0
			err := storeTx.Map(ctx, tx, func(ctx context.Context, key string, object TestObject) error {
				callCount++
				return nil
			})

			Expect(err).To(BeNil())
			Expect(callCount).To(Equal(0))
		})

		It("returns error when context is cancelled", func() {
			iterator.ValidReturns(true)

			cancelCtx, cancel := context.WithCancel(ctx)
			cancel()

			err := storeTx.Map(cancelCtx, tx, func(ctx context.Context, key string, object TestObject) error {
				return nil
			})

			Expect(err).To(Equal(context.Canceled))
		})

		It("returns error when function fails", func() {
			testObj := TestObject{Name: "John", Age: 30}
			data, _ := json.Marshal(testObj)

			iterator.ValidReturns(true)
			mockItem := &mocks.Item{}
			mockItem.KeyReturns([]byte("key1"))
			mockItem.ValueStub = func(fn func([]byte) error) error {
				return fn(data)
			}
			iterator.ItemReturns(mockItem)

			expectedErr := errors.New("function failed")
			err := storeTx.Map(ctx, tx, func(ctx context.Context, key string, object TestObject) error {
				return expectedErr
			})

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("call fn failed"))
		})

		It("returns error when unmarshal fails", func() {
			iterator.ValidReturns(true)
			mockItem := &mocks.Item{}
			mockItem.KeyReturns([]byte("key1"))
			mockItem.ValueStub = func(fn func([]byte) error) error {
				return fn([]byte("invalid json"))
			}
			iterator.ItemReturns(mockItem)

			err := storeTx.Map(ctx, tx, func(ctx context.Context, key string, object TestObject) error {
				return nil
			})

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("handle value failed"))
		})
	})

	Describe("Stream", func() {
		BeforeEach(func() {
			tx.BucketReturns(bucket, nil)
			bucket.IteratorReturns(iterator)
		})

		It("streams objects through channel", func() {
			testObj := TestObject{Name: "John", Age: 30}
			data, _ := json.Marshal(testObj)

			callCount := 0
			iterator.ValidStub = func() bool {
				callCount++
				return callCount <= 1
			}

			mockItem := &mocks.Item{}
			mockItem.KeyReturns([]byte("key1"))
			mockItem.ValueStub = func(fn func([]byte) error) error {
				return fn(data)
			}
			iterator.ItemReturns(mockItem)

			ch := make(chan TestObject, 10)
			defer close(ch)

			err := storeTx.Stream(ctx, tx, ch)
			Expect(err).To(BeNil())

			select {
			case obj := <-ch:
				Expect(obj).To(Equal(testObj))
			default:
				Fail("Expected object in channel")
			}
		})

		It("handles context cancellation", func() {
			iterator.ValidReturns(true)

			cancelCtx, cancel := context.WithCancel(ctx)
			cancel()

			ch := make(chan TestObject, 10)
			defer close(ch)

			err := storeTx.Stream(cancelCtx, tx, ch)
			Expect(err).To(Equal(context.Canceled))
		})
	})

	Describe("List", func() {
		BeforeEach(func() {
			tx.BucketReturns(bucket, nil)
			bucket.IteratorReturns(iterator)
		})

		It("returns list of all objects", func() {
			testObj1 := TestObject{Name: "John", Age: 30}
			testObj2 := TestObject{Name: "Jane", Age: 25}

			data1, _ := json.Marshal(testObj1)
			data2, _ := json.Marshal(testObj2)

			callCount := 0
			iterator.ValidStub = func() bool {
				callCount++
				return callCount <= 2
			}

			itemCallCount := 0
			iterator.ItemStub = func() kv.Item {
				itemCallCount++
				mockItem := &mocks.Item{}
				if itemCallCount == 1 {
					mockItem.KeyReturns([]byte("key1"))
					mockItem.ValueStub = func(fn func([]byte) error) error {
						return fn(data1)
					}
				} else {
					mockItem.KeyReturns([]byte("key2"))
					mockItem.ValueStub = func(fn func([]byte) error) error {
						return fn(data2)
					}
				}
				return mockItem
			}

			objects, err := storeTx.List(ctx, tx)
			Expect(err).To(BeNil())
			Expect(objects).NotTo(BeNil())
			Expect(len(objects)).To(Equal(2))
			Expect(objects[0]).To(Equal(testObj1))
			Expect(objects[1]).To(Equal(testObj2))
		})

		It("returns error when Map fails", func() {
			tx.BucketReturns(nil, errors.New("bucket error"))

			objects, err := storeTx.List(ctx, tx)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("map failed"))
			Expect(objects).To(BeNil())
		})

		It("returns empty list when no items exist", func() {
			iterator.ValidReturns(false)

			objects, err := storeTx.List(ctx, tx)
			Expect(err).To(BeNil())
			Expect(objects).NotTo(BeNil())
			Expect(len(objects)).To(Equal(0))
		})
	})
})
