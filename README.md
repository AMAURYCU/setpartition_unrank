
# Setpartition_unrank

A Go library to lexicographicaly unrank set partitions


## Authors

- Amaury CURIEL, Sorbonne Univerisity, LIP6, Paris
- Antoine GENITRINI, Sorbonne University, LIP6, Paris



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


## Executable application

The git repo has a ```main.go```file that can be executed entering this command : 

```go run main.go -operation [A/R/G] -mode [P/S] [arguments]```
where : 
```
Operations:
  A: to generate all partitions of n1 in n2 non-empty disjoints subsets - Requires 2 numeric arguments
  R: to randomly pickup one partition of n1 in n2 non-empty disjoints subsets - Requires 2 numeric arguments
  G: to have an overview of the performance of the algorithm 
  partitionning n1 in n2 non-empty disjoints subsets with n3 points - Requires 3 numeric arguments
Modes:
  P: parallel
  S: sequential
  -mode string
    	Specify mode: P or S
  -operation string
    	Specify operation: A, R or G
```
For example : 
```
go run main.go -operation R -mode P 50 5
```

will output:
```
[[1 5 11 32 34 39 43 44 45 48] [2 3 7 8 12 14 25 28 29 37] 
[4 20 22 24 26 35 38 47] [6 9 10 13 19 33 36 40] 
[15 16 17 18 21 23 27 30 31 41 42 46 49 50]]
```

warning : the ```G``` operation does not require ```-mode``` arguments
## Related

This project is related to the implementation of our paper : 
[Lexicographic unranking algorithms for the Twelvefold Way](link)

