package main

import (
	"Proxi1CConfigurationStorageServer/internal/config"
	"Proxi1CConfigurationStorageServer/internal/entity"
	"Proxi1CConfigurationStorageServer/internal/listenereventchan"
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

	wg := sync.WaitGroup{}
	wg.Add(1)

	f := func(ctx context.Context, ch <-chan entity.OneCEvents) {
		defer wg.Done()
		select {
		case val := <-ch:
			cast := val.(entity.CommitObject)
			if !(cast.Auth.User == "User.Test" && cast.Conf == "main" && cast.Params.Comment == "Comment for commit" && len(cast.Params.Changes.Value) == 3) {
				t.Fail()
			}
		case <-ctx.Done():
			t.Fail()
		}
	}

	workers := tcpxml.GetPoolWorkers(&config.Config{NumAnalizeWorkers: 1})
	defer workers.Close()
	workers.Analyze(string(xmlFile))

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()
	go f(ctx, workers.Eventchan)

	wg.Wait()

}

func TestEvent(t *testing.T) {

	var testevent entity.CommitObject
	testevent.Auth.User = "TestUser"
	testevent.Conf = "Main"
	testevent.Params.Comment = "Test comment"

	ctx, cancel := context.WithCancel(context.Background())

	eventWorker := listenereventchan.OScriptListener{
		entity.EventListen{
			Configuration: &entity.WorkConfiguration{
				Eventchan: make(chan entity.OneCEvents, 20),
			},
		},
	}
	go eventWorker.Listen(ctx, &config.Config{Scriptfile: map[string]string{"DevDepot_commitObjects": "CommitObject.os"}})

	eventWorker.Configuration.Eventchan <- testevent
	eventWorker.Configuration.Eventchan <- testevent

	time.Sleep(time.Duration(10 * time.Second))
	cancel()
	close(eventWorker.Configuration.Eventchan)

}
