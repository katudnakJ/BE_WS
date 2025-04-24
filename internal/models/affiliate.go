package models

// ใส่หลัง json ให้เหมือนใน DB
type Affiliates struct {
	Affiliate_ID       int    `json:"affiliate_id"`
	Affiliate_Name     string `json:"affiliate_name"`
	Affiliate_Email    string `json:"affiliate_email"`
	Affiliate_Password string `json:"affiliate_password"`
}

type Affiliate_Url struct {
	Url_id        int    `json:"url_id"`
	Affiliate_ID  int    `json:"affiliate_id"`
	Affiliate_Url string `json:"aff_url"`
	Clicks        int    `json:"clicks"`
	Parameter     string `json:"parameter"`
}

// type Response_AffRegister struct {
// 	Status  string
// 	Message string
// }

// 	"clicks": 150,
// 	"sales": 45,
// 	"commission_earned": 75.5,
// 	"conversion_rate": 0.3
