package main

import (
	tcpxml "Proxi1CConfigurationStorageServer/internal/xml"
	"io/ioutil"
	"testing"
)

func TestSampleXML(t *testing.T) {

	xmlFile, err := ioutil.ReadFile("sample.xml")
	if err != nil {
		panic(err)
	}

	eventchan := make(chan interface{}, 20)
	go tcpxml.Analyze(string(xmlFile), eventchan)
	_ = <-eventchan

}
