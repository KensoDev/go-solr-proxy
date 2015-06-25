package proxy

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestProxy(t *testing.T) { TestingT(t) }

type ProxySuite struct{}

var _ = Suite(&ProxySuite{})
