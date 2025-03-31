package seed

import (
	"embed"
	"fmt"
	"github.com/t2469/attendance-system.git/models"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"log"
	"math"
	"strconv"
	"strings"
)

func sanitizePercentage(s string) string {
	s = strings.TrimSpace(s)
	if idx := strings.Index(s, "%"); idx != -1 {
		return s[:idx+1]
	}
	return s
}

func parsePercentage(s string) (float64, error) {
	s = sanitizePercentage(s)
	s = strings.TrimSuffix(s, "%")
	return strconv.ParseFloat(s, 64)
}

func rmComma(s string) string {
	return strings.ReplaceAll(s, ",", "")
}

//go:embed insurance_rates.xlsx
var insuranceFile embed.FS

func SeedInsuranceRates(db *gorm.DB) error {
	fData, err := insuranceFile.Open("insurance_rates.xlsx")
	if err != nil {
		return fmt.Errorf("failed to open embedded Excel file: %v", err)
	}

	defer fData.Close()

	f, err := excelize.OpenReader(fData)
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %v", err)
	}
	defer f.Close()

	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		return fmt.Errorf("no sheet found")
	}

	for _, sheetName := range sheetList {
		var pref models.Prefecture
		if err := db.Where("name = ?", sheetName).First(&pref).Error; err != nil {
			return fmt.Errorf("prefecture '%s' not found", sheetName)
		}

		healthNoCareStr, _ := f.GetCellValue(sheetName, "F9")
		healthNoCareStr = strings.ToValidUTF8(healthNoCareStr, "")
		healthWithCareStr, _ := f.GetCellValue(sheetName, "H9")
		healthWithCareStr = strings.ToValidUTF8(healthWithCareStr, "")
		pensionStr, _ := f.GetCellValue(sheetName, "J9")
		pensionStr = strings.ToValidUTF8(pensionStr, "")

		healthNoCare, err1 := parsePercentage(healthNoCareStr)
		if err1 != nil {
			log.Printf("failed to parse HealthRateNoCare in sheet %s: %v", sheetName, err1)
		}
		healthWithCare, err2 := parsePercentage(healthWithCareStr)
		if err2 != nil {
			log.Printf("failed to parse HealthRateWithCare in sheet %s: %v", sheetName, err2)
		}
		pensionRate, err3 := parsePercentage(pensionStr)
		if err3 != nil {
			log.Printf("failed to parse PensionRate in sheet %s: %v", sheetName, err3)
		}

		pref.HealthRateNoCare = healthNoCare
		pref.HealthRateWithCare = healthWithCare
		pref.PensionRate = pensionRate
		db.Save(&pref)

		rows, err := f.GetRows(sheetName)
		if err != nil {
			log.Printf("failed to get rows in sheet %s: %v", sheetName, err)
			continue
		}

		startRow := 12
		endRow := 61

		for rowIdx := startRow; rowIdx <= endRow; rowIdx++ {
			rowData := rows[rowIdx-1]
			grade := strings.ToValidUTF8(strings.TrimSpace(rowData[0]), "")
			monthlyAmountStr := rmComma(strings.TrimSpace(rowData[1]))
			monthlyAmount, errMA := strconv.Atoi(monthlyAmountStr)
			if errMA != nil {
				log.Printf("error: invalid monthly amount at row %d in sheet %s: %v", rowIdx, sheetName, errMA)
				continue
			}

			minStr := rmComma(strings.TrimSpace(rowData[2]))
			maxStr := rmComma(strings.TrimSpace(rowData[4]))
			if minStr == "" {
				minStr = "0"
			}
			if maxStr == "" {
				maxStr = strconv.FormatInt(int64(math.MaxInt64), 10)
			}
			minAmt, err1 := strconv.Atoi(minStr)
			maxAmt, err2 := strconv.Atoi(maxStr)
			if err1 != nil || err2 != nil {
				log.Printf("failed to parse monthly amounts in sheet %s row %d", sheetName, rowIdx)
				continue
			}

			var hTotalNonCare, hHalfNonCare, hTotalWithCare, hHalfWithCare float64
			hTotalNonCare, _ = strconv.ParseFloat(rmComma(strings.TrimSpace(rowData[5])), 64)
			hHalfNonCare, _ = strconv.ParseFloat(rmComma(strings.TrimSpace(rowData[6])), 64)
			hTotalWithCare, _ = strconv.ParseFloat(rmComma(strings.TrimSpace(rowData[7])), 64)
			hHalfWithCare, _ = strconv.ParseFloat(rmComma(strings.TrimSpace(rowData[8])), 64)

			if strings.Contains(grade, "(") || strings.Contains(grade, "（") {
				runes := []rune(grade)
				var leftIdx, rightIdx int = -1, -1
				for i, ch := range runes {
					if ch == '（' || ch == '(' {
						leftIdx = i
						break
					}
				}
				for i, ch := range runes {
					if ch == '）' || ch == ')' {
						rightIdx = i
						break
					}
				}

				var healthGrade, pensionGrade string
				if leftIdx != -1 && rightIdx != -1 && rightIdx > leftIdx {
					healthGrade = strings.TrimSpace(string(runes[:leftIdx]))
					pensionGrade = strings.TrimSpace(string(runes[leftIdx+1 : rightIdx]))
				} else {
					healthGrade = grade
					pensionGrade = ""
				}

				var pTotal, pHalf float64
				pTotal, _ = strconv.ParseFloat(rmComma(strings.TrimSpace(rowData[9])), 64)
				pHalf, _ = strconv.ParseFloat(rmComma(strings.TrimSpace(rowData[10])), 64)

				hRecord := models.HealthInsuranceRate{
					PrefectureID:        pref.ID,
					Grade:               healthGrade,
					MonthlyAmount:       monthlyAmount,
					MinMonthlyAmount:    minAmt,
					MaxMonthlyAmount:    maxAmt,
					HealthTotalNonCare:  hTotalNonCare,
					HealthHalfNonCare:   hHalfNonCare,
					HealthTotalWithCare: hTotalWithCare,
					HealthHalfWithCare:  hHalfWithCare,
				}
				if err := db.Create(&hRecord).Error; err != nil {
					log.Printf("failed to create health record for sheet %s row %d: %v", sheetName, rowIdx, err)
				}

				if pensionGrade != "" {
					pRecord := models.PensionInsuranceRate{
						PrefectureID:     pref.ID,
						Grade:            pensionGrade,
						MonthlyAmount:    monthlyAmount,
						MinMonthlyAmount: minAmt,
						MaxMonthlyAmount: maxAmt,
						PensionTotal:     pTotal,
						PensionHalf:      pHalf,
					}
					if err := db.Create(&pRecord).Error; err != nil {
						log.Printf("failed to create pension record for sheet %s row %d: %v", sheetName, rowIdx, err)
					}
				}
			} else {
				hRecord := models.HealthInsuranceRate{
					PrefectureID:        pref.ID,
					Grade:               grade,
					MonthlyAmount:       monthlyAmount,
					MinMonthlyAmount:    minAmt,
					MaxMonthlyAmount:    maxAmt,
					HealthTotalNonCare:  hTotalNonCare,
					HealthHalfNonCare:   hHalfNonCare,
					HealthTotalWithCare: hTotalWithCare,
					HealthHalfWithCare:  hHalfWithCare,
				}
				if err := db.Create(&hRecord).Error; err != nil {
					log.Printf("failed to create health record for sheet %s row %d: %v", sheetName, rowIdx, err)
				}
			}
		}
		log.Printf("Sheet %s processed (rows %d - %d).", sheetName, startRow, endRow)
	}
	return nil
}
