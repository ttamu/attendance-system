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

// パーセンテージの文字列を整形
func sanitizePercentage(s string) string {
	s = strings.TrimSpace(s)
	if idx := strings.Index(s, "%"); idx != -1 {
		return s[:idx+1]
	}
	return s
}

// パーセンテージの文字列を数値に変換
func parsePercentage(s string) (float64, error) {
	s = sanitizePercentage(s)
	s = strings.TrimSuffix(s, "%")
	return strconv.ParseFloat(s, 64)
}

// 文字列からカンマを除去
func rmComma(s string) string {
	return strings.ReplaceAll(s, ",", "")
}

// 複数のExcelファイルをembed
// 各ファイルは1年度分の保険料データを含む（47都道府県分のデータ）
//
//go:embed insurance_rates_2021.xlsx insurance_rates_2022.xlsx insurance_rates_2023.xlsx insurance_rates_2024.xlsx insurance_rates_2025.xlsx
var insuranceRateFiles embed.FS

// InsuranceRateFiles Excelファイルと対応する西暦をマッピング
var InsuranceRateFiles = []struct {
	FileName string
	Year     int
}{
	{"insurance_rates_2021.xlsx", 2021}, // 令和3年度
	{"insurance_rates_2022.xlsx", 2022}, // 令和4年度
	{"insurance_rates_2023.xlsx", 2023}, // 令和5年度
	{"insurance_rates_2024.xlsx", 2024}, // 令和6年度
	{"insurance_rates_2025.xlsx", 2025}, // 令和7年度
}

// newMonth 新しい料率が適用される固定の月
const newMonth = 3

// SeedInsuranceRates は、埋め込みExcelファイルを読み込み、全都道府県の保険料データをシードデータとしてDBへ保存
func SeedInsuranceRates(db *gorm.DB) error {
	for _, fileMapping := range InsuranceRateFiles {
		fData, err := insuranceRateFiles.Open(fileMapping.FileName)
		if err != nil {
			return fmt.Errorf("failed to open embedded Excel file %s: %v", fileMapping.FileName, err)
		}

		f, err := excelize.OpenReader(fData)
		if err != nil {
			fData.Close()
			return fmt.Errorf("failed to open Excel file %s: %v", fileMapping.FileName, err)
		}
		defer f.Close()
		fData.Close()

		// Excelファイル内のシート一覧を取得
		sheetList := f.GetSheetList()
		if len(sheetList) == 0 {
			return fmt.Errorf("no sheet found in file %s", fileMapping.FileName)
		}

		for _, sheetName := range sheetList {
			var pref models.Prefecture
			if err := db.Where("name = ?", sheetName).First(&pref).Error; err != nil {
				return fmt.Errorf("prefecture '%s' not found", sheetName)
			}

			// 各シートのセルから、都道府県レベルの保険料率を取得
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

			// 都道府県ごとの保険料率を更新
			pref.HealthRateNoCare = healthNoCare
			pref.HealthRateWithCare = healthWithCare
			pref.PensionRate = pensionRate
			db.Save(&pref)

			rows, err := f.GetRows(sheetName)
			if err != nil {
				log.Printf("failed to get rows in sheet %s: %v", sheetName, err)
				continue
			}

			// Excelファイル内の税率がグレードごとに書かれた行の開始・終了位置
			startRow := 12
			endRow := 61

			// 期間：2025年度のデータなら2025年3月から2026年2月まで適用
			fromYear := fileMapping.Year
			fromMonth := newMonth
			toYear := fileMapping.Year + 1
			toMonth := newMonth - 1

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

				// グレードに健康保険と年金保険の情報が両方含まれているかチェック
				// Excelファイルには 4(1) のように書かれており、4は健康保険,1は年金保険のGrade
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

					// HealthInsuranceRate レコードを作成 or 更新
					var hRecord models.HealthInsuranceRate
					hResult := db.Where(&models.HealthInsuranceRate{
						PrefectureID: pref.ID,
						Grade:        healthGrade,
						FromYear:     fromYear,
						FromMonth:    fromMonth,
					}).Attrs(models.HealthInsuranceRate{
						MonthlyAmount:       monthlyAmount,
						MinMonthlyAmount:    minAmt,
						MaxMonthlyAmount:    maxAmt,
						HealthTotalNonCare:  hTotalNonCare,
						HealthHalfNonCare:   hHalfNonCare,
						HealthTotalWithCare: hTotalWithCare,
						HealthHalfWithCare:  hHalfWithCare,
						ToYear:              toYear,
						ToMonth:             toMonth,
					}).FirstOrCreate(&hRecord)
					if hResult.Error != nil {
						log.Printf("failed to create health record for sheet %s row %d: %v", sheetName, rowIdx, hResult.Error)
					}

					// 年金保険のグレードが存在する場合、PensionInsuranceRate レコードを作成 or 更新
					if pensionGrade != "" {
						var pRecord models.PensionInsuranceRate
						// grade が "32" の場合は最大値を設定する（それ以外はExcelからの値をそのまま利用）
						if pensionGrade == "32" {
							maxAmt = int(math.MaxInt64)
						}
						pResult := db.Where(&models.PensionInsuranceRate{
							PrefectureID: pref.ID,
							Grade:        pensionGrade,
							FromYear:     fromYear,
							FromMonth:    fromMonth,
						}).Attrs(models.PensionInsuranceRate{
							MonthlyAmount:    monthlyAmount,
							MinMonthlyAmount: minAmt,
							MaxMonthlyAmount: maxAmt,
							PensionTotal:     pTotal,
							PensionHalf:      pHalf,
							ToYear:           toYear,
							ToMonth:          toMonth,
						}).FirstOrCreate(&pRecord)
						if pResult.Error != nil {
							log.Printf("failed to create pension record for sheet %s row %d: %v", sheetName, rowIdx, pResult.Error)
						}
					}
				} else {
					// グレードに個別の年金情報が含まれていない場合
					var hRecord models.HealthInsuranceRate
					hResult := db.Where(&models.HealthInsuranceRate{
						PrefectureID: pref.ID,
						Grade:        grade,
						FromYear:     fromYear,
						FromMonth:    fromMonth,
					}).Attrs(models.HealthInsuranceRate{
						MonthlyAmount:       monthlyAmount,
						MinMonthlyAmount:    minAmt,
						MaxMonthlyAmount:    maxAmt,
						HealthTotalNonCare:  hTotalNonCare,
						HealthHalfNonCare:   hHalfNonCare,
						HealthTotalWithCare: hTotalWithCare,
						HealthHalfWithCare:  hHalfWithCare,
						ToYear:              toYear,
						ToMonth:             toMonth,
					}).FirstOrCreate(&hRecord)
					if hResult.Error != nil {
						log.Printf("failed to create health record for sheet %s row %d: %v", sheetName, rowIdx, hResult.Error)
					}
				}
			}
		}
		log.Printf("Finished processing file %s.", fileMapping.FileName)
	}
	return nil
}
