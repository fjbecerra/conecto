package core

import (
	"context"
	"fmt"
	"math"
	"github.com/vjeantet/govaluate"
	
)

type Compute struct {
	Schema SchemaConfig
}
type Computer interface{
	Compute(ctx context.Context, in <-chan Record) (<-chan Record, <-chan error)
}

func NewCompute(schema SchemaConfig) Computer{
	return &Compute{schema}
}

func (c *Compute) Compute(ctx context.Context, in <-chan Record) (<-chan Record, <-chan error){
	return StreamMap(ctx, in, func(row Record) (Record, error) {
	
			for field, ccfg := range c.Schema.Computed {
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

        	return row, nil
    	
	})

   
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