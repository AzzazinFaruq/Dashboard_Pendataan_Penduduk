package controllers

import (
	"backend_golang/models"
	"backend_golang/setup"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// getCurrentUser mengambil user yang sedang login dari klaim JWT (di-set oleh AuthMiddleware).
// Jika gagal, helper langsung menulis response error dan mengembalikan ok=false,
// sehingga caller cukup melakukan `return`.
func getCurrentUser(c *gin.Context) (models.User, bool) {
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found", "status": false})
		return models.User{}, false
	}

	var user models.User
	if err := setup.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found", "status": false})
		return models.User{}, false
	}

	return user, true
}

// isAdmin mengecek apakah user memiliki hak akses admin.
func isAdmin(user models.User) bool {
	return user.Level == "admin" || user.Level == "superAdmin"
}

// paginationParams membaca ?page= & ?limit= dari query string.
// Mengembalikan ok=false bila keduanya tidak ada, sehingga caller dapat
// mengembalikan seluruh data (backward compatible).
func paginationParams(c *gin.Context) (page int, limit int, ok bool) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	if pageStr == "" && limitStr == "" {
		return 0, 0, false
	}

	page, _ = strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit, _ = strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 10
	}
	return page, limit, true
}

// canAccessResource memastikan user hanya boleh mengakses resource miliknya sendiri,
// kecuali ia seorang admin. Jika tidak diizinkan, helper menulis response 403
// dan mengembalikan false.
func canAccessResource(c *gin.Context, user models.User, ownerID int64) bool {
	if isAdmin(user) || int64(user.Id) == ownerID {
		return true
	}
	c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses ke data ini", "status": false})
	return false
}
