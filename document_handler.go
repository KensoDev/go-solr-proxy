package proxy

import (
	"encoding/xml"
	"strings"
)

type SolrDocument struct {
	Id   string
	Name string
}

type Document struct {
	Field []DocField `xml:"field"`
}

type Add struct {
	Doc Document `xml:"doc"`
}

type DocField struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

func ParseXMLDocument(content []byte) (d *Add) {
	e := new(Add)
	xml.Unmarshal(content, e)
	return e
}

func (a *Add) getFieldValue(fieldName string) (v string) {
	for _, field := range a.Doc.Field {
		if field.Name == fieldName {
			return field.Value
		}
	}
	return ""
}

func (d *Add) GetSolrDocument() (solrDoc *SolrDocument) {
	name, id := d.GetNameAndId()
	return &SolrDocument{
		Id:   id,
		Name: name,
	}
}

func (d *Add) GetNameAndId() (n string, id string) {
	value := d.getFieldValue("id")
	splits := strings.Split(value, " ")

	if len(splits) <= 0 {
		return "", ""
	}
	return splits[0], splits[1]
}
