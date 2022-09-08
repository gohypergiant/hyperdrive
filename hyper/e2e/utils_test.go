package e2e_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/exec"

	"github.com/onsi/gomega/gexec"
)

func buildHyper() string {
	hyperPath, err := gexec.Build("github.com/gohypergiant/hyperdrive")
	Expect(err).NotTo(HaveOccurred())

	return hyperPath
}

func runHyper(path string, args string) *gexec.Session {
	cmd := exec.Command(fmt.Sprintf("%s %s", path, args))
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	return session
}
