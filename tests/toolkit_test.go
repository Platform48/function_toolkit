package toolkits_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

func TestToolkit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Toolkit OkJson")
	RunSpecs(t, "Toolkit ErrJson")
}
