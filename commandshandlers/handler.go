package handlers

import (
	parser "redis/RESP"
	"sync"
)

var SETs = make(map[string]string)
var SETMutex = sync.RWMutex{}
var HSETs = make(map[string]map[string]string)
var HSETMutex = sync.RWMutex{}
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

func hset(args []parser.Value) parser.Value {
	if len(args) != 3 {
		return parser.Value{Typ: "error", Str: "Error The number of arguments for hset is not equal to 3"}
	}

	key := args[0].Bulk
	field := args[1].Bulk
	value := args[2].Bulk

	HSETMutex.Lock()
	defer HSETMutex.Unlock()

	if _, exists := HSETs[key]; !exists {
		HSETs[key] = make(map[string]string)
	}
	HSETs[key][field] = value

	return parser.Value{Typ: "string", Str: "OK"}
}

func hget(args []parser.Value) parser.Value {
	if len(args) != 2 {
		return parser.Value{Typ: "error", Str: "Error The number of arguments for hget is not equal to 2"}
	}

	key := args[0].Bulk
	field := args[1].Bulk

	HSETMutex.RLock()
	defer HSETMutex.RUnlock()

	fields, exists := HSETs[key]
	if !exists {
		return parser.Value{Typ: "error", Str: "Key does not exist"}
	}

	value, exists := fields[field]
	if !exists {
		return parser.Value{Typ: "error", Str: "Field does not exist"}
	}

	return parser.Value{Typ: "string", Str: value}
}

func hgetall(args []parser.Value) parser.Value {
	if len(args) != 1 {
		return parser.Value{Typ: "error", Str: "Error The number of arguments for hgetall is not equal to 1"}
	}

	key := args[0].Bulk

	HSETMutex.RLock()
	defer HSETMutex.RUnlock()

	fields, exists := HSETs[key]
	if !exists {
		return parser.Value{Typ: "error", Str: "Key does not exist"}
	}

	result := parser.Value{Typ: "array", Array: make([]parser.Value, 0, len(fields)*2)}
	for field, value := range fields {
		result.Array = append(result.Array, parser.Value{Typ: "string", Str: field}, parser.Value{Typ: "string", Str: value})
	}

	return result
}

func LoadHandlers() {
	Handlers = make(map[string]func([]parser.Value) parser.Value)
	Handlers["PING"] = ping
	Handlers["GET"] = get
	Handlers["SET"] = set
	Handlers["HSET"] = hset
	Handlers["HGET"] = hget
	Handlers["HGETALL"] = hgetall
}
