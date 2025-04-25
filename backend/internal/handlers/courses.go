package handlers

import (
	"fmt"
	"log"
	"net/http"
	"onlinecourse/database"
	"onlinecourse/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

// ข้อ 3
func GetData(c *gin.Context) {

	//รับ Query search
	// ค้นหาได้ 3 ค่า คือ ชื่อ ประเภท คนสอน
	req := "SELECT id, content FROM data"
	if search := c.Query("s"); search != "" {
		req += fmt.Sprintf(" Where Course_Name ILIKE %s", search)
		req += fmt.Sprintf(" or Course_Type ILIKE %s", search)
		req += fmt.Sprintf(" or Course_Instructor ILIKE %s", search)
	}
	rows, err := database.DB.Query(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch data"})
		log.Println("error desc", err)
		return
	}
	defer rows.Close()

	// Initialize เป็น slice ว่างแทนที่จะเป็น nil
	results := make([]models.Courses, 0)
	for rows.Next() {
		var data models.Courses
		if err := rows.Scan(&data.Course_ID, &data.Course_Name, &data.Course_Type, &data.Course_Instructor, &data.Course_Price); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Data scan error"})
			return
		}
		results = append(results, data)
	}

	c.JSON(http.StatusOK, results)

	// func(c *gin.Context) {
	// 	username := c.GetString("username")
	// 	roles := c.GetStringSlice("roles")
	// 	c.JSON(200, gin.H{
	// 		"data":     "This is protected data",
	// 		"username": username,
	// 		"roles":    roles,
	// 	})
	// })
}

func GetAllCourses(c *gin.Context) {

	request := models.RequestLog{
		Affiliate_ID: c.GetString("affiliate_id"),
		Action:       "GET api/Allcourses",
		Parameter:    "No parameter",
		Timestamp:    time.Now().Format(time.RFC3339),
	}

	query := `INSERT INTO request_logs (affiliate_id, action, parameter, timestamp) VALUES ($1, $2, $3, $4)`

	_, err := database.DB.Exec(query, request.Affiliate_ID, request.Action, request.Parameter, request.Timestamp)
	if err != nil {
		log.Println("Error inserting request log:", err)
	}

	rows, err := database.DB.Query("select * from courses")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูล courses ได้"})
		return
	}
	defer rows.Close()

	var courses []models.Courses

	for rows.Next() {
		var course models.Courses
		err := rows.Scan(&course.Course_ID, &course.Course_Name, &course.Course_Desc, &course.Thumbnail_Url,
			&course.Course_Type, &course.Course_Instructor, &course.Profile_Url, &course.Course_Price,
			&course.Duration, &course.Rating, &course.Num_reviews, &course.Enrollment_count,
			&course.Created_at, &course.Updated_at)
		if err != nil {
			log.Println("Scan Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "เกิดข้อผิดพลาดในการอ่านข้อมูล"})
			return
		}
		courses = append(courses, course)
	}

	// ตรวจสอบ error หลังวน loop เสร็จ
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "เกิดข้อผิดพลาดในการอ่านข้อมูล (rows.Err)"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"courses": courses,
	})
}
