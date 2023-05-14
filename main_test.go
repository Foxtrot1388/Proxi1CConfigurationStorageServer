package main

import (
	"Proxi1CConfigurationStorageServer/internal/config"
	"Proxi1CConfigurationStorageServer/internal/entity/commitobject"
	"Proxi1CConfigurationStorageServer/internal/entity/reviseobject"
	"Proxi1CConfigurationStorageServer/internal/listenereventchan"
	tcpxml "Proxi1CConfigurationStorageServer/internal/xml"
	"context"
	"io/ioutil"
	"sync"
	"testing"
	"time"
)

func TestCommitSampleXML(t *testing.T) {

	xmlFile, err := ioutil.ReadFile("samplecommit.xml")
	if err != nil {
		panic(err)
	}

	eventchan := make(chan interface{}, 20)
	wg := sync.WaitGroup{}
	wg.Add(1)

	f := func(ctx context.Context, ch <-chan interface{}) {
		defer wg.Done()
		select {
		case val := <-ch:
			cast := val.(commitobject.CommitObject)
			if !(cast.Auth.User == "User.Test" && cast.Conf == "main" && cast.Params.Comment == "Comment for commit" && len(cast.Params.Changes.Value) == 3) {
				t.Fail()
			}
		case <-ctx.Done():
			t.Fail()
		}
	}

	workers := tcpxml.GetPoolWorkers(&config.Config{NumAnalizeWorkers: 1}, eventchan)
	defer workers.Close()
	workers.Analyze(string(xmlFile))

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()
	go f(ctx, eventchan)

	wg.Wait()

}

func TestReviseSampleXML(t *testing.T) {

	xmlFile, err := ioutil.ReadFile("samplerevise.xml")
	if err != nil {
		panic(err)
	}

	eventchan := make(chan interface{}, 20)
	wg := sync.WaitGroup{}
	wg.Add(1)

	f := func(ctx context.Context, ch <-chan interface{}) {
		defer wg.Done()
		select {
		case val := <-ch:
			cast := val.(reviseobject.ReviseObject)
			if !(cast.Auth.User == "User.Test" && cast.Conf == "main" && len(cast.Params.Objects.Value) == 2) {
				t.Fail()
			}
		case <-ctx.Done():
			t.Fail()
		}
	}

	workers := tcpxml.GetPoolWorkers(&config.Config{NumAnalizeWorkers: 1}, eventchan)
	defer workers.Close()
	workers.Analyze(string(xmlFile))

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()
	go f(ctx, eventchan)

	wg.Wait()

}

func TestEvent(t *testing.T) {

	var testevent commitobject.CommitObject
	testevent.Auth.User = "TestUser"
	testevent.Conf = "Main"
	testevent.Params.Comment = "Test comment"

	eventchan := make(chan interface{}, 20)
	ctx, cancel := context.WithCancel(context.Background())

	eventWorker := listenereventchan.GetListener(eventchan)
	go eventWorker.Listen(ctx, &config.Config{Scriptfile: map[string]string{"DevDepot_commitObjects": "CommitObject.os"}})

	eventchan <- testevent
	eventchan <- testevent

	time.Sleep(time.Duration(10 * time.Second))
	cancel()
	close(eventchan)

}
