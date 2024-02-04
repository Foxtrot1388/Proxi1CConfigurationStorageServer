package listenereventchan

import (
	"Proxi1CConfigurationStorageServer/internal/config"
	"Proxi1CConfigurationStorageServer/internal/entity"
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
	"time"
)

type ScriptListener struct {
	eventchan <-chan entity.OneCEvents
}

func NewListener(eventchan <-chan entity.OneCEvents) *ScriptListener {
	return &ScriptListener{eventchan: eventchan}
}

func (e *ScriptListener) Listen(ctx context.Context, cfg *config.Config) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			rawevent := e.readNextPart()
			e.doEvents(cfg, rawevent)
		}
	}
}

func (e *ScriptListener) readNextPart() []entity.OneCEvents {
	var rawevent []entity.OneCEvents
	for {
		select {
		case val, ok := <-e.eventchan:
			if !ok {
				return rawevent
			}
			rawevent = append(rawevent, val)
		default:
			return rawevent
		}
	}
}

func (e *ScriptListener) doEvents(cfg *config.Config, val []entity.OneCEvents) {

	aggevent := make(map[string]entity.Aggevents, len(cfg.Scriptfile))
	for i := 0; i < len(val); i++ {
		val[i].Append(aggevent)
	}

	for k, v := range aggevent {
		err := e.doEvent(cfg.Scriptfile[k], getJSON(v))
		if err != nil {
			panic(err)
		}
	}

}

func (e *ScriptListener) doEvent(script string, arg string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	switch {
	case strings.HasSuffix(script, ".os"):
		cmd := exec.CommandContext(ctx, "oscript", script, arg)
		err = cmd.Run()
	case strings.HasSuffix(script, ".sbsl"):
		cmd := exec.CommandContext(ctx, "executor", "-s "+script, arg)
		err = cmd.Run()
	default:
		err = errors.New("Unknown script type!")
	}
	return err
}

func getJSON(ob entity.Aggevents) string {
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
