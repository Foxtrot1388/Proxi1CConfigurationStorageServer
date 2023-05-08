package tcpxml

import (
	"Proxi1CConfigurationStorageServer/internal/config"
	"Proxi1CConfigurationStorageServer/internal/entity"
	"context"
	"encoding/xml"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type WorkersConfiguration struct {
	pool      chan int
	eventchan chan entity.OneCEvents
	cancelctx context.CancelFunc
}

func GetConfiguration(ctx context.Context, cfg *config.Config, f func(context.Context, *config.Config, <-chan entity.OneCEvents)) *WorkersConfiguration {

	workcfg := WorkersConfiguration{
		eventchan: make(chan entity.OneCEvents, 20), // to cfg?
		pool:      make(chan int, cfg.NumAnalizeWorkers),
	}

	for i := 0; i < cfg.NumAnalizeWorkers; i++ {
		workcfg.pool <- i
	}

	newctx, cancel := context.WithCancel(ctx)
	workcfg.cancelctx = cancel
	go f(newctx, cfg, workcfg.eventchan)

	return &workcfg

}

func (w *WorkersConfiguration) FreeLockPoolAnalize(str string) {
	select {
	case id := <-w.pool:
		go func(tokenid int) {
			w.analyze(str)
			w.pool <- tokenid
		}(id)
	default:
		return
	}
}

func (w *WorkersConfiguration) Close() {
	close(w.pool)
	w.cancelctx()
	close(w.eventchan)
}

func (w *WorkersConfiguration) analyze(xmlreqest string) {

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
			w.eventchan <- result
		}
	default:
		return
	}

}
