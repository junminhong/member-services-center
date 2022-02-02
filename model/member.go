package model

import (
	"gorm.io/gorm"
	"time"
)

// Member Table model
// ID = member table自動產生的index流水號
// uuid = 具備不重複索引值的功能，為了避免使用ID容易造成安全性問題，額外開出來的欄位
// Email = 用戶的信箱，也是個不重複的值
// AtomicToken = 該用戶目前使用的token
// RefreshAtomicToken = 用來重新刷新token的token
// EmailAuth = 用來判斷用戶是否經過信箱驗證
// ThirdAuthToken = 未來支援第三方登入用的
// SafePassword = 提供帳戶有安全密碼可以設定
// AtomicPoint = 點數，後續可以用來支付服務費
// ActivatedAt = 活躍時間
// CreatedAt = 創建該資料時間
// UpdatedAt = 資料更新時間
type Member struct {
	gorm.Model
	ID                 int    `gorm:"primaryKey"`
	uuid               string `gorm:"index;unique"`
	Email              string `gorm:"unique"`
	Password           string
	AtomicToken        string
	RefreshAtomicToken string
	EmailAuth          bool
	ThirdAuthToken     string
	SafePassword       string
	AtomicPoint        int
	MemberInfo         MemberInfo `gorm:"foreignKey:MemberInfoID"`
	ActivatedAt        time.Time
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}

// MemberInfo Table model
// ID = 資料流水號
// NickName = 用戶暱稱
// MugShotPath = 大頭貼路徑位置
// SocialInfo = 社群資料，未來支援freely產品要用的
type MemberInfo struct {
	gorm.Model
	MemberInfoID int `gorm:"primaryKey"`
	NickName     string
	MugShotPath  string
	SocialInfo   SocialInfo `gorm:"foreignKey:SocialInfoID"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime"`
}

// SocialInfo Table model
// ID = 流水號
// SocialType = 社群類型，facebook or google?
// SocialUrl = 社群連結
type SocialInfo struct {
	gorm.Model
	SocialInfoID int `gorm:"primaryKey"`
	SocialType   string
	SocialUrl    string
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
