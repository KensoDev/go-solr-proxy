package proxy

import (
	"encoding/xml"
	"fmt"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	"strings"
)

type SolrDocument struct {
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

func ParseXMLDocument(content []byte) *Add {
	a := new(Add)
	xml.Unmarshal(content, a)
	a.content = content
	return a
}

func (a *Add) getFieldValue(fieldName string) string {
	for _, field := range a.Doc.Field {
		if field.Name == fieldName {
			return field.Value
		}
	}
	return ""
}

func (d *Add) GetSolrDocument() *SolrDocument {
	name, id := d.GetNameAndId()
	return &SolrDocument{
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

func (d *SolrDocument) Cache(awsConfig *AWSConfig) {
	if d.Name == "" {
		return
	}
	documentName := fmt.Sprintf("%s/%s", d.Name, d.Id)
	auth, _ := aws.EnvAuth()
	region := aws.Region{Name: awsConfig.RegionName, S3Endpoint: awsConfig.S3Endpoint}
	svc := s3.New(auth, region)
	bucketName := awsConfig.BucketName
	bucket := svc.Bucket(bucketName)
	err := bucket.Put(documentName, d.content, "text/xml", s3.AuthenticatedRead, s3.Options{})
	if err != nil {
	}
}
