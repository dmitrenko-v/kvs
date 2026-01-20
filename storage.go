package main

import (
	"sync"
)

type KvsValue struct {
	dtype byte
	value []byte
}

type Kvs struct {
	mu      sync.Mutex
	storage map[string]*KvsValue
}

var kvs Kvs

func initStorage() {
	kvs.storage = make(map[string]*KvsValue)
}

func setHandler(args []*KvsValue) error {
	if len(args) != 2 {
		return ErrWrongSetArgsCount
	}

	key := args[0]
	value := args[1]

	if key.dtype != BulkStrSymbol {
		return ErrWrongKeyDtype
	}

	kvs.mu.Lock()
	kvs.storage[string(key.value)] = value
	kvs.mu.Unlock()

	return nil
}

func getHandler(args []*KvsValue) (value *KvsValue, err error) {
	if len(args) != 1 {
		return nil, ErrWrongGetArgsCount
	}

	key := args[0]

	if key.dtype != BulkStrSymbol {
		return nil, ErrWrongKeyDtype
	}

	res, ok := kvs.storage[string(key.value)]

	if !ok {
		return nil, nil
	}

	return res, nil
}

func deleteHandler(args []*KvsValue) error {
	if len(args) != 1 {
		return ErrWrongDelArgsCount
	}

	key := args[0]

	if key.dtype != BulkStrSymbol {
		return ErrWrongKeyDtype
	}

	kvs.mu.Lock()
	delete(kvs.storage, string(key.value))
	kvs.mu.Unlock()

	return nil
}
