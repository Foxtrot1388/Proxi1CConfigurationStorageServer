package tcp

import (
	"fmt"
	"log"
	"net"
)

type Listener struct {
	portlistener net.Listener
}

type analyzerWork interface {
	Analyze(string)
}

func New(listenPort string) (Listener, error) {

	fmt.Println("Use tcp storage...")

	portlistener, err := net.Listen("tcp", ":"+listenPort)
	if err != nil {
		return Listener{}, err
	}

	return Listener{portlistener: portlistener}, nil
	
}

func (l *Listener) Do(host, port string, workcfg interface{}, infologlocal, infologhost *log.Logger) {

	for {
		if conin, err := l.portlistener.Accept(); err == nil {

			go func() {
				conout, err := net.Dial("tcp", host+":"+port)
				if err != nil {
					panic(err)
				}
				done := make(chan struct{})
				go l.readwritetotcp(conin, conout, done, infologhost, workcfg.(analyzerWork))
				go l.readwritetotcp(conout, conin, done, infologlocal, nil)
				<-done
				<-done
			}()

		} else {
			panic(err)
		}
	}

}

func (l *Listener) readwritetotcp(conin net.Conn, connout net.Conn, done chan<- struct{}, logdebug *log.Logger, workcfg analyzerWork) {

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
