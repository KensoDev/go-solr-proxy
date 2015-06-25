package proxy

import (
	"github.com/mailgun/oxy/testutils"
	. "gopkg.in/check.v1"
	"net/http"
	"testing"
)

// /solr/{}/verb

func TestUpdater(t *testing.T) { TestingT(t) }

type UpdaterSuite struct{}

var _ = Suite(&UpdaterSuite{})

func (s *UpdaterSuite) TestUpdaterRoundRobin(c *C) {
	called1 := false
	srv1 := testutils.NewHandler(func(w http.ResponseWriter, req *http.Request) {
		called1 = true
	})
	defer srv1.Close()

	updater := NewUpdater(srv1.URL)

	proxyWrapper := testutils.NewHandler(func(w http.ResponseWriter, req *http.Request) {
		updater.ServeHTTP(w, req)
	})
	defer proxyWrapper.Close()

	re, _, err := testutils.Get(proxyWrapper.URL, testutils.Headers(http.Header{}))
	c.Assert(err, IsNil)
	c.Assert(re.StatusCode, Equals, http.StatusOK)
	c.Assert(called1, Equals, true)
}
