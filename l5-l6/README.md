# Mini Map Reduce

MapReduce is a framework which allows for a task to be easily and automatically parallelized onto a cluster of worker machines.

To use the MapReduce framework, programmers must break up their task into a `Map` function and a `Reduce` function.

You will be implementing a "mini" version of MapReduce, writing a Map function and a Reduce Function as well as a simple "driver" coordinator program. Unlike real MapReduce this will all run in a single process on a single machine. However, we will parallelize the code with goroutines!

## The Task: Max Temperature Per Year

Your task is to find the maximum temperature recorded each year.

The dataset you have is organized into 20 files, each has many logged temperatures across many days.

You could write a program that goes over each of these files-- however, we want to do this search in a highly-parallelized fashion following the MapReduce paradigm.

- Your Map function will take in a filepath and file contents and output key-value pairs of any found (year, temperature) readings.
- Your Reduce function will take in a year and a list of temperatures within that year and return the max temperature.
- Your coordinator will call Map/Reduce on the correct inputs (in parallel) and print the result.

## The Map Function

The map function should take in a file and output the following pair`(year, temperature)` for each found temperature.

Place this function in `worker.go`.

```go
func Map(key string, value string) []KVPair {

	output := make([]KVPair,0)

    // TODO: This loop iterates over each line of the "value" string
    // You will want to parse out the date and temperature from each line and add it to the "output" slice
	for _, line := range strings.Split(strings.TrimSuffix(value, "\n"), "\n") {
    	fmt.Println(line)
	}

	return output
}
```

> [!IMPORTANT]
> Complete the implementation of `Map`.

To test our map function, set `coordinator.go` to the following:

```go
package main

import (
	"fmt"
	"os"
)

type KVPair struct {
	key string
	value string
}

func main() {
	path := "data/Moscow.txt"
    input, err := os.ReadFile(path)
    if err != nil {
        panic(err)
    }
    map_out := Map(path, string(input))
    fmt.Println(map_out)
}
```

You can run your code with the command:
`go run coordinator.go worker.go`

The output should resemble:
```
[{2020 14.2} {2022 21.7} {2018 -3.1} {2023 -7.3} {2017 -0.7} {2021 22.6} {2019 16.5} {2023 8.0} {2017 19.2} {2022 7.4} {2023 1.3} {2018 -5.1} {2021 13.5} {2020 4.4} {2019 2.9} {2024 2.8} {2025 4.6} {2017 14.8} {2020 16.9} {2018 -4.9} {2025 14.6} {2023 12.3} {2017 13.1} {2020 9.7} {2019 4.2} {2018 23.1} {2016 0.1} {2016 12.5} {2022 11.9} {2023 -7.2} {2019 3.2} {2017 -0.7} {2020 2.3} {2026 0.5} ... 
```

## The Reduce Function

The reduce function takes in all pairs `(year, temperature)` for a given year and outputs the maximum temperature.

Place this function in `worker.go`.

```go
func Reduce(key string, value []string) float64 {

    // Converting from a string to float may be useful
	// val,err := strconv.ParseFloat(INPUT, 64)

}
```

> [!IMPORTANT]
> Complete the implementation of `Reduce`.

To test Reduce, add the following code to `coordinator.go`:

```go
res := Reduce("2024", []string{"12.1", "-10.5", "32.1", "30.6"})
fmt.Println(res)
```

You should see the output:
```
32.1
```

## The "Coordinator"

Paste the following starter code in `coordinator.go`:

```go

package main

import (
	"fmt"
	"os"
	"sync"
)

type KVPair struct {
	key string
	value string
}

func main() {
	cities := []string{"Tokyo", "Delhi", "Shanghai", "Sao_Paulo", "Mexico_City", "Cairo", "Mumbai", "Beijing", "Dhaka", "Osaka", "New_York", "Karachi", "Buenos_Aires", "Istanbul", "Kolkata", "Lagos", "Moscow", "London", "Paris", "Los_Angeles"}

    //Reads in input file for each city and calls Map function
	ch := make(chan KVPair)
	var wg sync.WaitGroup
	for _,city := range(cities) {
		wg.Add(1)
		go func(city string) {
			defer wg.Done()
			path := "data/"+ city+".txt"
			input, err := os.ReadFile(path)
			if err != nil {
				panic(err)
			}
			map_out := Map(path, string(input))
			for _,item := range(map_out) {
				ch<-item
			}
		}(city)
	}

    // Goroutine waits to close the channel
	go func() {
		wg.Wait()
		close(ch)
	}()

    // range over channel repeatedly reads from channel until it is closed
	kv_pairs := make(map[string][]string)
	for item := range(ch) {
        // TODO: correctly populate the map kv_pairs with the items read in on the channel "ch"
	}

	// TODO: add calling Reduce tasks
}

```

Typically, the coordinator has a big job of tracking tasks and responding to workers requests for work and detecting failed jobs. However, our simple implementation the coordinator will just call the map function in several go routines, wait, and then call the reduce function in go routines.

Currently:
- The coordinator opens the 20 files and call "Map" on each file.
- Each output item from a Map task is sent on the channel `ch`
- The coordinator receives on this channel, but does nothing with the output.

You will add:
- Storing the output items into a `map[string][]string`. This is a mapping strings (keys) to slices of strings (list of associated values). 
    - We will map a year (ie `2024`) to the list of temperatures occurring that year (ie `{"12.1", "-10.5", "32.1", "30.6"}`)
    - **Note: Storing the intermediate data on the coordinator is different from real MapReduce where intermediate key/value pairs are stored on worker machines and requested by reduce tasks as needed**
- Read more about Go's built-in map type to get started:
    - https://gobyexample.com/maps 


The coordinator should then wait, and call Reduce on each intermediate key and all associated values.

> [!IMPORTANT]
> Complete the implementation of `coordinator.go` by calling Reduce on the correct intermediate input.

The final output for each reduce task can be printed or logged.

You final output should resemble the following:
```
Highest Temp in 2026 was 38.3
Highest Temp in 2016 was 32.5
Highest Temp in 2020 was 36.9
Highest Temp in 2025 was 37.5
Highest Temp in 2023 was 36.5
Highest Temp in 2019 was 37.5
Highest Temp in 2017 was 38
Highest Temp in 2022 was 38.3
Highest Temp in 2024 was 39.2
Highest Temp in 2018 was 36.8
Highest Temp in 2021 was 37.7
```

Note, the exact order the years are reported should differ each run. This is a signal the task was correctly parallelized!