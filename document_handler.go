package proxy

import (
	"encoding/xml"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	"path/filepath"
	"strings"
)

type SolrDocument struct {
	Core    string
	Id      string
	Name    string
	content []byte
}

type Document struct {
	Field []DocField `xml:"field"`
}

type Add struct {
	Doc     Document `xml:"doc"`
	content []byte
}

type DocField struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

func ParseXMLDocument(content []byte) (*Add, error) {
	a := new(Add)
	a.content = content
	var err error
	if len(content) > 0 {
		err = xml.Unmarshal(content, a)
	}
	return a, err
}

func (a *Add) getFieldValue(fieldName string) string {
	for _, field := range a.Doc.Field {
		if field.Name == fieldName {
			return field.Value
		}
	}
	return ""
}

func (d *Add) GetSolrDocument(CoreName string) *SolrDocument {
	name, id := d.GetNameAndId()

	return &SolrDocument{
		Core:    CoreName,
		Id:      id,
		Name:    name,
		content: d.content,
	}
}

func (d *Add) GetNameAndId() (string, string) {
	value := d.getFieldValue("id")
	splits := strings.Split(value, " ")

	if len(splits) < 2 {
		return "", ""
	}

	return splits[0], splits[1]
}

func (d *SolrDocument) Cache(awsConfig *AWSConfig) error {
	if d.Name == "" {
		// TODO: how should this case be handled?
		return nil
	}
	documentName := d.GetDocumentName(awsConfig)
	auth, _ := aws.EnvAuth()
	region := aws.Region{Name: awsConfig.RegionName, S3Endpoint: awsConfig.S3Endpoint}
	svc := s3.New(auth, region)
	bucketName := awsConfig.BucketName
	bucket := svc.Bucket(bucketName)
	return bucket.Put(documentName, d.content, "text/xml", s3.AuthenticatedRead, s3.Options{})
}

func (d *SolrDocument) GetDocumentName(awsConfig *AWSConfig) string {
	return filepath.Join(awsConfig.BucketPrefix, d.Core, d.Name, d.Id)
}
