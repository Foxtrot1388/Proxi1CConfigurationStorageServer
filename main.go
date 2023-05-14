package main

import (
	"Proxi1CConfigurationStorageServer/internal/config"
	"Proxi1CConfigurationStorageServer/internal/listenereventchan"
	tcpxml "Proxi1CConfigurationStorageServer/internal/xml"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
)

var configname *string = flag.String("configname", "app.yaml", "target config")

type AnalyzeWork interface {
	Analyze(string)
	Close()
}

func GetConfiguration(ctx context.Context, cfg *config.Config) (AnalyzeWork, func()) {

	eventchan := make(chan interface{}, 20)
	workcfg := tcpxml.GetPoolWorkers(cfg, eventchan)
	eventlistener := listenereventchan.GetListener(eventchan)

	newctx, cancelctx := context.WithCancel(ctx)
	go eventlistener.Listen(newctx, cfg)

	return workcfg, func() {
		workcfg.Close()
		cancelctx()
	}

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

	portlistener, err := net.Listen("tcp", ":"+cfg.ListenPort)
	if err != nil {
		panic(err)
	}

	workcfg, close := GetConfiguration(context.Background(), cfg)
	defer close()

	for {
		if conin, err := portlistener.Accept(); err == nil {

			go func() {
				conout, err := net.Dial("tcp", cfg.Host+":"+cfg.Port)
				if err != nil {
					panic(err)
				}
				done := make(chan struct{})
				go readwritetotcp(conin, conout, done, infologhost, workcfg)
				go readwritetotcp(conout, conin, done, infologlocal, nil)
				<-done
				<-done
			}()

		} else {
			panic(err)
		}
	}

}

func readwritetotcp(conin net.Conn, connout net.Conn, done chan<- struct{}, logdebug *log.Logger, workcfg AnalyzeWork) {

	readbyte := make([]byte, 10240)
	for {

		n, err := conin.Read(readbyte)
		if err != nil {
			break
		}

		if n > 0 {

			if logdebug != nil {
				logdebug.Println(string(readbyte[:n]))
			}

			connout.Write(readbyte[:n])

			if workcfg != nil {
				workcfg.Analyze(string(readbyte[:n]))
			}

		}
	}

	done <- struct{}{}
}
