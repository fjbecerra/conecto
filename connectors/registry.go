package connectors

import (
	"errors"
)


var ErrConnectorNotFound =
	errors.New("connector not found")



type Registry struct {

	connectors map[string]Connector
}



func NewRegistry(items ...Connector) *Registry {

	r := &Registry{

		connectors:
		make(map[string]Connector),
	}

	for _,connector := range items {
		r.Register(
			connector,
		)
	}

	return r
}



func (r *Registry) Register(connector Connector) {

	r.connectors[connector.Name()] = connector
}



func (r *Registry) Get(name string) (Connector,error){

	connector,ok := r.connectors[name]
	if !ok {
		return nil,
			ErrConnectorNotFound
	}
	return connector,nil
}