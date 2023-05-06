package event

import (
	"Proxi1CConfigurationStorageServer/internal/entity"
	"fmt"
)

func EventListener(ch chan entity.OneCEvents) {
	// TODO
	for {
		val, ok := <-ch
		if !ok {
			break
		}
		DoEvent([]entity.OneCEvents{0: val})
	}
}

func DoEvent(val []entity.OneCEvents) {
	// TODO
	for i := 0; i < len(val); i++ {
		fmt.Println(val[i].GetJSON())
	}
}
