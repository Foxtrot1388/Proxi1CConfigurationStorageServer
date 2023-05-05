package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"

	"Proxi1CConfigurationStorageServer/internal/config"
	event "Proxi1CConfigurationStorageServer/internal/event"
	tcpxml "Proxi1CConfigurationStorageServer/internal/xml"
)

var configname *string = flag.String("configname", "app.yaml", "target config")

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

	pool := make(chan int, cfg.NumAnalizeWorkers)
	defer close(pool)
	for i := 0; i < cfg.NumAnalizeWorkers; i++ {
		pool <- i
	}

	eventchan := make(chan interface{}, 20)
	defer close(eventchan)
	go event.EventListener(eventchan)

	for {
		if conin, err := portlistener.Accept(); err == nil {

			go func() {
				conout, err := net.Dial("tcp", cfg.Host+":"+cfg.Port)
				if err != nil {
					panic(err)
				}
				done := make(chan struct{})
				go readwritetotcp(conin, conout, done, infologhost, pool, eventchan)
				go readwritetotcp(conout, conin, done, infologlocal, nil, nil)
				<-done
				<-done
			}()

		} else {
			panic(err)
		}
	}

}

func readwritetotcp(conin net.Conn, connout net.Conn, done chan<- struct{}, logdebug *log.Logger, poolworkers chan int, eventchan chan<- interface{}) {

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

			if poolworkers != nil && eventchan != nil {
				select {
				case id := <-poolworkers:
					go func(req string, tokenid int) {
						tcpxml.Analyze(req, eventchan)
						poolworkers <- tokenid
					}(string(readbyte[:n]), id)
				default:
					break
				}

			}

		}
	}

	done <- struct{}{}
}
