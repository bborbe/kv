// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/kv"
)

var _ = Describe("BucketName", func() {
	Context("NewBucketName", func() {
		It("creates bucket name from string", func() {
			bucketName := kv.NewBucketName("test-bucket")
			Expect(bucketName.String()).To(Equal("test-bucket"))
		})
	})

	Context("BucketFromStrings", func() {
		It("creates bucket name from single string", func() {
			bucketName := kv.BucketFromStrings("test")
			Expect(bucketName.String()).To(Equal("test"))
		})

		It("creates bucket name from multiple strings", func() {
			bucketName := kv.BucketFromStrings("test", "bucket", "name")
			Expect(bucketName.String()).To(Equal("test_bucket_name"))
		})

		It("creates bucket name from empty strings", func() {
			bucketName := kv.BucketFromStrings("", "bucket", "")
			Expect(bucketName.String()).To(Equal("_bucket_"))
		})
	})

	Context("String", func() {
		It("returns string representation", func() {
			bucketName := kv.NewBucketName("my-bucket")
			Expect(bucketName.String()).To(Equal("my-bucket"))
		})

		It("handles empty bucket name", func() {
			bucketName := kv.NewBucketName("")
			Expect(bucketName.String()).To(Equal(""))
		})
	})

	Context("Bytes", func() {
		It("returns byte representation", func() {
			bucketName := kv.NewBucketName("test")
			Expect(bucketName.Bytes()).To(Equal([]byte("test")))
		})

		It("handles empty bucket name", func() {
			bucketName := kv.NewBucketName("")
			Expect(bucketName.Bytes()).To(Equal([]byte("")))
		})
	})

	Context("Equal", func() {
		It("returns true for equal bucket names", func() {
			bucketName1 := kv.NewBucketName("test")
			bucketName2 := kv.NewBucketName("test")
			Expect(bucketName1.Equal(bucketName2)).To(BeTrue())
		})

		It("returns false for different bucket names", func() {
			bucketName1 := kv.NewBucketName("test1")
			bucketName2 := kv.NewBucketName("test2")
			Expect(bucketName1.Equal(bucketName2)).To(BeFalse())
		})

		It("returns true for empty bucket names", func() {
			bucketName1 := kv.NewBucketName("")
			bucketName2 := kv.NewBucketName("")
			Expect(bucketName1.Equal(bucketName2)).To(BeTrue())
		})
	})
})

var _ = Describe("BucketNames", func() {
	Context("Contains", func() {
		var bucketNames kv.BucketNames

		BeforeEach(func() {
			bucketNames = kv.BucketNames{
				kv.NewBucketName("bucket1"),
				kv.NewBucketName("bucket2"),
				kv.NewBucketName("bucket3"),
			}
		})

		It("returns true when bucket name exists", func() {
			Expect(bucketNames.Contains(kv.NewBucketName("bucket2"))).To(BeTrue())
		})

		It("returns false when bucket name does not exist", func() {
			Expect(bucketNames.Contains(kv.NewBucketName("bucket4"))).To(BeFalse())
		})

		It("returns false for empty slice", func() {
			emptyBucketNames := kv.BucketNames{}
			Expect(emptyBucketNames.Contains(kv.NewBucketName("bucket1"))).To(BeFalse())
		})

		It("handles empty bucket name", func() {
			bucketNamesWithEmpty := kv.BucketNames{
				kv.NewBucketName(""),
				kv.NewBucketName("bucket1"),
			}
			Expect(bucketNamesWithEmpty.Contains(kv.NewBucketName(""))).To(BeTrue())
		})
	})
})
