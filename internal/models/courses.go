package models

type Courses struct {
	//ID ใช้ int ไม่ใช่หรอ หรือป่าวอ่ะ
	Course_ID         int     `json:"course_id"`
	Course_Name       string  `json:"course_name"`
	Course_Desc       string  `json:"course_desc"`
	Thumbnail_url     string  `json:"thumbnail_url"`
	Course_Type       string  `json:"course_type"`
	Course_Instructor string  `json:"course_instructor"`
	Profile_url       string  `json:"profile_url"`
	Course_Price      int     `json:"course_price"`
	Duration          int     `json:"duration"`
	Rating            float64 `json:"rating"`
	Num_reviews       int     `json:"num_reviews"`
	Enrollment_count  int     `json:"enrollment_count"`
	Created_at        string  `json:"created_at"`
	Updated_at        string  `json:"updated_at"`
}

// course_id SERIAL PRIMARY KEY,
// course_name VARCHAR(255) NOT NULL,
// description TEXT,
// thumbnail_url VARCHAR(255),
// instructor_name VARCHAR(255),
// profile_url VARCHAR(255),
// duration VARCHAR(255),
// price DECIMAL(10,2),
// detail_url VARCHAR(255),
// rating DECIMAL(2,1),
// num_reviews INT,
// enrollment_count INT,
// created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
// updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
