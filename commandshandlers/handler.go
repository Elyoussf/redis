package handlers

import (
	parser "redis/RESP"
	"sync"
)

var Handlers map[string]func([]parser.Value) parser.Value

func ping(v []parser.Value) parser.Value {
	if len(v) == 0 {
		return parser.Value{Typ: "string", Str: "PONG"}
	}
	if len(v) > 1 {
		return parser.Value{Typ: "error", Bulk: "Error The number of arguement for ping is greater than 1 "}
	}

	return parser.Value{Typ: "string", Str: v[0].Bulk}

}

var SETs = make(map[string]string)
var SETMutex = sync.RWMutex{}

func set(args []parser.Value) parser.Value {
	if len(args) != 2 {
		return parser.Value{Typ: "error", Str: "Error The number of arguments for set is not equal to 2"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	SETMutex.Lock()
	defer SETMutex.Unlock()
	SETs[key] = value

	return parser.Value{Typ: "string", Str: "OK"}
}
func get(args []parser.Value) parser.Value {
	if len(args) != 1 {
		return parser.Value{Typ: "error", Bulk: "Error The number of arguments for get is not equal to 1"}
	}

	key := args[0].Bulk

	SETMutex.RLock()
	defer SETMutex.RUnlock()
	value, exists := SETs[key]
	if !exists {
		return parser.Value{Typ: "error", Str: "Key does not exist"}
	}

	return parser.Value{Typ: "string", Str: value}
}

func LoadHandlers() {
	Handlers = make(map[string]func([]parser.Value) parser.Value)
	Handlers["PING"] = ping
	Handlers["GET"] = get
	Handlers["SET"] = set
}
