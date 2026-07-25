package factories

import (
	"conecto/core/retry"
	"conecto/core/statestores"
)

type Conecto struct{
	connections  Connections
	random       retry.Random
	stateStore 	 statestores.StateStore
}