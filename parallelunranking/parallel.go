// Package parallelunranking provides functions to unrank set partitions lexicographicaly
//
// The main function of this package is the function UnrankDicho that generates
// a set partition in exactly k part in the lexicographic order.
// Our package provides a formula with 4 levels of optimization and paralellism to avoid precomputation step
package parallelunranking

import (
	"math/big"
	"time"

	"github.com/AMAURYCU/setpartition_unrank/types"
)

var WaitingTime int64

var StirlingColumn0 []big.Int
var StirlingColumn1 []big.Int

var TimePreviousColumn []int64
var TimePreviousColumnWithK []int64

var vs3 = [5](func(n, k int, swap bool, d int) big.Int){S3v1, S3v2, S3v3, S3v4, S3v5}
var TimeTotal int64

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

/*
Direct implementation of the s3 formula (read the readme section Related propositon 7 for more details)
n - number of elements in the set
k - number of blocks desired
swap - flag
d - last element of the unranked prefix
u - length of the prefix
*/
func S3v1(n, k int, swap bool, d int) big.Int {
	if d < 0 {
		return *big.NewInt(0)
	}
	if d == 0 && !(k-1 <= n && k-1 >= 0) {
		return *big.NewInt(0)
	}
	var res *big.Int
	if !swap {
		res = new(big.Int).Set(&StirlingColumn0[n])
	} else {
		res = new(big.Int).Set(&StirlingColumn1[n])
	}
	u := 0
	b := big.NewInt(1)
	for u < min(n-d, n-k+1) {
		b.Mul(b, big.NewInt(int64(n-d-u)))
		u++
		b.Div(b, big.NewInt(int64(u)))

		if !swap {
			res.Add(res, new(big.Int).Mul(&StirlingColumn0[n-u], b))
		} else {
			res.Add(res, new(big.Int).Mul(&StirlingColumn1[n-u], b))
		}
	}
	return *res
}

/*
	    Implementation of the s3 formula (read the readme section Related proposition 7 for more details)
		taking in consideration the symmetry of the binomials coefficients to improve the
		computation time of the formula.
		n - number of elements in the set
		k - number of blocks desired
		swap - flag
		d - last element of the unranked prefix
		u - length of the prefix
*/
func S3v2(n, k int, swap bool, d int) big.Int {
	if d < 0 {
		return *big.NewInt(0)
	}
	if d == 0 && !(k-1 <= n && k-1 >= 0) {
		return *big.NewInt(0)
	}
	var res *big.Int
	if !swap {
		res = new(big.Int).Set(&StirlingColumn0[n])
	} else {
		res = new(big.Int).Set(&StirlingColumn1[n])
	}

	if d >= k-1 {
		if !swap {
			res.Add(res, &StirlingColumn0[d])
		} else {
			res.Add(res, &StirlingColumn1[d])
		}
	}
	u := 0
	b := big.NewInt(1)
	for u < min((n-d)/2, n-k+1) {
		b.Mul(b, big.NewInt(int64(n-d-u)))
		u++
		b.Div(b, big.NewInt(int64(u)))

		if (d+u >= k-1) && (u < (n-d)/2 || (u == (n-d)/2 && (n-d)%2 == 1)) {
			var tmp *big.Int
			if !swap {
				tmp = new(big.Int).Add(&StirlingColumn0[n-u], &StirlingColumn0[d+u])
			} else {
				tmp = new(big.Int).Add(&StirlingColumn1[n-u], &StirlingColumn1[d+u])
			}
			tmp.Mul(tmp, b)
			res.Add(res, tmp)
		} else {
			if !swap {
				res.Add(res, new(big.Int).Mul(&StirlingColumn0[n-u], b))
			} else {
				res.Add(res, new(big.Int).Mul(&StirlingColumn1[n-u], b))
			}
		}
	}
	return *res
}

/*
	    Direct implementation of the s3 formula (read the readme section Related proposition 9 for more details)
		n - number of elements in the set
		k - number of blocks desired
		swap - flag
		d - last element of the unranked prefix
		u - length of the prefix
*/
func S3v3(n, k int, swap bool, d int) big.Int {
	if d < 0 {
		return *big.NewInt(0)
	}
	if d == 0 {
		if k-1 <= n && k-1 >= 0 {
			if !swap {
				return StirlingColumn1[n+1]
			} else {
				return StirlingColumn0[n+1]
			}
		}
		return *big.NewInt(0)
	}
	var res *big.Int
	if !swap {
		res = new(big.Int).Set(&StirlingColumn1[n+1])
	} else {
		res = new(big.Int).Set(&StirlingColumn0[n+1])
	}

	u := 0
	b := big.NewInt(1)
	pm1 := big.NewInt(-1)
	for u < d {
		b.Mul(b, big.NewInt(int64(d-u)))
		u++
		b.Div(b, big.NewInt(int64(u)))
		var tmp *big.Int
		if !swap {
			tmp = new(big.Int).Mul(&StirlingColumn1[n+1-u], b)
		} else {
			tmp = new(big.Int).Mul(&StirlingColumn0[n+1-u], b)
		}
		tmp.Mul(tmp, pm1)
		res.Add(res, tmp)
		pm1.Neg(pm1)
	}
	return *res
}

/*
	    Direct implementation of the s3 formula (read the readme section Related proposition 9 for more details)
		taking in acount the symmetry of the binomial coefficients to improve computation time.
		n - number of elements in the set
		k - number of blocks desired
		swap - flag
		d - last element of the unranked prefix
		u - length of the prefix
*/
func S3v4(n, k int, swap bool, d int) big.Int {
	if d < 0 {
		return *big.NewInt(0)
	}
	if d == 0 {
		if k-1 <= n && k-1 >= 0 {
			if !swap {
				return StirlingColumn1[n+1]
			} else {
				return StirlingColumn0[n+1]
			}
		}
		return *big.NewInt(0)
	}
	var res big.Int
	if d%2 == 1 {
		if !swap {
			res = *new(big.Int).Sub(&StirlingColumn1[n+1], &StirlingColumn1[n+1-d])
		} else {
			res = *new(big.Int).Sub(&StirlingColumn0[n+1], &StirlingColumn0[n+1-d])
		}
	} else {
		if !swap {
			res = *new(big.Int).Add(&StirlingColumn1[n+1], &StirlingColumn1[n+1-d])
		} else {
			res = *new(big.Int).Add(&StirlingColumn0[n+1], &StirlingColumn0[n+1-d])
		}
	}
	u := 0
	b := big.NewInt(1)
	pm1 := big.NewInt(-1)
	for u < min(d/2, n-k+1) {
		b.Mul(b, big.NewInt(int64(d-u)))
		u++
		b.Div(b, big.NewInt(int64(u)))

		if u < d/2 || (u == d/2 && d%2 == 1) {
			if d%2 == 1 {
				var tmp *big.Int
				if !swap {
					tmp = new(big.Int).Sub(&StirlingColumn1[n+1-u], &StirlingColumn1[n+1-d+u])
				} else {
					tmp = new(big.Int).Sub(&StirlingColumn0[n+1-u], &StirlingColumn0[n+1-d+u])
				}
				tmp.Mul(tmp, b)
				tmp.Mul(tmp, pm1)
				res.Add(&res, tmp)
			} else {
				var tmp *big.Int
				if !swap {
					tmp = new(big.Int).Add(&StirlingColumn1[n+1-u], &StirlingColumn1[n+1-d+u])
				} else {
					tmp = new(big.Int).Add(&StirlingColumn0[n+1-u], &StirlingColumn0[n+1-d+u])
				}
				tmp.Mul(tmp, b)
				tmp.Mul(tmp, pm1)
				res.Add(&res, tmp)
			}
		} else {
			var tmp *big.Int
			if !swap {
				tmp = new(big.Int).Mul(&StirlingColumn1[n+1-u], b)
			} else {
				tmp = new(big.Int).Mul(&StirlingColumn0[n+1-u], b)
			}
			tmp.Mul(tmp, pm1)
			res.Add(&res, tmp)
		}
		pm1.Neg(pm1)
	}
	return res
}

/*
	    Take the best of s3v4 and s3v2
		n - number of elements in the set
		k - number of blocks desired
		swap - flag
		d - last element of the unranked prefix
*/
func S3v5(n, k int, swap bool, d int) big.Int {
	if 2*d < n {
		return S3v4(n, k, swap, d)
	} else {
		return S3v2(n, k, swap, d)
	}
}

//  Unrank set partition lexicographicaly.
/*
This function is the main function of the package and takes 4 arguments as parameters :
- n : int, the cardinal of the set to be partitioned.
- k : int, the number of desired blocks in the result.
- rank : big.Int, the rank of the desired set partition in the lexicographical order.
- whichS3: [|0,4|], the desired version of S3 formula to use (4 is the fastest, 0 is the slowest).

For the time complexity, read the README section Related, theorem 12 of the article
Example usage:

    result := parallelunranking.UnrankDicho(5, 3, *big.NewInt(10), 4)
    fmt.Println(result) // Output: [[1 2 3] [4] [5]]
*/
func UnrankDicho(n, k int, rank big.Int, whichS3 int) [][]int {
	listK := make([]int64, 0)
	blockSizes := make([]int64, 0)
	blockTime := make([]int64, 0)
	TimeTotal = 0
	if k == 1 {
		res := make([][]int, 0)
		tmp := make([]int, 0)
		for d := 1; d <= n; d++ {
			tmp = append(tmp, d)
		}
		res = append(res, tmp)
		return res
	}

	n0 := n
	res := make([][]int, 0)
	r := *new(big.Int).Set(&rank)
	chanRes := make(chan []big.Int)
	startTime := time.Now().UnixMilli()
	couple := *Stirling2Columns(n, k)
	endTime := time.Now().UnixMilli()
	WaitingTime = endTime - startTime
	StirlingColumn0 = couple.Col0[:n]
	StirlingColumn1 = couple.Col1

	go computePreviousColumn(StirlingColumn0, n-1, k-1, chanRes)

	swap := false

	for k > 1 {
		//fmt.Println(k)
		listK = append(listK, int64(k))
		var startTime int64
		var endTime int64
		startTime = time.Now().UnixMicro()
		block, acc := optimizedBlockDicho(n, k, swap, r, whichS3)
		endTime = time.Now().UnixMicro()
		blockTime = append(blockTime, endTime-startTime)
		res = append(res, block)
		r.Sub(&r, &acc)
		n -= len(block)
		k--

		blockSizes = append(blockSizes, int64(len(block)))
		if !swap {
			startTime = time.Now().UnixMicro()
			StirlingColumn1 = <-chanRes
			endTime = time.Now().UnixMicro()
		} else {
			startTime = time.Now().UnixMicro()
			StirlingColumn0 = <-chanRes
			endTime = time.Now().UnixMicro()
		}
		TimeTotal += endTime - startTime

		if k > 1 {

			if !swap {
				go computePreviousColumn(StirlingColumn1, n-1, k-1, chanRes)
			} else {
				go computePreviousColumn(StirlingColumn0, n-1, k-1, chanRes)
			}
		}
		swap = !swap
	}
	res = append(res, make([]int, n))
	res = lexicographicPermutationUnrank(n0, res)
	return res
}

func lexicographicPermutationUnrank(n int, Pos [][]int) [][]int {
	L := make([]int, n)
	for i := 0; i < n; i++ {
		L[i] = i + 1
	}
	var P [][]int
	for b := 0; b < len(Pos); b++ {
		p := []int{}
		for _, i := range Pos[b] {
			p = append(p, L[i])
			L = append(L[:i], L[i+1:]...)
		}
		P = append(P, p)
	}
	return P
}

func computePreviousColumn(column []big.Int, n, k int, resultChan chan []big.Int) {
	TimePreviousColumnWithK = append(TimePreviousColumnWithK, int64(k))
	startTime := time.Now().UnixMicro()
	if k == 1 {
		res := make([]big.Int, n+1)
		res[0] = *big.NewInt(1)
		resultChan <- res
		return
	}
	if k == 2 {
		res := make([]big.Int, n+1)
		res[0] = *big.NewInt(0)
		for i := 1; i < len(res); i++ {
			res[i] = *big.NewInt(1)
		}
		resultChan <- res
		return
	}
	res := make([]big.Int, n+1)
	res[0] = *big.NewInt(0)
	for i := 1; i < n+1; i++ {
		res[i-1].Sub(&column[i], big.NewInt(0).Mul(big.NewInt(int64(k)), &column[i-1]))
	}
	endTime := time.Now().UnixMicro()
	TimePreviousColumn = append(TimePreviousColumn, endTime-startTime)
	resultChan <- res

}

func optimizedBlockDicho(n, k int, swap bool, rank big.Int, whichS3 int) ([]int, big.Int) {
	res := make([]int, 1)
	var acc *big.Int
	if !swap {
		acc = new(big.Int).Set(&StirlingColumn0[n-1])
	} else {
		acc = new(big.Int).Set(&StirlingColumn1[n-1])
	}

	if rank.Cmp(acc) < 0 {
		return res, *big.NewInt(0)
	}
	d0 := 1
	position := 2
	limitMin := 2
	limitMax := n
	completed := false
	for !completed {
		s3 := vs3[whichS3](n+1-position, k, swap, d0+1-position)

		tmp := new(big.Int).Sub(&rank, &s3)
		tmp.Sub(tmp, acc)

		var limitMiddle int
		for limitMin < limitMax {
			limitMiddle = (limitMin + limitMax) / 2
			tmpS3 := vs3[whichS3](n+1-position, k, swap, limitMiddle+1-position)
			tmpS3 = *tmpS3.Neg(&tmpS3)
			if tmp.Cmp(&tmpS3) >= 0 {
				limitMin = limitMiddle + 1
			} else {
				limitMax = limitMiddle
			}
		}
		limitMiddle = limitMin
		tmp2S3 := vs3[whichS3](n+1-position, k, swap, limitMiddle-position)
		middleRank := new(big.Int).Sub(&s3, &tmp2S3)
		middleRank.Add(middleRank, acc)
		res = append(res, limitMiddle-1-len(res))
		acc = middleRank
		var stirling big.Int
		if !swap {
			stirling = StirlingColumn0[n-position]
		} else {
			stirling = StirlingColumn1[n-position]
		}
		toCompare := new(big.Int).Add(&stirling, acc)
		if rank.Cmp(toCompare) < 0 {
			completed = true
		} else {
			position++
			d0 = limitMiddle
			limitMin = d0 + 1
			limitMax = n
			acc.Add(acc, &stirling)
		}
	}
	return res, *acc
}

// Return the k and the k-1th stirling triangle collumn until the line n.
func Stirling2Columns(n, k int) *types.CoupleColumns {
	// renvoie 2 colonnes de Stirling, k-1 et k jusqu'aux lignes n et n
	// on suppose k >= 1
	// il faut n-k+1 valeurs dans chaque colonne
	if k == 1 {
		c0 := make([]big.Int, n+1)
		c1 := make([]big.Int, n+1)
		c0[0] = *big.NewInt(1)
		for i := 1; i < n; i++ {
			c1[i] = *big.NewInt(1)
		}
		couple := types.CoupleColumns{Col0: c0, Col1: c1}
		return &couple
	}
	prev := make([]*big.Int, n+1)
	curr := make([]*big.Int, n+1)
	for i := range prev {
		prev[i] = big.NewInt(1)
	}
	prev[0] = big.NewInt(0)

	for j := 2; j < k+1; j++ {
		if j%2 == 0 {
			curr[j-2] = big.NewInt(0)
			curr[j-1] = big.NewInt(0)
			curr[j] = big.NewInt(1)

			for i := j + 1; i < n-k+1+j; i++ { // on peut améliorer borne max avec min...
				curr[i] = big.NewInt(0)
				curr[i].Mul(big.NewInt(int64(j)), curr[i-1])
				curr[i].Add(curr[i], prev[i-1])
			}
		} else {
			prev[j-2] = big.NewInt(0)
			prev[j-1] = big.NewInt(0)
			prev[j] = big.NewInt(1)

			for i := j + 1; i < n-k+1+j; i++ { // on peut améliorer borne max avec min...
				prev[i].Mul(big.NewInt(int64(j)), prev[i-1])
				prev[i].Add(prev[i], curr[i-1])
			}
		}

	}
	c0 := make([]big.Int, n+1)
	c1 := make([]big.Int, n+1)

	for i := 0; i < n; i++ {
		if k%2 == 0 {
			c0[i] = *big.NewInt(0).Set(prev[i])
			c1[i] = *big.NewInt(0).Set(curr[i])
		} else {
			c0[i] = *big.NewInt(0).Set(curr[i])
			c1[i] = *big.NewInt(0).Set(prev[i])
		}
	}
	if k%2 == 0 {
		c1[n] = *big.NewInt(0).Set(curr[n])
	} else {
		c1[n] = *big.NewInt(0).Set(prev[n])
	}

	couple := types.CoupleColumns{Col0: c0, Col1: c1}
	return &couple
}
