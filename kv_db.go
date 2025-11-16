// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv

import (
	"context"
	"errors"
)

// ErrTransactionAlreadyOpen is returned when attempting to open a transaction while one is already active.
var ErrTransactionAlreadyOpen = errors.New("transaction already open")

// TransactionAlreadyOpenError is deprecated: use ErrTransactionAlreadyOpen instead.
var TransactionAlreadyOpenError = ErrTransactionAlreadyOpen

//counterfeiter:generate -o mocks/db.go --fake-name DB . DB

// DB represents a key-value database that supports transactions and lifecycle management.
type DB interface {
	// Update opens a write transaction
	Update(ctx context.Context, fn func(ctx context.Context, tx Tx) error) error

	// View opens a read only transaction
	View(ctx context.Context, fn func(ctx context.Context, tx Tx) error) error

	// Sync database to disk
	Sync() error

	// Close database
	Close() error

	// Remove database files from disk
	Remove() error
}
