// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

// NewResetBucketHandler returns a http.Handler
// that allow delete a bucket
func NewResetBucketHandler(db DB, cancel context.CancelFunc) http.Handler {
	var lock sync.Mutex
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		bucketName := BucketName(vars["BucketName"])
		glog.V(2).Infof("reset bucket %s started", bucketName)

		if lock.TryLock() == false {
			http.Error(resp, fmt.Sprintf("reset db already running"), http.StatusInternalServerError)
			return
		}
		defer lock.Unlock()

		ctx := context.Background()
		if len(bucketName) == 0 {
			http.Error(resp, "parameter bucket missing", http.StatusBadRequest)
			return
		}
		err := db.Update(ctx, func(ctx context.Context, tx Tx) error {
			return tx.DeleteBucket(ctx, bucketName)
		})
		if err != nil {
			http.Error(resp, fmt.Sprintf("remove bucket failed: %v", err), http.StatusInternalServerError)
			return
		}
		resp.WriteHeader(http.StatusOK)
		fmt.Fprintf(resp, "reset bucket %s successful\n", bucketName)
		glog.V(2).Infof("reset bucket %s successful", bucketName)
		cancel()
	})
}
