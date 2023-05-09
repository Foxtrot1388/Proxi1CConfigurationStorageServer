package entity

import (
	"Proxi1CConfigurationStorageServer/internal/config"
	"context"
)

type OneCEvents interface {
	GetCompactEvent() interface{}
}

const (
	AttrCommitObjectEvent         = 2
	AttrCommitObjectConfiguration = 1
)

type AnalyzeWork interface {
	Analyze(string)
	Close()
}

type WorkConfiguration struct {
	AnalyzeWork
	Eventchan chan OneCEvents
}

type EventListener interface {
	Listen(context.Context, *config.Config)
}

type EventListen struct {
	EventListener
	Configuration *WorkConfiguration
}

type CommitObject struct {
	OneCEvents
	Conf string
	Auth struct {
		User string `xml:"user,attr"`
	} `xml:"auth"`
	Params struct {
		Changes struct {
			Value []struct {
				Second struct {
					Super struct {
						Name struct {
							Value string `xml:"value,attr"`
						} `xml:"name"`
					} `xml:"_super"`
				} `xml:"second"`
			} `xml:"value"`
		} `xml:"changes"`
		Comment string `xml:"comment"`
	} `xml:"params"`
}

func (com CommitObject) GetCompactEvent() interface{} {
	dat := make(map[string]interface{})
	objects := make([]string, len(com.Params.Changes.Value))
	dat["user"] = com.Auth.User
	dat["comment"] = com.Params.Comment
	dat["configuration"] = com.Conf
	for i := 0; i < len(com.Params.Changes.Value); i++ {
		objects[i] = com.Params.Changes.Value[i].Second.Super.Name.Value
	}
	dat["objects"] = objects
	return dat
}
