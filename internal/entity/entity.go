package entity

import "encoding/json"

type OneCEvents interface {
	GetJSON() (string, error)
}

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

func (com CommitObject) GetJSON() (string, error) {
	dat := make(map[string]interface{})
	objects := make([]string, len(com.Params.Changes.Value))
	dat["user"] = com.Auth.User
	dat["comment"] = com.Params.Comment
	dat["configuration"] = com.Conf
	for i := 0; i < len(com.Params.Changes.Value); i++ {
		objects[i] = com.Params.Changes.Value[i].Second.Super.Name.Value
	}
	dat["objects"] = objects
	str, err := json.Marshal(dat)
	return string(str), err
}
