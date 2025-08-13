// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv

import (
	"context"
	"errors"
)

// BucketNotFoundError is returned when attempting to access a bucket that does not exist.
var BucketNotFoundError = errors.New("bucket not found")

// BucketAlreadyExistsError is returned when attempting to create a bucket that already exists.
var BucketAlreadyExistsError = errors.New("bucket already exists")

//counterfeiter:generate -o mocks/tx.go --fake-name Tx . Tx

// Tx represents a database transaction that provides bucket management operations.
type Tx interface {
	Bucket(ctx context.Context, name BucketName) (Bucket, error)
	CreateBucket(ctx context.Context, name BucketName) (Bucket, error)
	CreateBucketIfNotExists(ctx context.Context, name BucketName) (Bucket, error)
	DeleteBucket(ctx context.Context, name BucketName) error
	ListBucketNames(ctx context.Context) (BucketNames, error)
}
