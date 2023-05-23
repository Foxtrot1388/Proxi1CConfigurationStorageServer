package listenereventchan

import (
	"Proxi1CConfigurationStorageServer/internal/config"
	"context"
	"encoding/json"
	"os/exec"
	"time"
)

type oneCEvents interface {
	GetCompactEvent() interface{}
	Append(map[string]aggevents)
}

type aggevents []oneCEvents

type OScriptListener struct {
	eventchan <-chan interface{}
}

func GetListener(eventchan <-chan interface{}) *OScriptListener {
	return &OScriptListener{eventchan: eventchan}
}

func (e *OScriptListener) Listen(ctx context.Context, cfg *config.Config) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			rawevent := e.readNextPart()
			e.doEvent(cfg, rawevent)
			time.Sleep(time.Duration(5 * time.Minute))
		}
	}
}

func (e *OScriptListener) readNextPart() []oneCEvents {
	var rawevent []oneCEvents
	for {
		select {
		case val, ok := <-e.eventchan:
			if !ok {
				return rawevent
			}
			rawevent = append(rawevent, val.(oneCEvents))
		default:
			return rawevent
		}
	}
}

func (e *OScriptListener) doEvent(cfg *config.Config, val []oneCEvents) {

	aggevent := make(map[string]aggevents, len(cfg.Scriptfile))
	for i := 0; i < len(val); i++ {
		val[i].Append(aggevent)
	}

	for k, v := range aggevent {
		cmd := exec.Command("oscript", cfg.Scriptfile[k], v.getJSON())
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}

}

func (ob aggevents) getJSON() string {
	dat := make([]interface{}, len(ob))
	for i := 0; i < len(ob); i++ {
		dat[i] = ob[i].GetCompactEvent()
	}
	js, err := json.Marshal(dat)
	if err != nil {
		panic(err)
	}
	return string(js)
}
