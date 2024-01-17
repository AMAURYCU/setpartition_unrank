package statistic

import (
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"github.com/AMAURYCU/setpartition_unrank/parallelunranking"
	"github.com/AMAURYCU/setpartition_unrank/precalcul"
)

func StirlingTriangle(n, k int) [2000][2000]*big.Int {
	if k < 0 || k > n {
		return [2000][2000]*big.Int{}
	}

	triangle := [2000][2000]*big.Int{}
	for i := 0; i <= n; i++ {
		triangle[i] = [2000]*big.Int{}
		for j := 0; j <= k; j++ {
			triangle[i][j] = new(big.Int)
		}
	}

	triangle[0][0].SetInt64(1)

	j := 1
	for i := j; i < n-k+1+j; i++ { // on peut améliorer borne max avec min...
		triangle[i][j].SetInt64(1)
	}
	for j := 2; j < k+1; j++ {
		for i := j; i < n-k+1+j; i++ { // on peut améliorer borne max avec min...  n-k+1+j
			temp1 := new(big.Int)
			if i-1 < j {
				temp1.SetInt64(0)
			} else {
				temp1.Mul(big.NewInt(int64(j)), triangle[i-1][j])
			}
			temp2 := new(big.Int)
			temp2.Set(triangle[i-1][j-1])
			triangle[i][j].Add(temp1, temp2)
		}
	}
	return triangle
}

func Stat(bsup int, nbPoints int, repetitions int, verbose bool) {

	file, err := os.Create("benchmark.py")
	fmt.Println(err)
	defer file.Close()
	buf := "import matplotlib.pyplot as plt \nfig = plt.figure(figsize=(12,6))\nax = fig.add_subplot(1,1,1)\n"
	dt := bsup / nbPoints
	ord := make([][]float64, 5)
	values := make([][]int64, 5)
	valuepre := make([]float64, 0)
	waiting := make([]float64, 0)
	col1 := make([]float64, 0)
	acol1 := int64(0)
	sumtimepre := int64(0)
	at := int64(0)
	abs := make([]int64, 0)
	sg := rand.NewSource(time.Now().UnixNano())
	rg := rand.New(sg)
	buf += "seed = " + strconv.Itoa(int(sg.Int63())) + "\n"
	for k := 2; k < bsup; k = k + dt {
		if verbose {
			out := exec.Command("clear")
			out.Stdout = os.Stdout
			out.Run()
			fmt.Println(k, " / ", bsup)
		}
		sumTime := make([]int64, 5)
		for useless := 0; useless < repetitions; useless++ {
			var r big.Int

			couple := *parallelunranking.Stirling2Columns(bsup, k)
			StirlingColumn1 := couple.Col1

			r.Rand(rg, &StirlingColumn1[bsup])
			for o := 0; o < 5; o++ {
				startTime := time.Now().UnixMicro()
				parallelunranking.UnrankDicho(bsup, k, r, o)
				endTime := time.Now().UnixMicro()
				sumTime[o] += endTime - startTime
				values[o] = append(values[o], endTime-startTime)
				acol1 += parallelunranking.WaitingTime
				parallelunranking.WaitingTime = 0

			}
			precalcul.StirlingMatrix = StirlingTriangle(bsup, k)
			startTime := time.Now().UnixMicro()
			precalcul.UnrankDichoPre(bsup, k, r, 0)
			endTime := time.Now().UnixMicro()
			sumtimepre += endTime - startTime
			at += parallelunranking.TimeTotal

		}
		waiting = append(waiting, float64(at)/float64(repetitions))
		col1 = append(col1, float64(acol1)/float64(repetitions))
		acol1 = 0
		at = 0
		abs = append(abs, int64(k))
		for s := 0; s < 5; s++ {
			ord[s] = append(ord[s], float64(sumTime[s])/float64(repetitions))
		}
		valuepre = append(valuepre, float64(sumtimepre)/float64(repetitions))
		sumtimepre = 0
	}
	
	buf += "t = " + ListToString(abs) + "\n"
	for u := 0; u < 5; u++ {
		buf += "l" + strconv.Itoa(u) + "=" + ListToStringFloat(ord[u]) + "\n"
	}
	buf += "lpre " + "= " + ListToStringFloat(valuepre) + "\n"
	buf += "latt " + "= " + ListToStringFloat(waiting) + "\n"
	buf += "wait_col1 " + "= " + ListToStringFloat(col1) + "\n"
	for u := 0; u < 5; u++ {
		buf += "ax.plot(t,l" + strconv.Itoa(u) + ",label = \"s3v" + strconv.Itoa(u) + "\")\n"
	}
	buf += "ax.plot(t,lpre,label = \"precalculs\")\n"
	buf += "ax.plot(t, wait_col1 , label = \"attente premières colones \")\n"
	buf += "latt = [k/1000 for k in latt]\n"
	buf += "sumPrecAtt = lpre.copy()\nbuf = []\n"
	buf += "for k in range(len(lpre)):\n"
	buf += "   buf+=[lpre[k]+latt[k]]\n"
	buf += "   wait_col1[k]+=lpre[k]\n"
	buf += "   wait_col1[k]+=latt[k]\n"
	buf += "ax.plot(t,wait_col1, label = \"precalcul+attente prev.Col + attente 1ere col\")\n"
	buf += "ax.plot(t,latt, label = \"temps d'attente \")\n"
	buf += "ax.plot(t,buf, label = \"precalculs+attente\")\n"
	buf += "ax.grid()\n"
	buf += "ax.legend()\n"

	buf += "ax.set_ylabel(\"temps d execution en μs\")\nax.set_xlabel(\"valeur de k\")\n"
	buf += "fig.savefig(\"courbe\")\nplt.show()"
	_, a := file.WriteString(buf)
	fmt.Println(a)
	fmt.Println(buf)
	file.Close()
	cmd := exec.Command("python3", "benchmark.py")
	out, erro := cmd.Output()
	fmt.Println(out, erro)

}

func ListToStringFloat(liste []float64) string {
	
	elements := make([]string, len(liste))
	for i, v := range liste {
		elements[i] = fmt.Sprintf("%f", v)
	}
	
	chaine := "[" + strings.Join(elements, ", ") + "]"
	return chaine
}

func ListToString(liste []int64) string {
	
	elements := make([]string, len(liste))
	for i, v := range liste {
		elements[i] = fmt.Sprintf("%d", v)
	}

	chaine := "[" + strings.Join(elements, ", ") + "]"

	return chaine
}

func PrintMatrix(mat [][]int64) string {
	buf := "["
	buf += (ListToString(mat[0]))
	for k := 1; k < len(mat); k++ {
		buf += ","
		buf += ListToString(mat[k])
	}
	buf += "]"
	return buf
}

func Graph3d(nmax int, dn, dk int, repetition int) ([][]int64, []int64, []int64) {
	sg := rand.NewSource(time.Now().UnixNano())
	rg := rand.New(sg)
	mat := make([][]int64, 0)
	valn := make([]int64, 0)
	valk := make([]int64, 0)
	for n := 200; n < nmax+dn; n += dn {
		fmt.Println(n, "/", nmax)
		res := make([]int64, 0)
		for k := 2; k < n; k += dk {
			sumtime := int64(0)
			for useless := 0; useless < repetition; useless++ {

				couple := *parallelunranking.Stirling2Columns(n, k)
				parallelunranking.StirlingColumn0 = couple.Col0
				parallelunranking.StirlingColumn1 = couple.Col1
				var r big.Int
				r.Rand(rg, &parallelunranking.StirlingColumn0[n])
				startTime := time.Now().UnixMilli()
				parallelunranking.UnrankDicho(n, k, r, 4)
				endTime := time.Now().UnixMilli()
				sumtime += endTime - startTime
			}
			res = append(res, sumtime/int64(repetition))

		}
		valn = append(valn, int64(n))
		mat = append(mat, res)

	}
	for k := 2; k < nmax+dn; k += dk {
		valk = append(valk, int64(k))
	}
	return mat, valn, valk

}
