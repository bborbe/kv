// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kv_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/kv"
	"github.com/bborbe/kv/mocks"
)

var _ = Describe("RelationStoreTxString", func() {
	var relationStoreTxString kv.RelationStoreTxString
	BeforeEach(func() {
		relationStoreTxString = &mocks.RelationStoreTxString{}
	})
	It("returns relationStoreTxString", func() {
		Expect(relationStoreTxString).NotTo(BeNil())
	})
})
