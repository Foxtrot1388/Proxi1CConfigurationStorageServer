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

func GetListener(eventchan <-chan entity.OneCEvents) *ScriptListener {
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
			time.Sleep(time.Duration(20 * time.Second))
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

func (e *ScriptListener) doEvent(cfg *config.Config, val []entity.OneCEvents) {

	aggevent := make(map[string]entity.Aggevents, len(cfg.Scriptfile))
	for i := 0; i < len(val); i++ {
		val[i].Append(aggevent)
	}

	for k, v := range aggevent {
		var err error
		switch {
		case strings.HasSuffix(cfg.Scriptfile[k], ".os"):
			cmd := exec.Command("oscript", cfg.Scriptfile[k], getJSON(v))
			err = cmd.Run()
		case strings.HasSuffix(cfg.Scriptfile[k], ".sbsl"):
			cmd := exec.Command("executor", "-s "+cfg.Scriptfile[k], getJSON(v))
			err = cmd.Run()
		default:
			err = errors.New("Unknown script type!")
		}
		if err != nil {
			panic(err)
		}
	}

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
