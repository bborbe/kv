// Copyright (c) 2026 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv

type Stats struct {
	Backend string        `json:"backend"`
	SizeB   int64         `json:"size_bytes,omitempty"`
	Buckets []BucketStats `json:"buckets"`
}

type BucketStats struct {
	Name     BucketName `json:"name"`
	KeyCount int64      `json:"keys"`
	SizeB    int64      `json:"size_bytes,omitempty"`
}
