package main

import (
	"Proxi1CConfigurationStorageServer/internal/entity"
	"Proxi1CConfigurationStorageServer/internal/event"
	tcpxml "Proxi1CConfigurationStorageServer/internal/xml"
	"io/ioutil"
	"testing"
)

func TestSampleXML(t *testing.T) {

	xmlFile, err := ioutil.ReadFile("sample.xml")
	if err != nil {
		panic(err)
	}

	eventchan := make(chan entity.OneCEvents, 20)
	tcpxml.Analyze(string(xmlFile), eventchan)
	val := <-eventchan
	event.DoEvent([]entity.OneCEvents{0: val})
	close(eventchan)

}
