package transforms

import (
	"conecto/core"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
)

type Fields map[string]struct{
	Path    string    
	Type    string
	Default interface{}
    Required bool
}

type Extractor struct {
    Fields Fields
}

func (e *Extractor) Extract(in json.RawMessage) (core.Record, error) {
    parsed := gjson.ParseBytes(in)

    if !parsed.Exists() || parsed.Type == gjson.Null {
        return core.Record{}, fmt.Errorf("no json to parse")
    }
    row := make(map[string]interface{})

    for field, cfg := range e.Fields {
        val := parsed.Get(cfg.Path)

        if val.Exists() {
            row[field] = castGJSON(val, cfg.Type)
        } else {
            row[field] = cfg.Default
        }
    }

    return row, nil
}

func castGJSON(val gjson.Result, typ string) interface{} {
    switch typ {
    case "int":
        return val.Int()
    case "float":
        return val.Float()
    case "bool":
        return val.Bool()
    default:
        return val.String()
    }
}