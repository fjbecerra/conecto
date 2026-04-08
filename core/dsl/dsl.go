package transformer

import (
	"encoding/json"
	"fmt"
	"math"
	"os"

	"github.com/tidwall/gjson"
	"github.com/vjeantet/govaluate"
	// "github.com/vjeantet/govaluate"
)


type FieldConfig struct {
	Path     string      `json:"path"`
	Type     string      `json:"type"`
	Default  interface{} `json:"default"`
	Nullable bool        `json:"nullable"`
}

type SchemaConfig struct {
	Path        string                      `json:"path"`
    Is_array    bool                        `json:"is_array"`
	Fields      map[string]FieldConfig      `json:"fields"`
    Computed    map[string]ComputedConfig   `json:"computed"`
}

type ComputedConfig struct {
    Expr    string  `json:"expr"`
    Type    string  `json:"type"`
    Default interface{} `json:"default"`
}


//type Schema map[string]SchemaConfig

func Extract(jsonStr string, schema SchemaConfig) ([]map[string]interface{}) {
    var records []gjson.Result

    if schema.Path == "" {
        records = gjson.Parse(jsonStr).Array()
    } else {
        res := gjson.Get(jsonStr, schema.Path)

        if res.IsArray() {
            records = res.Array()
        }
    }  

    results := make([]map[string]interface{}, 0, len(records))

    for _, record := range records {
        row := make(map[string]interface{})

        // 1. Extract fields
        for field, fcfg := range schema.Fields {
            val := record.Get(fcfg.Path)

            if val.Exists() {
                row[field] = castGJSON(val, fcfg.Type)
            } else {
                row[field] = fcfg.Default
            }
        }
         // 2. Compute fields
        for field, ccfg := range schema.Computed {
            val := safeEval(
                ccfg.Expr,
                row,
                ccfg.Default,
                func(v interface{}) interface{} {
                    return castExpr(v, ccfg.Type)
                },
            )
            row[field] = val
        }
        results = append(results, row)    
          
    }

   
    return results
}

func safeEval(expr string, row map[string]interface{}, defaultVal interface{}, cast func(interface{}) interface{}) interface{} {
    result, err := evalExpr(expr, row)
    if err != nil {
        return defaultVal
    }

    switch v := result.(type) {
    case float64:
        if math.IsNaN(v) || math.IsInf(v, 0) {
            return defaultVal
        }
    }

    return cast(result)
}

func evalExpr(expression string, params map[string]interface{}) (interface{},error){
    expr, err := govaluate.NewEvaluableExpression(expression)
    if err != nil {
        return nil, err
    }

    return expr.Evaluate(params)

}

func castGJSON(v gjson.Result, typ string) interface{} {
    switch typ {
    case "int":
        return int(v.Int())
    case "float":
        return v.Float()
    case "string":
        return v.String()
    default:
        return v.Value() // fallback
    }
}

func castExpr(val interface{}, typ string) interface{} {
    switch typ {
    case "int":
        if f, ok := val.(float64); ok {
            return int64(f)
        }
        return val
    case "float":
        if f, ok := val.(float64); ok {
            return f
        }
        return val
    case "string":
        return fmt.Sprintf("%v", val)
    case "bool":
        if b, ok := val.(bool); ok {
            return b
        }
    }
    return val
}

func LoadSchema() SchemaConfig{
    data, err := os.ReadFile("schema.json")
    if err != nil {
        panic(err)
    }
    var schema SchemaConfig

    err = json.Unmarshal(data, &schema)
    if err != nil {
        panic(err)
    }
    return schema
}

func LoadResponse(json string) string {
    data, err := os.ReadFile(json)
    if err != nil {
        panic(err)
    }

    return string(data)
} 