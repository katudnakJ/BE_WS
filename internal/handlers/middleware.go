// ยังไม่แก้
package handlers

// ต้องมี Middle -> Best practice
import (
	"net/http"
	"onlinecourse/database"
	"time"

	"github.com/gin-gonic/gin"
)

// ข้อ 5 request logs ----> ยังไม่ทดลอง
func RequestLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		search := c.Query("s")
		if search == "" {
			search = "All Course"
		}
		affID, exists := c.Get("affiliate_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		_, err := database.DB.Exec("INSERT INTO request_logs (affiliate_id, action,parameter,timestamp) VALUES ($1, $2, $3, $4)",
			affID, c.Request.URL.Path, search, time.Now())
		// c.Request.Method

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not log request"})
			c.Abort()
			return
		}

		c.Next()
	}
}
