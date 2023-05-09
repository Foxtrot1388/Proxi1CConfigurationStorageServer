package listenereventchan

import (
	"Proxi1CConfigurationStorageServer/internal/config"
	"Proxi1CConfigurationStorageServer/internal/entity"
	"context"
	"encoding/json"
	"os/exec"
	"time"
)

type aggevents []entity.OneCEvents

type OScriptListener struct {
	entity.EventListen
}

func (e *OScriptListener) Listen(ctx context.Context, cfg *config.Config) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			rawevent := e.readNextPart(e.Configuration.Eventchan)
			e.doEvent(cfg, rawevent)
			time.Sleep(time.Duration(5 * time.Minute))
		}
	}
}

func (e *OScriptListener) readNextPart(ch <-chan entity.OneCEvents) []entity.OneCEvents {
	var rawevent []entity.OneCEvents
	for {
		select {
		case val, ok := <-ch:
			if !ok {
				return rawevent
			}
			rawevent = append(rawevent, val)
		default:
			return rawevent
		}
	}
}

func (e *OScriptListener) doEvent(cfg *config.Config, val []entity.OneCEvents) {

	aggevent := make(map[string]aggevents, len(cfg.Scriptfile))
	for i := 0; i < len(val); i++ {
		switch val[i].(type) {
		case entity.CommitObject:
			aggevent["DevDepot_commitObjects"] = append(aggevent["DevDepot_commitObjects"], val[i])
		}
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
