package tcpxml

import (
	"bytes"
	"encoding/xml"
)

type CommitObject struct {
	Conf string
	Auth struct {
		User string `xml:"user,attr"`
	} `xml:"auth"`
	Params struct {
		Changes struct {
			Value []struct {
				Second struct {
					Super struct {
						Name struct {
							Value string `xml:"value,attr"`
						} `xml:"name"`
					} `xml:"_super"`
				} `xml:"second"`
			} `xml:"value"`
		} `xml:"changes"`
		Comment string `xml:"comment"`
	} `xml:"params"`
}

func Analyze(xmlreqest string, eventchan chan<- interface{}) {

	xmlData := bytes.NewBufferString(xmlreqest)
	d := xml.NewDecoder(xmlData)

	_, err := d.Token()
	if err != nil {
		return
	}

	_, err = d.Token()
	if err != nil {
		return
	}

	t, err := d.Token()
	if err != nil {
		return
	}

	switch se := t.(type) {
	case xml.StartElement:
		if se.Name.Local == "call" && len(se.Attr) == 4 && se.Attr[2].Value == "DevDepot_commitObjects" {
			var result CommitObject
			d.DecodeElement(&result, &se)
			result.Conf = se.Attr[1].Value
			eventchan <- result
		}
	default:
		return
	}

}
