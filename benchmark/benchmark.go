// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package benchmark

import (
	"context"
	"strconv"
	"time"

	"github.com/bborbe/errors"
	"github.com/bborbe/log"
	"github.com/golang/glog"

	"github.com/bborbe/kv"
)

type Result struct {
	Amount        int
	ValueSize     int
	WriteDuration time.Duration
	ReadDuration  time.Duration
}

type Benchmark interface {
	Benchmark(
		ctx context.Context,
		amount int, // how many read and write are performt
		valueSize int, // how many character the value is long
		batchSize int, // how many record read and write per transaction
	) (*Result, error)
}

func NewBenchmark(
	db kv.DB,
	logSamplerFactory log.SamplerFactory,
) Benchmark {
	return &benchmark{
		db:              db,
		readLogSampler:  logSamplerFactory.Sampler(),
		writeLogSampler: logSamplerFactory.Sampler(),
		bucketName:      kv.BucketName("test"),
	}
}

type benchmark struct {
	db              kv.DB
	writeLogSampler log.Sampler
	readLogSampler  log.Sampler
	bucketName      kv.BucketName
}

func (b *benchmark) Benchmark(
	ctx context.Context,
	amount int,
	valueSize int,
	batchSize int,
) (*Result, error) {

	keys := generateKeys(amount)
	result := Result{
		Amount:    amount,
		ValueSize: valueSize,
	}

	{
		startTime := time.Now()
		if err := b.write(ctx, keys, batchSize, valueSize); err != nil {
			return nil, errors.Wrapf(ctx, err, "write failed")
		}
		result.WriteDuration = time.Since(startTime)
		glog.V(2).Infof("write %d records complete after: %v", amount, result.WriteDuration)
	}

	{
		startTime := time.Now()
		if err := b.read(ctx, keys, batchSize); err != nil {
			return nil, errors.Wrapf(ctx, err, "write failed")
		}
		result.ReadDuration = time.Since(startTime)
		glog.V(2).Infof("read %d records complete after: %v", amount, result.ReadDuration)
	}

	return &result, nil
}

func (b *benchmark) read(
	ctx context.Context,
	keys []string,
	batchSize int,
) error {
	ShuffleSlice(keys)
	amount := len(keys)
	counter := 0
	for {
		if counter == amount {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := b.db.View(ctx, func(ctx context.Context, tx kv.Tx) error {
				bucket, err := tx.Bucket(ctx, b.bucketName)
				if err != nil {
					return errors.Wrapf(ctx, err, "get bucket failed")
				}

				for t := 0; t < batchSize; t++ {
					item, err := bucket.Get(ctx, []byte(keys[counter]))
					if err != nil {
						return errors.Wrapf(ctx, err, "get failed")
					}
					err = item.Value(func(value []byte) error {
						if len(value) == 0 {
							return errors.Errorf(ctx, "empty value")
						}
						return nil
					})
					if err != nil {
						return errors.Wrapf(ctx, err, "value failed")
					}
					counter++
					if counter == amount {
						return nil
					}
				}
				return nil
			})
			if err != nil {
				return errors.Wrapf(ctx, err, "read failed")
			}
			if b.readLogSampler.IsSample() {
				glog.Infof("read %d of %d completed", counter+1, amount)
			}
		}
	}
}

func (b *benchmark) write(
	ctx context.Context,
	keys []string,
	batchSize int,
	valueSize int,
) error {
	value := RandString(valueSize)
	amount := len(keys)
	ShuffleSlice(keys)
	counter := 0
	for {
		if counter == amount {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := b.db.Update(ctx, func(ctx context.Context, tx kv.Tx) error {
				bucket, err := tx.CreateBucketIfNotExists(ctx, b.bucketName)
				if err != nil {
					return errors.Wrapf(ctx, err, "create bucket failed")
				}

				for t := 0; t < batchSize; t++ {
					if err := bucket.Put(ctx, []byte(keys[counter]), []byte(value)); err != nil {
						return errors.Wrapf(ctx, err, "write failed")
					}
					counter++
					if counter == amount {
						return nil
					}
				}
				return nil
			})
			if err != nil {
				return errors.Wrapf(ctx, err, "write failed")
			}
			if b.writeLogSampler.IsSample() {
				glog.Infof("write %d of %d completed", counter, amount)
			}
		}
	}
}

func generateKeys(amount int) []string {
	keys := make([]string, amount)
	for s := 0; s < amount; s++ {
		keys[s] = strconv.Itoa(s)
	}
	return keys
}
