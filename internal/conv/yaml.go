package conv

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

func YAMLtoJSON(data []byte) ([]byte, error) {
	var m any
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return json.Marshal(m)
}
