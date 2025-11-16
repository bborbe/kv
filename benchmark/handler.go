// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package benchmark

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/bborbe/errors"
	libhttp "github.com/bborbe/http"
	"github.com/bborbe/log"
	"github.com/bborbe/run"
	"github.com/golang/glog"

	"github.com/bborbe/kv"
)

const timeout = time.Minute * 10

func NewBoltBenchmarkHandler(
	provider kv.Provider,
	backgroundRunner run.BackgroundRunner,
	logSamplerFactory log.SamplerFactory,
) libhttp.WithError {
	return libhttp.WithErrorFunc(
		func(ctx context.Context, resp http.ResponseWriter, req *http.Request) error {
			amount, _ := strconv.Atoi(req.FormValue("amount"))
			if amount <= 0 {
				amount = 1000
			}
			valueLength, _ := strconv.Atoi(req.FormValue("valueLength"))
			if valueLength <= 0 {
				valueLength = 10000
			}
			batchSize, _ := strconv.Atoi(req.FormValue("batchSize"))
			if batchSize <= 0 {
				batchSize = 1
			}
			err := backgroundRunner.Run(func(ctx context.Context) error {
				ctx, cancel := context.WithTimeout(ctx, timeout)
				defer cancel()

				db, err := provider.Get(ctx)
				if err != nil {
					return errors.Wrap(ctx, err, "open bolt db failed")
				}
				defer func() {
					_ = db.Close()
				}()

				benchmark := NewBenchmark(
					db,
					logSamplerFactory,
				)
				result, err := benchmark.Benchmark(ctx, amount, valueLength, batchSize)
				if err != nil {
					return errors.Wrapf(ctx, err, "benchmark failed")
				}

				avgWriteDuration := result.WriteDuration / time.Duration(amount)
				avgReadDuration := result.ReadDuration / time.Duration(amount)
				glog.V(2).
					Infof("amount: %d valueLength: %d batchSize: %d", amount, valueLength, batchSize)
				glog.V(2).
					Infof("write duration total: %v avg: %d us", result.WriteDuration, avgWriteDuration.Microseconds())
				glog.V(2).
					Infof("read duration total: %v avg: %d us", result.ReadDuration, avgReadDuration.Microseconds())
				return nil
			})
			if err != nil {
				return errors.Wrapf(ctx, err, "run in background failed")
			}
			_, _ = libhttp.WriteAndGlog(
				resp,
				"benchmark triggered with amount %d, valueLength %d and batchSize %d",
				amount,
				valueLength,
				batchSize,
			)
			return nil
		},
	)
}
