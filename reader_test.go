package proxy

import (
	"github.com/mailgun/oxy/testutils"
	. "gopkg.in/check.v1"
	"net/http"
	"testing"
)

// /solr/{}/verb

func TestReader(t *testing.T) { TestingT(t) }

type ReaderSuite struct{}

var _ = Suite(&ReaderSuite{})

func (s *ReaderSuite) TestReaderRoundRobin(c *C) {
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

	reader := NewReader([]string{target1, target2})

	proxyWrapper := testutils.NewHandler(func(w http.ResponseWriter, req *http.Request) {
		reader.ServeHTTP(w, req)
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
