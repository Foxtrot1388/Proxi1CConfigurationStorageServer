package main

import (
	"Proxi1CConfigurationStorageServer/internal/config"
	"Proxi1CConfigurationStorageServer/internal/entity"
	tcpxml "Proxi1CConfigurationStorageServer/internal/xml"
	"context"
	"io/ioutil"
	"sync"
	"testing"
	"time"
)

func TestSampleXML(t *testing.T) {

	xmlFile, err := ioutil.ReadFile("sample.xml")
	if err != nil {
		panic(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	wg := sync.WaitGroup{}

	f := func(ch chan entity.OneCEvents) {
		defer wg.Done()
		select {
		case val := <-ch:
			if js, err := val.GetJSON(); err != nil || js != "{\"comment\":\"Comment for commit\",\"configuration\":\"main\",\"objects\":[\"Object1\",\"Object2\",\"Object3\"],\"user\":\"User.Test\"}" {
				t.Fail()
			}
		case <-ctx.Done():
			t.Fail()
		}
	}

	wg.Add(1)
	workers := tcpxml.GetConfiguration(&config.Config{NumAnalizeWorkers: 1}, f)
	defer workers.Close()
	workers.FreeLockPoolAnalize(string(xmlFile))
	wg.Wait()

}
