package graphql

import "conecto/connectors"

type GraphQLCursorInjector struct {
    Variables map[string]any
}

func (g *GraphQLCursorInjector) Inject(c *connectors.PageCursor) error {

    if c == nil {
        return nil
    }

    g.Variables["cursor"] = c.Value

    return nil
}