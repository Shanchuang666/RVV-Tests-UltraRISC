package generator

import (
	"fmt"
	"math"
	"strings"
)

func (i *Insn) genCodeVs3Rs1mVm(pos int) []string {
	nfields := getNfieldsRoundedUp(i.Name)
	combinations := i.combinations(
		nfieldsLMULs(nfields),
		[]SEW{getEEW(i.Name)},
		[]bool{false, true},
		i.vxrms(),
	)
	res := make([]string, 0, len(combinations))

	for _, c := range combinations[pos:] {
		builder := strings.Builder{}
		builder.WriteString(c.initialize())

		builder.WriteString(i.gWriteRandomData(LMUL(1)))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(0, LMUL(1), SEW(32)))

		lmul1 := LMUL(math.Max(float64(c.LMUL1)*float64(nfields), 1))
		vs3, _, _ := getVRegs(lmul1, false, i.Name)
		builder.WriteString(i.gWriteIntegerTestData(lmul1, c.SEW, 0))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs3, lmul1, c.SEW))

		builder.WriteString(i.gResultDataAddr())

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s v%d, (a0)%s\n", i.Name, vs3, v0t(c.Mask)))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs3, lmul1, c.SEW))
		builder.WriteString(i.gMagicInsn(vs3, lmul1))

		res = append(res, builder.String())
	}
	return res
}
