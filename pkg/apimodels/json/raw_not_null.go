package json

import "encoding/json"

// RawMessageNotNull is an alias for json.RawMessage that must be initialized
// and must not contain the JSON null value.
type RawMessageNotNull = json.RawMessage
