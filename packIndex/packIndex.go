// Code generated by xgen. DO NOT EDIT.

package packIndex

import (
	"encoding/xml"
)

// SemanticVersionType ...
type SemanticVersionType string

// RestrictedString ...
type RestrictedString string

// VidxType ...
type VidxType struct {
	UrlAttr    string `xml:"url,attr"`
	VendorAttr string `xml:"vendor,attr"`
	DateAttr   string `xml:"date,attr,omitempty"`
}

// PdscType ...
type PdscType struct {
	UrlAttr         string `xml:"url,attr"`
	VendorAttr      string `xml:"vendor,attr"`
	NameAttr        string `xml:"name,attr"`
	VersionAttr     string `xml:"version,attr"`
	DateAttr        string `xml:"date,attr,omitempty"`
	DeprecatedAttr  string `xml:"deprecated,attr,omitempty"`
	ReplacementAttr string `xml:"replacement,attr,omitempty"`
	SizeAttr        uint32 `xml:"size,attr,omitempty"`
}

// PindexType ...
type PindexType struct {
	Pdsc []*PdscType `xml:"pdsc"`
}

// VindexType ...
type VindexType struct {
	Pidx []*VidxType `xml:"pidx"`
}

// Index ...
type Index struct {
	XMLName           xml.Name      `xml:"index"`
	SchemaVersionAttr string        `xml:"schemaVersion,attr"`
	Vendor            string        `xml:"vendor"`
	Url               string        `xml:"url"`
	Timestamp         string        `xml:"timestamp"`
	Pindex            []*PindexType `xml:"pindex"`
	Vindex            []*VindexType `xml:"vindex"`
}
