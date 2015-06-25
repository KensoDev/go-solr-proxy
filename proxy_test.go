package proxy

import (
	"github.com/mailgun/oxy/testutils"
	. "gopkg.in/check.v1"
	"net/http"
	"testing"
)

// /solr/{}/verb

func TestProxy(t *testing.T) { TestingT(t) }

type ProxySuite struct{}

var _ = Suite(&ProxySuite{})

func (s *ProxySuite) TestRoundRobin(c *C) {
	// create target1, target2
	// create Proxy
	// call /query
	//assert target1 or 2 called/
	// call /query
	//assert the other was called

	called1 := false
	srv1 := testutils.NewHandler(func(w http.ResponseWriter, req *http.Request) {
		called1 = true
	})
	defer srv1.Close()

	called2 := false
	srv2 := testutils.NewHandler(func(w http.ResponseWriter, req *http.Request) {
		called2 = true
	})
	defer srv2.Close()
	target1 := srv1.URL
	target2 := srv2.URL
	proxy := NewProxy("MASTER", []string{target1, target2})

	proxyWrapper := testutils.NewHandler(func(w http.ResponseWriter, req *http.Request) {
		proxy.ServeHTTP(w, req)
	})
	defer proxyWrapper.Close()

	re, _, err := testutils.Get(proxyWrapper.URL, testutils.Headers(http.Header{}))
	c.Assert(err, IsNil)
	c.Assert(re.StatusCode, Equals, http.StatusOK)
	c.Assert(called1, Equals, true)
	c.Assert(called2, Equals, false)

	re, _, err = testutils.Get(proxyWrapper.URL, testutils.Headers(http.Header{}))
	c.Assert(err, IsNil)
	c.Assert(re.StatusCode, Equals, http.StatusOK)
	c.Assert(called2, Equals, true)
}
