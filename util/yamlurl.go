package util

import "net/url"

// url.URL does not adhear to the yaml.Unmarshaler interface
// Need to write a wrapper to url.URL and impliment UnmarshalYAML and MarshalYAML
// https://povilasv.me/yaml-url-parsing/
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
