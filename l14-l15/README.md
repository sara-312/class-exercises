# Transaction Serializability and 2PL

In this exercise you will be experimenting with Two-Phase Locking! This is a partner project! Please find a partner and record their name in a file called `names.txt` which you will submit alongside your code.

Groups of three are acceptable if necessary. 

## lockstore.go

In the exercise directory you have been provided a `lockstore.go` file which starts an RPC server.

This server has several endpoints:
- Get
- Put
- Lock
- Unlock

You will not need to modify `lockstore.go` in anyway!

## txn.go

Below is the starting code for your transaction manager. This code will run two transactions-- a writing transaction and a reading transaction.

The writing transaction must modify at least two keys in the KVStore with `Put`. 

The reading transaction should read all keys involved in the transaction with `Get` and verify if the output of concurrent transactions in serializable.

```golang
package main

import (
	"log"
	"net/rpc"
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

func Lock(client *rpc.Client, client_id string, key string) {
	args := &Args{Key:key, ClientID: client_id}
	var reply Reply
	reply.Success = false
	
	for reply.Success != true {
		err := client.Call("KVServer.Lock",args, &reply)
		if(err != nil) {
			log.Fatal(err)
		}
	}
}

func Unlock(client *rpc.Client, client_id string, key string) {
	args := &Args{Key:key, ClientID: client_id}
	var reply Reply
	err := client.Call("KVServer.Unlock",args, &reply)
	if(err != nil) {
		log.Fatal(err)
	}
}

func Get(client *rpc.Client, key string) string {
	args := &Args{Key:key}
	var reply Reply
	err := client.Call("KVServer.Get",args, &reply)
	if(err != nil) {
		log.Fatal(err)
	}
	return reply.Value
}

func Put(client *rpc.Client, key string, value string) {
	args := &Args{Key:key, Value: value}
	var reply Reply
	err := client.Call("KVServer.Put",args, &reply)
	if(err != nil) {
		log.Fatal(err)
	}
}


// This function should get/put values into the key value store
func run_txn(client *rpc.Client, client_id string) {

}

// This function should get values from the keyvalue store server 
// This should check if you get serializable results and print an error if the results are unserializable.
func run_read_txn(client *rpc.Client, client_id string) {

}

func main() {
	lockserver_address := "localhost:1234"
	client_id := "anna"

	client, err := rpc.DialHTTP("tcp", lockserver_address)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	for range(1000) {
		run_txn(client, client_id)
        run_read_txn(client, client_id)
	}
}
```

**Do not call Lock/Unlock YET**

Create a transaction in `run_txn` that does some combination of puts on certain variables. The choice of which variables and what values is up to you!
Talk to your partner as they also create a writing transaction that also performs a series of puts.

**Work together:** What are the serializable results from concurrently running your two transactions? 
You should design your two transactions so that they **overlap** in 2+ variables. It should be possible to observe unserializable results if the two transactions run concurrently (i.e. you should perform conflicting updates to the same variables).

For example:

```
T1:
x = 10
y = 20
```

```
T2:
x = 20
y = 40
```

Unserializable results from concurrently running these two transactions would be:

```
x=20 y=20
```
or
```
x=10 y=40
```

Complete `run_read_txn` to read relevant variables prints an error message if it observes unserializable results.

> [!IMPORTANT]
> - Complete the implementation of `run_txn` as described
> - Complete the implementation of `run_read_txn` to print an error message on observation of unserializable results

## Testing your Code

You will test your code with your partner. 

**Make sure you are connected to the same wifi**

Find the ip address of your current machine.

On Linux:
```
hostname -I
```

On Mac:
```
ipconfig getifaddr en0
```

On Windows:
```
ipconfig
```

- **One** partner will run the lockserver.go code from their machine.
- Put the address of this partner in `lockserver_address` in `txn.go`
- Put unique string in `client_id` variable in `txn.go`

The chosen partner should run:


```
go run lockserver.go 1234
```

Then, both partners should simultaneously run:

```
go run txn.go
```

Note: Testing may work better if you start transaction on the computer that is **not** running the lockserver just slightly earlier.
This way, there is a decreased a possibility that the quicker local transaction will occur entirely before the remote partner's transaction starts.

### Unserializable Results

Your `run_read_txn` function should flag unserializable results as we have no concurrency control mechanism implemented!

If you do not observe the error message, try again until you do.


> [!IMPORTANT]
> - Observe the printed error message from `run_read_txn`

### Implementing 2PL

You will now implement two-phase locking in both of your transactions.
Once your partner also completes the implementation re-run your tests.

> [!IMPORTANT]
> - Add the correct `Lock` and `Unlock` calls to your `run_txn` and `run_read_txn`. 
> You should no longer observe unserializable results.


### Deadlock

Purposefully trigger a case with **deadlock** across you and your partner's transactions.

- Hint: You may need to insert `time.Sleep` calls to make the deadlock more likely.
- Hint: Carefully plan out the order of locks acquired such that deadlock is possible across the transactions.

