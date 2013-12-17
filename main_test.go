package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Coordinator", func() {
	Describe("on startup", func() {
		It("discovers Warden services", func() {
		})

		It("populates the pool with the existing services", func() {
		})

		It("subscribes to app.start", func() {
		})
	})

	Describe("when a Warden service leaves the cluster", func() {
		It("removes it from the pool", func() {
		})

		It("starts any apps it had on that service", func() {
		})
	})

	Describe("when a Warden service joins the cluster", func() {
		It("adds it to the pool", func() {
		})
	})

	Context("when a Warden service's capacity changes", func() {
		It("updates its capacity in the pool", func() {
		})
	})

	Describe("when app.start is received", func() {
		Context("and the pool has capacity", func() {
			It("places it in the Warden zone with fewer instances", func() {
			})

			It("places it on a Warden server with fewer instances", func() {
			})

			It("places it on a Warden server with more available capacity", func() {
			})
		})

		Context("but the pool is full", func() {
			It("drops the request on the floor", func() {
			})
		})

		Context("but there are no members of the Warden pool", func() {
			It("drops the request on the floor", func() {
			})
		})
	})
})
