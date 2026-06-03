package engines

import "conecto/core"


type Engine struct {	
	ConnectorRunnable  	ConnectorRunnable
	Transformer  		core.Transformer
	SinkCommiter 		SinkCommiter
}



