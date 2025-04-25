package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"onlinecourse/database"
	"onlinecourse/internal/config"
	"onlinecourse/internal/handlers"
)

// Public Key ในรูปแบบ PEM (คัดลอกจาก Keycloak)
var publicKeyPEM = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzfiOGXueSlSBOU5l/vRG8jIl9kIxffHQZgWFVR1LmdEttq6mZWA+T3Ry8XmR//IHiUKtW8ypiLxrD8Xnpbd02NWfZoiXcTL94rHAZAgEi5Kz136iaD0pjZ54OzUcsjNzF6nV9Qq+VgvGhtHBs8VyCPzgyFGVXRiRje4layFhCtTUeuNFbEoxqN4Ua3xyd8k21726OIOKfHPY6LCCUiaSIjBxp6OdX5fFpnNss5EJABt7C1pF9/Hk3vKPa4ivqLisnQcT9+fJw22NLiCAfjMDtcfJXLWD+8mt3aNv+BYRkx7FFvwdSGR7NDL5e9wbmNitdMqmnsCJ20yKKDm8OFp8hwIDAQAB
-----END PUBLIC KEY-----`

// Global variable สำหรับเก็บ *rsa.PublicKey ที่แปลงแล้ว
var rsaPublicKey *rsa.PublicKey

// Casbin enforcer
var enforcer *casbin.Enforcer

func init() {
	var err error
	// สร้าง Casbin enforcer
	enforcer, err = casbin.NewEnforcer("model.conf", "policy.csv")
	if err != nil {
		panic("failed to create casbin enforcer: " + err.Error())
	}

	// Decode PEM และ parse public key เพียงครั้งเดียวตอนเริ่มต้น
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		panic("failed to parse PEM block containing the public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic("failed to parse DER encoded public key: " + err.Error())
	}
	var ok bool
	rsaPublicKey, ok = pub.(*rsa.PublicKey)
	if !ok {
		panic("key is not of type *rsa.PublicKey")
	}
}

func main() {
	// Connect to database
	// _ = godotenv.Load()
	cfg := config.LoadConfig()
	database.ConnectDB(cfg)

	r := gin.Default()
	r.POST("/register", handlers.Register)
	// เพิ่ม CORS middleware - ใช้ server ของ keycloak
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// ข้อ 1 request ได้เฉพาะคนที่สมัคร ----> ยังไม่ทดลอง
	protected := r.Group("/api")

	protected.Use(JWTAuthMiddleware())
	{
		// protected.GET("/dataLog", handlers.RequestLogMiddleware(), handlers.GetData)
		// protected.GET("/data", handlers.GetData)
		protected.GET("/Allcourses", handlers.GetAllCourses)

		protected.GET("/data", func(c *gin.Context) {
			username := c.GetString("username")
			roles := c.GetStringSlice("roles")
			c.JSON(200, gin.H{
				"data":     "This is protected data",
				"username": username,
				"roles":    roles,
			})
		})

		protected.GET("/token", func(c *gin.Context) {
			tokenRaw := c.GetHeader("Authorization")
			if tokenRaw == "" {
				c.JSON(400, gin.H{"error": "Authorization header not found"})
				return
			}

			// ตัดคำว่า "Bearer " ออก
			token := strings.TrimPrefix(tokenRaw, "Bearer ")

			c.JSON(200, gin.H{
				"token": token,
			})
		})
	}

	r.Run(":8081")
}

// Middleware สำหรับตรวจสอบ JWT และสิทธิ์ด้วย Casbin
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse JWT โดยใช้ rsaPublicKey ที่แปลงแล้วจากขั้นตอน initial setup
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return rsaPublicKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// ดึง claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// ตรวจสอบเวลาหมดอายุของโทเค็น
		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			c.JSON(401, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		// ตรวจสอบผู้ออกโทเค็น
		if !claims.VerifyIssuer("http://localhost:8082/realms/auth101", true) {
			c.JSON(401, gin.H{"error": "Invalid token issuer"})
			c.Abort()
			return
		}

		// ดึง username จาก claims
		username, ok := claims["preferred_username"].(string)
		if !ok {
			c.JSON(401, gin.H{"error": "Username not found in token"})
			c.Abort()
			return
		}

		// ดึง roles จาก realm_access.roles ซึ่งอาจมีหลาย role
		realmAccess, ok := claims["realm_access"].(map[string]interface{})
		if !ok {
			c.JSON(401, gin.H{"error": "Roles not found in token"})
			c.Abort()
			return
		}
		rawRoles, ok := realmAccess["roles"].([]interface{})
		if !ok || len(rawRoles) == 0 {
			c.JSON(401, gin.H{"error": "No roles found in token"})
			c.Abort()
			return
		}

		// ดึง role ทั้งหมดจาก payload
		var rolesList []string
		for _, r := range rawRoles {
			if roleStr, ok := r.(string); ok {
				rolesList = append(rolesList, roleStr)
			}
		}

		// ตรวจสอบสิทธิ์ด้วย Casbin: ให้ตรวจสอบว่ามี role ใดที่อนุญาตให้เข้าถึง resource ได้หรือไม่
		resource := c.Request.URL.Path // เช่น /api/data
		action := c.Request.Method     // เช่น GET
		allowed := false
		for _, role := range rolesList {
			permit, err := enforcer.Enforce(role, resource, action)
			if err != nil {
				c.JSON(500, gin.H{"error": "Error checking permission"})
				c.Abort()
				return
			}
			if permit {
				allowed = true
				break
			}
		}
		if !allowed {
			c.JSON(403, gin.H{"error": "Forbidden: Insufficient permissions"})
			c.Abort()
			return
		}

		// ส่ง username และ roles ไปยัง handler
		c.Set("username", username)
		c.Set("roles", rolesList)
		c.Next()
	}
}

// func JWTAuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		fmt.Println("🔍 AUTH HEADER:", authHeader)

// 		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
// 			fmt.Println("❌ ไม่พบ Bearer token ใน header")
// 			c.JSON(401, gin.H{"error": "Unauthorized"})
// 			c.Abort()
// 			return
// 		}
// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 		fmt.Println("🔐 JWT:", tokenString[:30]+"...") // แสดงแค่บางส่วนกันยาว

// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
// 				fmt.Println("❌ Signing method ไม่ถูกต้อง")
// 				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 			}
// 			return rsaPublicKey, nil
// 		})

// 		if err != nil || !token.Valid {
// 			fmt.Println("❌ Token ไม่ valid:", err)
// 			c.JSON(401, gin.H{"error": "Invalid token"})
// 			c.Abort()
// 			return
// 		}

// 		claims, ok := token.Claims.(jwt.MapClaims)
// 		if !ok {
// 			fmt.Println("❌ ดึง claims ไม่ได้")
// 			c.JSON(401, gin.H{"error": "Invalid token claims"})
// 			c.Abort()
// 			return
// 		}

// 		fmt.Println("✅ CLAIMS:", claims)

// 		// ตรวจ issuer
// 		if !claims.VerifyIssuer("http://localhost:8082/realms/auth101", true) {
// 			fmt.Println("❌ Issuer ไม่ตรง:", claims["iss"])
// 			c.JSON(401, gin.H{"error": "Invalid token issuer"})
// 			c.Abort()
// 			return
// 		}

// 		// ตรวจ username
// 		username, ok := claims["preferred_username"].(string)
// 		if !ok {
// 			fmt.Println("❌ ไม่พบ preferred_username ใน claims")
// 			c.JSON(401, gin.H{"error": "Username not found in token"})
// 			c.Abort()
// 			return
// 		}
// 		fmt.Println("👤 USER:", username)

// 		// ตรวจ roles
// 		realmAccess, ok := claims["realm_access"].(map[string]interface{})
// 		if !ok {
// 			fmt.Println("❌ realm_access ไม่ถูกต้อง")
// 			c.JSON(401, gin.H{"error": "Roles not found in token"})
// 			c.Abort()
// 			return
// 		}
// 		rawRoles, ok := realmAccess["roles"].([]interface{})
// 		if !ok {
// 			fmt.Println("❌ ไม่พบ roles array")
// 			c.JSON(401, gin.H{"error": "No roles found in token"})
// 			c.Abort()
// 			return
// 		}
// 		var rolesList []string
// 		for _, r := range rawRoles {
// 			if roleStr, ok := r.(string); ok {
// 				rolesList = append(rolesList, roleStr)
// 			}
// 		}
// 		fmt.Println("🛡️ ROLES:", rolesList)

// 		// ตรวจสิทธิ์ Casbin
// 		resource := c.Request.URL.Path
// 		action := c.Request.Method
// 		fmt.Printf("🔒 CHECK PERMISSION: role(s): %v → %s %s\n", rolesList, action, resource)

// 		allowed := false
// 		for _, role := range rolesList {
// 			permit, err := enforcer.Enforce(role, resource, action)
// 			if err != nil {
// 				fmt.Println("❌ ERROR จาก Casbin:", err)
// 				c.JSON(500, gin.H{"error": "Error checking permission"})
// 				c.Abort()
// 				return
// 			}
// 			if permit {
// 				allowed = true
// 				break
// 			}
// 		}
// 		if !allowed {
// 			fmt.Println("🚫 ไม่อนุญาตตาม policy")
// 			c.JSON(403, gin.H{"error": "Forbidden: Insufficient permissions"})
// 			c.Abort()
// 			return
// 		}

// 		c.Set("username", username)
// 		c.Set("roles", rolesList)
// 		c.Next()
// 	}
// }
