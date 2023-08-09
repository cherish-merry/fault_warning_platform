package outdoor

import (
	"github.com/RaymondCode/simple-demo/models"
)

func outdoorMachineMask(outdoorUnits []models.OutdoorDevice) map[string]float64 {
	outdoorNums := len(outdoorUnits)
	outdoorOnNum := 0.0
	avgDict := make(map[string]float64)

	for i := 0; i < outdoorNums; i++ {
		outdoorOnNum += outdoorUnits[i].Status
	}

	labels := []string{"H1", "Fo", "OE", "Pd", "Ps", "Td1", "TdSH", "Te1", "Ta", "Tfin", "A12", "A1", "Info1"}

	for _, label := range labels {
		sum := 0.0
		for i := 0; i < outdoorNums; i++ {
			sum += outdoorUnits[i].Status * outdoorUnits[i].H1
		}
		var avg float64
		if outdoorOnNum != 0 {
			avg = sum / outdoorOnNum
		} else {
			avg = 0
		}
		avgDict[label+"_arr"] = avg
	}
	return avgDict
}
