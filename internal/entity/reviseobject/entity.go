package reviseobject

const (
	AttrReviseObjectEvent         = 2
	AttrReviseObjectConfiguration = 1
)

type oneCEvents interface {
	GetCompactEvent() interface{}
	Append(map[string]aggevents)
}

type aggevents []oneCEvents

type ReviseObject struct {
	Conf string
	Auth struct {
		User string `xml:"user,attr"`
	} `xml:"auth"`
	Params struct {
		Objects struct {
			Value []struct {
				Second struct {
					Super struct {
						Name struct {
							Value string `xml:"value,attr"`
						} `xml:"name"`
					} `xml:"_super"`
				} `xml:"second"`
			} `xml:"value"`
		} `xml:"objects"`
	} `xml:"params"`
}

func (com ReviseObject) GetCompactEvent() interface{} {
	dat := make(map[string]interface{})
	objects := make([]string, len(com.Params.Objects.Value))
	dat["user"] = com.Auth.User
	dat["configuration"] = com.Conf
	for i := 0; i < len(com.Params.Objects.Value); i++ {
		objects[i] = com.Params.Objects.Value[i].Second.Super.Name.Value
	}
	dat["objects"] = objects
	return dat
}

func (val ReviseObject) Append(collection map[string]aggevents) {
	collection["DevDepot_reviseDevObjects"] = append(collection["DevDepot_reviseDevObjects"], val)
}
