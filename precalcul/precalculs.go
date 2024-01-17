package precalcul

import (
	"math/big"
	"github.com/AMAURYCU/setpartition_unrank/parallelunranking"
)

/*_____________________________PRE CALCULS____________________________________*/
var StirlingMatrix [2000][2000]*big.Int
var vs3pre []func(n, k, d int) *big.Int

func Init() {
	vs3pre = append(vs3pre, s3v5pre)

}

func UnrankDichoPre(n, k int, rank big.Int, vs3 int) [][]int {
	n0 := n
	res := make([][]int, 0)
	r := *new(big.Int).Set(&rank)

	if k == 1 {
		res := make([][]int, 0)
		tmp := make([]int, 0)
		for d := 1; d <= n; d++ {
			tmp = append(tmp, d)
		}
		res = append(res, tmp)
		return res
	}

	for k > 1 {
		block, acc := optimizedBlockDichoPre(n, k, r, vs3)
		res = append(res, block)
		r.Sub(&r, &acc)
		n -= len(block)
		k--
	}

	res = append(res, make([]int, n))
	res = parallelunranking.LexicographicPermutationUnrank(n0, res)

	return res

}

func optimizedBlockDichoPre(n, k int, rank big.Int, whichS3 int) ([]int, big.Int) {
	res := make([]int, 1)
	acc := new(big.Int).Set(StirlingMatrix[n-1][k-1])
	if rank.Cmp(acc) < 0 {
		return res, *big.NewInt(0)
	}
	d0 := 1
	position := 2
	limitMin := 2
	limitMax := n
	completed := false
	for !completed {
		s3 := vs3pre[whichS3](n+1-position, k, d0+1-position)
		tmp := new(big.Int).Sub(&rank, s3)
		tmp.Sub(tmp, acc)

		var limitMiddle int
		for limitMin < limitMax {
			limitMiddle = (limitMin + limitMax) / 2
			tmpS3 := vs3pre[whichS3](n+1-position, k, limitMiddle+1-position)
			tmpS3 = tmpS3.Neg(tmpS3)
			if tmp.Cmp(tmpS3) >= 0 {
				limitMin = limitMiddle + 1
			} else {
				limitMax = limitMiddle
			}

		}
		limitMiddle = limitMin
		tmp2S3 := vs3pre[whichS3](n+1-position, k, limitMiddle-position)
		middleRank := new(big.Int).Sub(s3, tmp2S3)
		middleRank.Add(middleRank, acc)
		res = append(res, limitMiddle-1-len(res))
		acc = middleRank
		stirling := StirlingMatrix[n-position][k-1]
		toCompare := new(big.Int).Add(stirling, acc)
		if rank.Cmp(toCompare) < 0 {
			completed = true
		} else {
			position++
			d0 = limitMiddle
			limitMin = d0 + 1
			limitMax = n
			acc.Add(acc, stirling)
		}

	}
	return res, *acc
}

func s3v4pre(n, k, d int) *big.Int {
	if d < 0 {
		return big.NewInt(0)
	}
	if d == 0 {
		if k-1 <= n && k-1 >= 0 {

			return StirlingMatrix[n+1][k]
		}
		return big.NewInt(0)
	}
	var res *big.Int
	if d%2 == 1 {
		res = new(big.Int).Sub(StirlingMatrix[n+1][k], StirlingMatrix[n+1-d][k])
	} else {
		res = new(big.Int).Add(StirlingMatrix[n+1][k], StirlingMatrix[n+1-d][k])
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
				tmp := new(big.Int).Sub(StirlingMatrix[n+1-u][k], StirlingMatrix[n+1-d+u][k])
				tmp.Mul(tmp, b)
				tmp.Mul(tmp, pm1)
				res.Add(res, tmp)
			} else {
				tmp := new(big.Int).Add(StirlingMatrix[n+1-u][k], StirlingMatrix[n+1-d+u][k])
				tmp.Mul(tmp, b)
				tmp.Mul(tmp, pm1)
				res.Add(res, tmp)
			}
		} else {
			tmp := new(big.Int).Mul(StirlingMatrix[n+1-u][k], b)
			tmp.Mul(tmp, pm1)
			res.Add(res, tmp)

		}
		pm1.Neg(pm1)
	}
	return res
}

func s3v2pre(n, k, d int) *big.Int {
	if d < 0 {
		return big.NewInt(0)
	}
	if d == 0 && !(k-1 <= n && k-1 >= 0) {
		return big.NewInt(0)
	}
	res := new(big.Int).Set(StirlingMatrix[n][k-1])
	if d >= k-1 {
		res.Add(res, StirlingMatrix[d][k-1])
	}
	u := 0
	b := big.NewInt(1)
	for u < min((n-d)/2, n-k+1) {
		b.Mul(b, big.NewInt(int64(n-d-u)))
		u++
		b.Div(b, big.NewInt(int64(u)))

		if (d+u >= k-1) && (u < (n-d)/2 || (u == (n-d)/2 && (n-d)%2 == 1)) {
			tmp := new(big.Int).Add(StirlingMatrix[n-u][k-1], StirlingMatrix[d+u][k-1])
			tmp.Mul(tmp, b)
			res.Add(res, tmp)
		} else {
			res.Add(res, new(big.Int).Mul(StirlingMatrix[n-u][k-1], b))
		}
	}
	return res
}

func s3v5pre(n, k, d int) *big.Int {
	if 2*d < n { //d < (n - d) {
		return s3v4pre(n, k, d)
	} else {
		return s3v2pre(n, k, d)
	}
}
