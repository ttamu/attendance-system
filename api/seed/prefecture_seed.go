package seed

import (
	"github.com/t2469/attendance-system.git/models"
	"log"

	"gorm.io/gorm"
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

type Prefecture = models.Prefecture

func SeedPrefectures(db *gorm.DB) error {
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
