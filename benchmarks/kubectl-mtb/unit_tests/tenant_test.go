package unittests

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestTenant(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	// DestroyTenant(g)
	CreateCrds()
	CreateTenant(t, g)
	CreateTenantNS(t, g)
	DestroyTenant(g)
}
