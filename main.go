package main

import (
	"Proxi1CConfigurationStorageServer/internal/config"
	"Proxi1CConfigurationStorageServer/internal/entity"
	storagehttp "Proxi1CConfigurationStorageServer/internal/input/http"
	storagetcp "Proxi1CConfigurationStorageServer/internal/input/tcp"
	"Proxi1CConfigurationStorageServer/internal/listenereventchan"
	tcpxml "Proxi1CConfigurationStorageServer/internal/xml"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
)

var configname *string = flag.String("configname", "app.yaml", "target config")

type analyzerWork interface {
	Analyze(string)
	Close()
}

type listener interface {
	Do(host, port string, workcfg interface{}, infologlocal, infologhost *log.Logger)
}

func main() {

	fmt.Println("Launching server...")

	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	cfg := config.Get(configname)
	var infologlocal, infologhost *log.Logger
	if cfg.Debug {
		infologlocal = log.New(os.Stdout, "to localhost: ", log.LstdFlags)
		infologhost = log.New(os.Stdout, "to "+cfg.Host+": ", log.LstdFlags)
	}

	workcfg, closecfg := newConfiguration(context.Background(), cfg)
	defer closecfg()

	newlistener := newListener(cfg.ListenPort, cfg.Type)
	newlistener.Do(cfg.Host, cfg.Port, workcfg, infologlocal, infologhost)

}

func newConfiguration(ctx context.Context, cfg *config.Config) (analyzerWork, func()) {

	eventchan := make(chan entity.OneCEvents, 20)
	workcfg := tcpxml.NewPoolWorkers(cfg, eventchan)
	eventlistener := listenereventchan.NewListener(eventchan)

	newctx, cancelctx := context.WithCancel(ctx)
	go eventlistener.Listen(newctx, cfg)

	return workcfg, func() {
		workcfg.Close()
		cancelctx()
	}

}

func newListener(listenport, typeinput string) listener {

	if typeinput == "tcp" {
		newlistener, err := storagetcp.New(listenport)
		if err != nil {
			panic(err)
		}
		return &newlistener
	} else if typeinput == "http" {
		newlistener, err := storagehttp.New(listenport)
		if err != nil {
			panic(err)
		}
		return &newlistener
	} else {
		log.Fatal("Wrong input storage!")
		return nil
	}

}
