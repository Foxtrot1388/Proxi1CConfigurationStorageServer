package listenereventchan

import (
	"Proxi1CConfigurationStorageServer/internal/config"
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
	"time"
)

type oneCEvents interface {
	GetCompactEvent() interface{}
	Append(map[string]aggevents)
}

type aggevents []oneCEvents

type ScriptListener struct {
	eventchan <-chan interface{}
}

func GetListener(eventchan <-chan interface{}) *ScriptListener {
	return &ScriptListener{eventchan: eventchan}
}

func (e *ScriptListener) Listen(ctx context.Context, cfg *config.Config) {
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

func (e *ScriptListener) readNextPart() []oneCEvents {
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

func (e *ScriptListener) doEvent(cfg *config.Config, val []oneCEvents) {

	aggevent := make(map[string]aggevents, len(cfg.Scriptfile))
	for i := 0; i < len(val); i++ {
		val[i].Append(aggevent)
	}

	for k, v := range aggevent {
		var err error
		switch {
		case strings.HasSuffix(cfg.Scriptfile[k], ".os"):
			cmd := exec.Command("oscript", cfg.Scriptfile[k], v.getJSON())
			err = cmd.Run()
		case strings.HasSuffix(cfg.Scriptfile[k], ".sbsl"):
			cmd := exec.Command("executor", "-s "+cfg.Scriptfile[k], v.getJSON())
			err = cmd.Run()
		default:
			err = errors.New("Unknown script type!")
		}
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
