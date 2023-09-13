package outdoor

import (
	"github.com/RaymondCode/simple-demo/models"
)

func MachineMask(outdoorUnits []*models.OutdoorDevice) map[string]float64 {
	outdoorNums := len(outdoorUnits)
	outdoorOnNum := 0.0
	avgDict := make(map[string]float64)

	for i := 0; i < outdoorNums; i++ {
		outdoorOnNum += outdoorUnits[i].Status
	}

	sumH1 := 0.0
	sumFo := 0.0
	sumOE := 0.0
	sumPd := 0.0
	sumPs := 0.0
	sumTd1 := 0.0
	sumTdSH := 0.0
	sumTe1 := 0.0
	sumTa := 0.0
	sumTfin := 0.0
	sumA12 := 0.0
	sumA1 := 0.0
	sumInfo1 := 0.0

	for i := 0; i < outdoorNums; i++ {
		sumH1 += outdoorUnits[i].Status * outdoorUnits[i].H1
		sumFo += outdoorUnits[i].Status * outdoorUnits[i].Fo
		sumOE += outdoorUnits[i].Status * outdoorUnits[i].OE
		sumPd += outdoorUnits[i].Status * outdoorUnits[i].Pd
		sumPs += outdoorUnits[i].Status * outdoorUnits[i].Ps
		sumTd1 += outdoorUnits[i].Status * outdoorUnits[i].Td1
		sumTdSH += outdoorUnits[i].Status * outdoorUnits[i].TdSH
		sumTe1 += outdoorUnits[i].Status * outdoorUnits[i].Te1
		sumTa += outdoorUnits[i].Status * outdoorUnits[i].Ta
		sumTfin += outdoorUnits[i].Status * outdoorUnits[i].Tfin
		sumA12 += outdoorUnits[i].Status * outdoorUnits[i].A12
		sumA1 += outdoorUnits[i].Status * outdoorUnits[i].A1
		sumInfo1 += outdoorUnits[i].Status * outdoorUnits[i].Info1
	}

	avgH1 := 0.0
	avgFo := 0.0
	avgOE := 0.0
	avgPd := 0.0
	avgPs := 0.0
	avgTd1 := 0.0
	avgTdSH := 0.0
	avgTe1 := 0.0
	avgTa := 0.0
	avgTfin := 0.0
	avgA12 := 0.0
	avgA1 := 0.0
	avgInfo1 := 0.0

	if outdoorOnNum != 0 {
		avgOE = sumOE / outdoorOnNum
		avgPd = sumPd / outdoorOnNum
		avgPs = sumPs / outdoorOnNum
		avgTd1 = sumTd1 / outdoorOnNum
		avgTdSH = sumTdSH / outdoorOnNum
		avgTe1 = sumTe1 / outdoorOnNum
		avgTa = sumTa / outdoorOnNum
		avgTfin = sumTfin / outdoorOnNum
		avgA12 = sumA12 / outdoorOnNum
		avgA1 = sumA1 / outdoorOnNum
		avgInfo1 = sumInfo1 / outdoorOnNum
		avgH1 = sumH1 / outdoorOnNum
		avgFo = sumFo / outdoorOnNum
	}

	avgDict["H1"] = avgH1
	avgDict["Fo"] = avgFo
	avgDict["OE"] = avgOE
	avgDict["Pd"] = avgPd
	avgDict["Ps"] = avgPs
	avgDict["Td1"] = avgTd1
	avgDict["TdSH"] = avgTdSH
	avgDict["Te1"] = avgTe1
	avgDict["Ta"] = avgTa
	avgDict["Tfin"] = avgTfin
	avgDict["A12"] = avgA12
	avgDict["A1"] = avgA1
	avgDict["Info1"] = avgInfo1

	return avgDict
}
