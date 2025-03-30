package kv_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/kv"
	"github.com/bborbe/kv/mocks"
)

var _ = Describe("DBWithMetrics", func() {
	var ctx context.Context
	var err error
	var db *mocks.DB
	BeforeEach(func() {
		ctx = context.Background()
		db = &mocks.DB{}
	})
	Context("Update", func() {
		JustBeforeEach(func() {
			err = kv.NewDBWithMetrics(db, kv.NewMetrics()).Update(ctx, func(ctx context.Context, tx kv.Tx) error {
				return nil
			})
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("call Update", func() {
			Expect(db.UpdateCallCount()).To(Equal(1))
		})
	})
	Context("View", func() {
		JustBeforeEach(func() {
			err = kv.NewDBWithMetrics(db, kv.NewMetrics()).View(ctx, func(ctx context.Context, tx kv.Tx) error {
				return nil
			})
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("call View", func() {
			Expect(db.ViewCallCount()).To(Equal(1))
		})
	})
	Context("Sync", func() {
		JustBeforeEach(func() {
			err = kv.NewDBWithMetrics(db, kv.NewMetrics()).Sync()
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("call Sync", func() {
			Expect(db.SyncCallCount()).To(Equal(1))
		})
	})
	Context("Close", func() {
		JustBeforeEach(func() {
			err = kv.NewDBWithMetrics(db, kv.NewMetrics()).Close()
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("call Close", func() {
			Expect(db.CloseCallCount()).To(Equal(1))
		})
	})
	Context("Remove", func() {
		JustBeforeEach(func() {
			err = kv.NewDBWithMetrics(db, kv.NewMetrics()).Remove()
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("call Remove", func() {
			Expect(db.RemoveCallCount()).To(Equal(1))
		})
	})
})
