package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)

type Args struct {
	Key      string
	Value    string
	ClientID string
}

type Reply struct {
	Success bool
	Value   string
}

type KVServer struct {
	store map[string]string

	// Protects the store and lock table.
	lock sync.Mutex

	// Maps each key to the client currently holding its lock.
	locks map[string]string
}

// Lock acquires an exclusive lock on a key.
func (t *KVServer) Lock(args *Args, reply *Reply) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	log.Println(t.locks)

	owner, exists := t.locks[args.Key]

	if exists && owner != args.ClientID {
		reply.Success = false
		return nil
	}

	t.locks[args.Key] = args.ClientID

	reply.Success = true
	return nil
}

// Unlock releases a lock on a key.
func (t *KVServer) Unlock(args *Args, reply *Reply) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	owner, exists := t.locks[args.Key]

	log.Println("UNLOCKING")

	// The key is not locked.
	if !exists {
		reply.Success = false
		return nil
	}

	// Only the owner can unlock the key.
	if owner != args.ClientID {
		reply.Success = false
		return nil
	}

	delete(t.locks, args.Key)
	log.Println(t.locks)

	reply.Success = true
	return nil
}

func (t *KVServer) Get(args *Args, reply *Reply) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	value, exists := t.store[args.Key]

	reply.Success = exists
	reply.Value = value

	return nil
}

func (t *KVServer) Put(args *Args, reply *Reply) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.store[args.Key] = args.Value

	reply.Success = true
	reply.Value = ""

	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run server.go <port>")
	}

	server := &KVServer{
		store: make(map[string]string),
		locks: make(map[string]string),
	}

	rpc.Register(server)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":"+os.Args[1])
	if err != nil {
		log.Fatal("listen error:", err)
	}

	fmt.Println("KV server listening on port", os.Args[1])

	err = http.Serve(l, nil)
	if err != nil {
		log.Fatal("server error:", err)
	}
}

