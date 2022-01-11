package util

import "net/url"

type Yamlurl struct {
	*url.URL
}

func (j *Yamlurl) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}
	url, err := url.Parse(s)
	j.URL = url
	return err
}

func (j Yamlurl) MarshalYAML() (interface{}, error) {
	return j.String(), nil
}
