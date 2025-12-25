package conv

import (
	"bytes"
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

func YAMLtoJSON(data []byte) ([]byte, error) {
	var m any
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return json.Marshal(m)
}

func JSONtoYAML(object any) ([]byte, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	var node yaml.Node
	if err = yaml.Unmarshal(data, &node); err != nil {
		return nil, err
	}

	clearNodeStyle(&node)

	var b bytes.Buffer
	encoder := yaml.NewEncoder(&b)
	encoder.SetIndent(2)
	if err = encoder.Encode(&node); err != nil {
		return nil, fmt.Errorf("encoder encode: %w", err)
	}
	return b.Bytes(), nil
}

func clearNodeStyle(node *yaml.Node) {
	if node == nil {
		return
	}
	node.Style = 0
	for _, n := range node.Content {
		clearNodeStyle(n)
	}
}
