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
	content, _ := ioutil.ReadFile("fixtures/document_sample.xml")
	doc, err := ParseXMLDocument(content)
	c.Assert(err, IsNil)
	c.Assert(len(doc.Doc.Field), Equals, 56)
}

func (s *DocumentSuite) TestDocumentParserWithBadXml(c *C) {
	content := []byte(`<commit />`)
	doc, err := ParseXMLDocument(content)
	c.Assert(err, IsNil)
	c.Assert(len(doc.Doc.Field), Equals, 0)
}

func (s *DocumentSuite) TestDocumentParserWithEmptyBody(c *C) {
	content := []byte("")
	doc, err := ParseXMLDocument(content)
	c.Assert(err, IsNil)
	c.Assert(len(doc.Doc.Field), Equals, 0)
}

func (s *DocumentSuite) TestGetFieldValue(c *C) {
	content, _ := ioutil.ReadFile("fixtures/document_sample.xml")
	doc, err := ParseXMLDocument(content)
	c.Assert(err, IsNil)
	c.Assert(doc.getFieldValue("id"), Equals, "Hotel 5000000000000")
}

func (s *DocumentSuite) TestGetFieldValueBlank(c *C) {
	content, _ := ioutil.ReadFile("fixtures/document_sample.xml")
	doc, err := ParseXMLDocument(content)
	c.Assert(err, IsNil)
	c.Assert(doc.getFieldValue("id-missing"), Equals, "")
}

func (s *DocumentSuite) TestGetNameAndId(c *C) {
	content, _ := ioutil.ReadFile("fixtures/document_sample.xml")
	doc, err := ParseXMLDocument(content)
	c.Assert(err, IsNil)
	name, id := doc.GetNameAndId()
	c.Assert(name, Equals, "Hotel")
	c.Assert(id, Equals, "5000000000000")
}

func (s *DocumentSuite) TestSolrDocument(c *C) {
	content, _ := ioutil.ReadFile("fixtures/document_sample.xml")
	doc, err := ParseXMLDocument(content)
	c.Assert(err, IsNil)
	solrDoc := doc.GetSolrDocument("collection1")
	c.Assert(solrDoc.Id, Equals, "5000000000000")
	c.Assert(solrDoc.Name, Equals, "Hotel")
}

func (s *DocumentSuite) TestDocumentNameWithBucketPrevix(c *C) {
	doc := &SolrDocument{
		Id:      "1",
		Name:    "Hotel",
		content: []byte(""),
	}
	config := &AWSConfig{
		BucketPrefix: "staging",
	}
	documentName := doc.GetDocumentName(config)
	c.Assert(documentName, Equals, "staging/Hotel/1")
}

func (s *DocumentSuite) TestDocumentNameWithoutBucketPrefix(c *C) {
	doc := &SolrDocument{
		Id:      "1",
		Name:    "Hotel",
		content: []byte(""),
	}
	config := &AWSConfig{}
	documentName := doc.GetDocumentName(config)
	c.Assert(documentName, Equals, "Hotel/1")
}

func (s *DocumentSuite) TestDocumentNameWithCoreName(c *C) {
	doc := &SolrDocument{
		Core:    "collection1",
		Id:      "1",
		Name:    "Hotel",
		content: []byte(""),
	}
	config := &AWSConfig{}
	documentName := doc.GetDocumentName(config)
	c.Assert(documentName, Equals, "collection1/Hotel/1")
}
