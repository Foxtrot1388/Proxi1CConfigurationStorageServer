package main

import (
	"Proxi1CConfigurationStorageServer/app/internal/config"
	"flag"
	"log"
	"net"
	"os"
	"runtime"
)
import "fmt"

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

	for {
		if conin, err := portlistener.Accept(); err == nil {

			go func() {
				conout, err := net.Dial("tcp", cfg.Host+":"+cfg.Port)
				if err != nil {
					panic(err)
				}
				done := make(chan struct{})
				go readwritetotcp(conin, conout, done, infologhost)
				go readwritetotcp(conout, conin, done, infologlocal)
				<-done
				<-done
			}()

		} else {
			panic(err)
		}
	}

}

func readwritetotcp(conin net.Conn, connout net.Conn, done chan<- struct{}, logdebug *log.Logger) {
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
		}
	}
	done <- struct{}{}
}
