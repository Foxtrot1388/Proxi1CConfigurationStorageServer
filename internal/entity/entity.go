package entity

const (
	AttrCommitObjectEvent         = 2
	AttrCommitObjectConfiguration = 1
	AttrReviseObjectEvent         = 2
	AttrReviseObjectConfiguration = 1
)

type OneCEvents interface {
	GetCompactEvent() interface{}
	Append(map[string]Aggevents)
}

type Aggevents []OneCEvents

type CommitObject struct {
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

func (val CommitObject) Append(collection map[string]Aggevents) {
	collection["DevDepot_commitObjects"] = append(collection["DevDepot_commitObjects"], val)
}

type ReviseObject struct {
	Conf string
	Auth struct {
		User string `xml:"user,attr"`
	} `xml:"auth"`
	Params struct {
		Revise struct {
			Value string `xml:"value,attr"`
		} `xml:"revise"`
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
	dat["revise"] = (com.Params.Revise.Value == "true")
	for i := 0; i < len(com.Params.Objects.Value); i++ {
		objects[i] = com.Params.Objects.Value[i].Second.Super.Name.Value
	}
	dat["objects"] = objects
	return dat
}

func (val ReviseObject) Append(collection map[string]Aggevents) {
	collection["DevDepot_reviseDevObjects"] = append(collection["DevDepot_reviseDevObjects"], val)
}
