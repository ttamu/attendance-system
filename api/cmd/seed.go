package main

import (
	"fmt"
	"github.com/t2469/labor-management-system.git/config"
	"github.com/t2469/labor-management-system.git/db"
	"github.com/t2469/labor-management-system.git/models"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
	"time"
)

var prefectureNames = []string{
	"北海道", "青森", "岩手", "宮城", "秋田", "山形", "福島",
	"茨城", "栃木", "群馬", "埼玉", "千葉", "東京", "神奈川",
	"新潟", "富山", "石川", "福井", "山梨", "長野", "岐阜",
	"静岡", "愛知", "三重", "滋賀", "京都", "大阪", "兵庫",
	"奈良", "和歌山", "鳥取", "島根", "岡山", "広島", "山口",
	"徳島", "香川", "愛媛", "高知", "福岡", "佐賀", "長崎",
	"熊本", "大分", "宮崎", "鹿児島", "沖縄",
}

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

func seedPrefectures(db *gorm.DB) error {
	for _, name := range prefectureNames {
		var pref = models.Prefecture{
			Name:               name,
			HealthRateNoCare:   0,
			HealthRateWithCare: 0,
			PensionRate:        0,
		}
		if err := db.Create(&pref).Error; err != nil {
			log.Fatal(err)
		} else {
			log.Println("Created prefecture", pref.Name)
		}
	}
	return nil
}

func seedInsuranceRates(db *gorm.DB, filePath string) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %v", err)
	}
	defer f.Close()

	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		return fmt.Errorf("no sheet found")
	}

	for _, sheetName := range sheetList {
		log.Printf("Processing sheet %s\n", sheetName)
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
		endRow := len(rows)

		for rowIdx := startRow; rowIdx <= endRow; rowIdx++ {
			rowData := rows[rowIdx-1]
			if len(rowData) < 10 {
				continue
			}

			grade := strings.ToValidUTF8(strings.TrimSpace(rowData[0]), "")
			if grade == "" {
				continue
			}

			minStr := rmComma(strings.TrimSpace(rowData[2]))
			if minStr == "" {
				minStr = "0"
			}
			maxStr := rmComma(strings.TrimSpace(rowData[4]))
			minAmt, err1 := strconv.Atoi(minStr)
			maxAmt, err2 := strconv.Atoi(maxStr)
			if err1 != nil || err2 != nil {
				log.Printf("failed to parse monthly amounts in sheet %s row %d", sheetName, rowIdx)
				continue
			}

			var hTotalNonCare, hHalfNonCare, hTotalWithCare, hHalfWithCare float64
			hTotalNonCare, _ = strconv.ParseFloat(strings.TrimSpace(rowData[5]), 64)
			hHalfNonCare, _ = strconv.ParseFloat(strings.TrimSpace(rowData[6]), 64)
			hTotalWithCare, _ = strconv.ParseFloat(strings.TrimSpace(rowData[7]), 64)
			hHalfWithCare, _ = strconv.ParseFloat(strings.TrimSpace(rowData[8]), 64)

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
				var pTotal, pHalf float64
				pTotal, _ = strconv.ParseFloat(strings.TrimSpace(rowData[9]), 64)
				pHalf, _ = strconv.ParseFloat(strings.TrimSpace(rowData[10]), 64)

				var healthGrade, pensionGrade string
				if leftIdx != -1 && rightIdx != -1 && rightIdx > leftIdx {
					healthGrade = strings.TrimSpace(string(runes[:leftIdx]))
					pensionGrade = strings.TrimSpace(string(runes[leftIdx+1 : rightIdx]))
				} else {
					healthGrade = grade
					pensionGrade = ""
				}

				hRecord := models.HealthInsuranceRate{
					PrefectureID:        pref.ID,
					Grade:               healthGrade,
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

func main() {
	_ = config.LoadEnv()
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	time.Local = loc

	db.InitDB()
	if err := db.DB.AutoMigrate(
		&models.Prefecture{},
		&models.HealthInsuranceRate{},
		&models.PensionInsuranceRate{},
		&models.User{},
		&models.Attendance{},
	); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	if err := seedPrefectures(db.DB); err != nil {
		log.Fatalf("seedPrefectures failed: %v", err)
	}

	if err := seedInsuranceRates(db.DB, "seed/insurance_rates.xlsx"); err != nil {
		log.Fatalf("seedInsuranceRates failed: %v", err)
	}

	log.Println("Seeding completed successfully.")
}
