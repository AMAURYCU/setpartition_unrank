package main

import (
	"flag"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/AMAURYCU/setpartition_unrank/parallelunranking"
	"github.com/AMAURYCU/setpartition_unrank/precalcul"
	"github.com/AMAURYCU/setpartition_unrank/statistic"
)

func main() {

	operation := flag.String("operation", "", "Specify operation: A, R or G")
	mode := flag.String("mode", "", "Specify mode: P or S")
	flag.Parse()

	if *operation == "" {
		printUsageAndExit()
	}

	switch *operation {
	case "A":
		if *mode == "" {
			printUsageAndExit()
		}
		handleOperationA(*mode, flag.Args())
	case "R":
		if *mode == "" {
			printUsageAndExit()
		}
		handleOperationR(*mode, flag.Args())
	case "G":
		handleOperationG(flag.Args())
	default:
		printUsageAndExit()
	}

}

func handleOperationA(mode string, args []string) {

	if len(args) != 2 {
		fmt.Println("Error: Operation A requires exactly 2 arguments.")
		printUsageAndExit()
	}

	n, err1 := strconv.Atoi(args[0])
	k, err2 := strconv.Atoi(args[1])

	if err1 != nil || err2 != nil {
		fmt.Println("Error: Arguments for Operation A must be numeric.")
		printUsageAndExit()
	}

	switch mode {
	case "P":
		c := parallelunranking.Stirling2Columns(n, k).Col1[n]
		c.Sub(&c, big.NewInt(1))
		for k2 := big.NewInt(0); k2.Cmp(&c) < 1; k2.Add(k2, big.NewInt(1)) {
			fmt.Println(parallelunranking.UnrankDicho(n, k, *k2, 4), k2)
		}
	case "S":
		precalcul.Init()
		c := parallelunranking.Stirling2Columns(n, k).Col1[n]
		c.Sub(&c, big.NewInt(1))
		precalcul.StirlingMatrix = statistic.StirlingTriangle(n, k)
		for k2 := big.NewInt(0); k2.Cmp(&c) < 1; k2.Add(k2, big.NewInt(1)) {
			fmt.Println(precalcul.UnrankDichoPre(n, k, *k2, 0), k2)
		}
	default:
		fmt.Printf("Error: Invalid mode %s for operation A.\n", mode)
		printUsageAndExit()
	}

}

func handleOperationR(mode string, args []string) {

	if len(args) != 2 {
		fmt.Println("Error: Operation R requires exactly 2 arguments.")
		printUsageAndExit()
	}

	n, err1 := strconv.Atoi(args[0])
	k, err2 := strconv.Atoi(args[1])

	if err1 != nil || err2 != nil {
		fmt.Println("Error: Arguments for Operation R must be numeric.")
		printUsageAndExit()
	}

	sg := rand.NewSource(time.Now().UnixNano())
	rg := rand.New(sg)
	var r big.Int

	switch mode {
	case "P":
		c := parallelunranking.Stirling2Columns(n, k).Col1[n]
		c.Sub(&c, big.NewInt(1))
		r.Rand(rg, &c)
		parallelunranking.UnrankDicho(n, k, r, 0)
		fmt.Println("temps calcul prev col", parallelunranking.ListToString(parallelunranking.TimePreviousColumn))
		fmt.Println("-----------------------------")
		fmt.Println("k", statistic.ListToString(parallelunranking.TimePreviousColumnWithK))
	case "S":
		precalcul.Init()
		precalcul.StirlingMatrix = statistic.StirlingTriangle(n, k)
		r.Rand(rg, precalcul.StirlingMatrix[n][k])
		fmt.Println(precalcul.UnrankDichoPre(n, k, r, 0), &r)
	default:
		fmt.Printf("Error: Invalid mode %s for operation R.\n", mode)
		printUsageAndExit()
	}

}

func handleOperationG(args []string) {

	if len(args) != 3 {
		fmt.Println("Error: Operation G requires exactly 3 arguments.")
		printUsageAndExit()
	}

	n, err1 := strconv.Atoi(args[0])
	k, err2 := strconv.Atoi(args[1])
	r, err3 := strconv.Atoi(args[2])

	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Println("Error: Arguments for Operation G must be numeric.")
		printUsageAndExit()
	}
	precalcul.Init()
	statistic.Stat(n, k, r, true)
	a, b, c := statistic.Graph3d(n, 10, 10, r)
	file, err := os.Create("graph3d.g")
	fmt.Println(err)
	defer file.Close()
	_, v := file.WriteString(statistic.PrintMatrix(a) + "// \n" + statistic.ListToString(b) + "// \n" + statistic.ListToString(c))
	fmt.Println(v)

}

func printUsageAndExit() {
	fmt.Println("Usage: program_name -operation [A/R/G] -mode [P/S] [arguments]")
	fmt.Println("Operations:")
	fmt.Println("  A: to generate all partitions of n1 in n2 non-empty disjoints subsets - Requires 2 numeric arguments")
	fmt.Println("  R: to randomly pickup one partition of n1 in n2 non-empty disjoints subsets - Requires 2 numeric arguments")
	fmt.Println("  G: to have an overview of the performance of the algorithm partitionning n1 in n2 non-empty disjoints subsets with n3 points - Requires 3 numeric arguments")
	fmt.Println("Modes:")
	fmt.Println("  P: parallel")
	fmt.Println("  S: sequential")
	flag.PrintDefaults()
	os.Exit(1)
}
