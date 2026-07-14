# Vector Timestamps

You will be implementing a vector timestamp onto the following code.

You will create a file `process.go` representing a single process who is implementing a vector clock.
You will also eventually create a tester for this process code in a file `process_test.go`.

# Process

A process contains 4 methods:
- `Start(N int, PID int)` -- which sets up a process with N other processes in the system.
- `Internal()` -- Which should increment the vector clock correctly for this process.
- `Send() -> Vector Timestamp` -- Which should update the process's vector timestamp and return a vector timestamp.
- `Receive(Vector Timestamp)` -- Which should receive a vector timestamp and perform the correct logic for updating the process's timestamp

Copy the following code into process.go and complete the above methods.

```golang
package main


type Process struct {
	clock []int
	pid int
	N int
}

func (p *Process) Start(N int, PID int)  {
	
}

func (p *Process) Internal() {
	
}

func (p *Process) Send() []int {

}

func (p *Process) Receive(ts []int) {

}

func Compare(ts1 []int, ts2 []int) int {

}
```

Notice that this code is built around calling a method defined on a struct type `Process`.

Read about how go handles methods here:
- https://gobyexample.com/methods


> [!IMPORTANT]
> Complete the implementation of the four vector timestamp functions, updating the fields of the process struct as needed.



# Comparison

We also want the ability to correctly compare two vector timestamps.

In `process.go` implement the `Compare` function:

```golang

func Compare(ts1 []int, ts2 []int) int {

}
```

This function should return:
- an int < 0 if ts1 < ts2
- an int > 0 if ts1 > ts2
- 0 if ts1 < > ts2 or ts1 == ts2.

> [!IMPORTANT]
> Complete the implementation of the Compare vector timestamp functions

# Testing Code

To test your `process.go` copy the following code into `process_test.go`.

You will need to complete the implementation of `TestProcess`.

To understand how the testing package works, read: https://gobyexample.com/testing-and-benchmarking 

```golang
package main

import (
	"testing"
	"reflect"
)

/*

Run the following test:

p1: .   .   .        .
         \          /
p2:       .  .     /
              \   /
p3:  . .       . .

*/

func TestProcesses(t *testing.T) {

	var p1, p2, p3 Process

	p1.Start(3,0)
	p2.Start(3,1)
	p3.Start(3,2)

	// TODO: fill in calls to p.Internal, p.Send(), p.Receive() correctly
    // The execution should match the comment above.

	expected_p1 := []int{4,2,4}
	expected_p2 := []int{2,2,0}
	expected_p3 := []int{2,2,4} 

	got1 := p1.clock
	if !reflect.DeepEqual(got1, expected_p1) {
		t.Fatalf(
			"expected %v\ngot      %v", expected_p1, got1,
		)
	}
	got2 := p2.clock
	if !reflect.DeepEqual(got2, expected_p2) {
		t.Fatalf(
			"expected %v\ngot      %v", expected_p2, got2,
		)
	}
	got3 := p3.clock
	if !reflect.DeepEqual(got3, expected_p3) {
		t.Fatalf(
			"expected %v\ngot      %v", expected_p3, got3,
		)
	}

}


func TestCompare3(t *testing.T) {
	t1 := []int{1,2,3}
	t2 := []int{2,1,2}

	expected := 0
	result := Compare(t1,t2)
	if(result != expected) {
		t.Errorf("Expected %v, Got %v\n", expected, result)
	}

	t1 = []int{1,2,3}
	t2 = []int{0,1,0}

	expected = 1
	result = Compare(t1,t2)
	if(result <= 0) {
		t.Errorf("Expected %v, Got %v\n", expected, result)
	}

	expected = -1
	result = Compare(t2,t1)
	if(result >= 0) {
		t.Errorf("Expected %v, Got %v\n", expected, result)
	}
}

func TestCompare2(t *testing.T) {
	if Compare([]int{1, 2}, []int{2, 3}) >= 0 {
		t.Error("expected first timestamp to be less than second")
	}

	if Compare([]int{3, 4}, []int{2, 3}) <= 0 {
		t.Error("expected first timestamp to be greater than second")
	}

	if Compare([]int{2, 1}, []int{1, 2}) != 0 {
		t.Error("expected concurrent timestamps")
	}
}

```

> [!IMPORTANT]
> Complete the implementation of the `TestProcesses` function by correctly calling Start, Internal, Send, Receive to match the execution graph below.

```
p1: .   .   .        .
         \          /
p2:       .  .     /
              \   /
p3:  . .       . .
```



You may need to run `go mod init l7-l8` to set-up the go module file.

When you are finished, `go test` should pass.

```shell
go test
PASS
ok      l7-l8    0.182s
```
