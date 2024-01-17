
# Setpartition_unrank

A Go library to lexicographicaly unrank set partitions


## Authors

- Antoine GENITRINI, CNRS, LIP6
- Amaury CURIEL, CNRS, LIP6


## Installation

To import the library you need first to create a Go project. In your terminal enter the following command : 
```bash
go mod init your_project
```
Then, you need to get the setpartition_unrank library. In your project folder, where is stored the go.mod file enter this command :

```bash
go get github.com/AMAURYCU/setpartition_unrank
```
You are now ready to use our library
## License

[GPL-3.0]("https://github.com/AMAURYCU/setpartition_unrank/blob/main/LICENSE")
## Demo

First of all, ensure to have followed the installation step.
In the folder where your go.mod file is, enter the following command : 
```bash
go get github.com/AMAURYCU/setpartition_unrank
```
then you can import the library puting the following line in the head of your Go project
```go
import "github.com/AMAURYCU/setpartition_unrank/XXXXX"
```
where 
```
XXXXX
```
should be replaced with one of 
```
parallelunranking //to run the efficient parallel algorithm
precalcul //to run the algorithm with precomputation step
statistic //to make some statistics on the library
```
An example of program that lists all set partitions of the set [|1,10|] in 5 blocks : 
```go
package main

import (
	"fmt"
	"math/big"

	"github.com/AMAURYCU/setpartition_unrank/parallelunranking"
)

func main() {
	c := parallelunranking.Stirling2Columns(10, 5).Col1[10]
	c.Sub(&c, big.NewInt(1))
	for k2 := big.NewInt(0); k2.Cmp(&c) < 1; k2.Add(k2, big.NewInt(1)) {
        // 10 stand for [|1,10|], 5 for 5 blocks, *k2 to iterate over sets partitions
        // and for to use S3V5 (you should always use it)
		fmt.Println(parallelunranking.UnrankDicho(10, 5, *k2, 4), k2)
	}
}
//output : 
/*.
.
.
[[1 10] [2 9] [3 7 8] [4] [5 6]] 42515
[[1 10] [2 9] [3 7 8] [4 5] [6]] 42516
[[1 10] [2 9] [3 7 8] [4 6] [5]] 42517
[[1 10] [2 9] [3 8] [4] [5 6 7]] 42518
[[1 10] [2 9] [3 8] [4 5] [6 7]] 42519
[[1 10] [2 9] [3 8] [4 5 6] [7]] 42520
[[1 10] [2 9] [3 8] [4 5 7] [6]] 42521
[[1 10] [2 9] [3 8] [4 6] [5 7]] 42522
[[1 10] [2 9] [3 8] [4 6 7] [5]] 42523
[[1 10] [2 9] [3 8] [4 7] [5 6]] 42524
*/
```


## Documentation

[Documentation](https://pkg.go.dev/github.com/AMAURYCU/setpartition_unrank)


## Related

This project is related to the implementation of our paper : 
[Lexicographic unranking algorithms for the Twelvefold Way](link)

