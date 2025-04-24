package handlers

import (
	"fmt"
	"log"
	"net/http"
	"onlinecourse/database"
	"onlinecourse/internal/models"

	"github.com/gin-gonic/gin"
)

// ข้อ 3
func GetData(c *gin.Context) {
	//รับ Query search
	// ค้นหาได้ 3 ค่า คือ ชื่อคอส วิชา คนสอน
	var search string
	req := "SELECT id, content FROM data WHERE 1=1"

	if search = c.Query("n"); search != "" {
		req += fmt.Sprintf(" and Course_Name ILIKE %s", search)
	}
	if search = c.Query("t"); search != "" {
		req += fmt.Sprintf(" and Course_Type ILIKE %s", search)
	}
	if search = c.Query("i"); search != "" {
		req += fmt.Sprintf(" and Course_Instructor ILIKE %s", search)
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
