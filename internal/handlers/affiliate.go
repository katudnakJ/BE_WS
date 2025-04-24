package handlers

// pass รอแก้เป็น hash แล้วเก็บไว้

// Gen APIKey
// func generateAPIKey() string {
// 	bytes := make([]byte, 16)
// 	rand.Read(bytes)
// 	return hex.EncodeToString(bytes)
// }

// func Register(c *gin.Context) {
// 	var aff models.Affiliates
// 	if err := c.ShouldBindJSON(&aff); err != nil {

// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		log.Println("error", err.Error())
// 		return
// 	}

// 	Affiliate_APIKey := generateAPIKey()
// 	c.JSON(http.StatusOK, gin.H{"message": "You are registered, Thanks to join Us!", "api_key": Affiliate_APIKey})
// }
