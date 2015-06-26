package proxy

import (
	. "gopkg.in/check.v1"
	"io/ioutil"
	"testing"
)

func TestDocument(t *testing.T) { TestingT(t) }

type DocumentSuite struct{}

var _ = Suite(&DocumentSuite{})

func (s *DocumentSuite) TestDocumentParser(c *C) {
	content, _ := ioutil.ReadFile("document_sample.xml")
	doc := ParseXMLDocument(content)
	c.Assert(len(doc.Doc.Field), Equals, 56)
}

func (s *DocumentSuite) TestDocumentParserWithBadXml(c *C) {
	content := []byte(`<commit />`)
	doc := ParseXMLDocument(content)
	c.Assert(len(doc.Doc.Field), Equals, 0)
}

func (s *DocumentSuite) TestGetFieldValue(c *C) {
	content, _ := ioutil.ReadFile("document_sample.xml")
	doc := ParseXMLDocument(content)
	c.Assert(doc.getFieldValue("id"), Equals, "Hotel 5000000000000")
}

func (s *DocumentSuite) TestGetFieldValueBlank(c *C) {
	content, _ := ioutil.ReadFile("document_sample.xml")
	doc := ParseXMLDocument(content)
	c.Assert(doc.getFieldValue("id-missing"), Equals, "")
}

func (s *DocumentSuite) TestGetNameAndId(c *C) {
	content, _ := ioutil.ReadFile("document_sample.xml")
	doc := ParseXMLDocument(content)
	name, id := doc.GetNameAndId()
	c.Assert(name, Equals, "Hotel")
	c.Assert(id, Equals, "5000000000000")
}

func (s *DocumentSuite) TestSolrDocument(c *C) {
	content, _ := ioutil.ReadFile("document_sample.xml")
	doc := ParseXMLDocument(content)
	solrDoc := doc.GetSolrDocument()
	c.Assert(solrDoc.Id, Equals, "5000000000000")
	c.Assert(solrDoc.Name, Equals, "Hotel")
}
