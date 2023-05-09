package tcpxml

import (
	"Proxi1CConfigurationStorageServer/internal/config"
	"Proxi1CConfigurationStorageServer/internal/entity"
	"encoding/xml"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type WorkersConfiguration struct {
	pool      chan int
	Eventchan chan entity.OneCEvents
}

func GetPoolWorkers(cfg *config.Config) *WorkersConfiguration {

	workcfg := WorkersConfiguration{
		Eventchan: make(chan entity.OneCEvents, 20), // to cfg?
		pool:      make(chan int, cfg.NumAnalizeWorkers),
	}

	for i := 0; i < cfg.NumAnalizeWorkers; i++ {
		workcfg.pool <- i
	}

	return &workcfg
}

func (w *WorkersConfiguration) Analyze(str string) {
	select {
	case id := <-w.pool:
		go func(tokenid int) {
			w.analyzeXML(str)
			w.pool <- tokenid
		}(id)
	default:
		return
	}
}

func (w *WorkersConfiguration) Close() {
	close(w.pool)
	close(w.Eventchan)
}

func (w *WorkersConfiguration) analyzeXML(xmlreqest string) {

	firstindex := strings.Index(xmlreqest, "<?xml")
	if firstindex == -1 {
		return
	}

	lastindex := strings.LastIndex(xmlreqest, ">")
	if lastindex == -1 {
		return
	}

	decoder := charmap.Windows1251.NewDecoder()
	reader := decoder.Reader(strings.NewReader(xmlreqest[firstindex:lastindex]))
	d := xml.NewDecoder(reader)

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
		if se.Name.Local == "call" && len(se.Attr) == 4 && se.Attr[entity.AttrCommitObjectEvent].Value == "DevDepot_commitObjects" {
			var result entity.CommitObject
			d.DecodeElement(&result, &se)
			result.Conf = se.Attr[entity.AttrCommitObjectConfiguration].Value
			w.Eventchan <- result
		}
	default:
		return
	}

}
