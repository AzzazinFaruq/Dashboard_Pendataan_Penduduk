package controllers

import (
	"backend_golang/models"
	"backend_golang/setup"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// computeDataCount menghitung jumlah penduduk & keluarga sesuai hak akses user.
func computeDataCount(user models.User) gin.H {
	var pendudukCount int64
	var keluargaCount int64

	if isAdmin(user) {
		setup.DB.Model(&models.Penduduk{}).Count(&pendudukCount)
		setup.DB.Model(&models.Keluarga{}).Count(&keluargaCount)
	} else {
		setup.DB.Model(&models.Penduduk{}).Where("user_id = ?", user.Id).Count(&pendudukCount)
		setup.DB.Model(&models.Keluarga{}).Where("user_id = ?", user.Id).Count(&keluargaCount)
	}

	return gin.H{"pendudukCount": pendudukCount, "keluargaCount": keluargaCount}
}

// computeAliveCount menghitung jumlah penduduk aktif & tidak aktif sesuai hak akses user.
func computeAliveCount(user models.User) gin.H {
	var alivepend int64
	var nopen int64

	if isAdmin(user) {
		setup.DB.Model(&models.Penduduk{}).Where("status = ?", 1).Count(&alivepend)
		setup.DB.Model(&models.Penduduk{}).Where("status = ?", 2).Count(&nopen)
	} else {
		setup.DB.Model(&models.Penduduk{}).Where("user_id = ?", user.Id).Where("status = ?", 1).Count(&alivepend)
		setup.DB.Model(&models.Penduduk{}).Where("user_id = ?", user.Id).Where("status = ?", 2).Count(&nopen)
	}

	return gin.H{"alivepend": alivepend, "nopen": nopen}
}

func DataCount(c *gin.Context) {
	user, ok := getCurrentUser(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, computeDataCount(user))
}

func AliveCount(c *gin.Context) {
	user, ok := getCurrentUser(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, computeAliveCount(user))
}

func MarryCount(c *gin.Context) {
	// Mengambil data user yang sedang login
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "status": false})
		return
	}

	var user models.User
	if err := setup.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "status": false})
		return
	}

	var belum int64
	var kawin int64
	var ceraihidup int64
	var ceraimati int64

	if user.Level == "admin" {
		setup.DB.Model(&models.Penduduk{}).Where("stat_kawin = ?", 1).Count(&kawin)
		setup.DB.Model(&models.Penduduk{}).Where("stat_kawin = ?", 2).Count(&ceraihidup)
		setup.DB.Model(&models.Penduduk{}).Where("stat_kawin = ?", 3).Count(&ceraimati)
		setup.DB.Model(&models.Penduduk{}).Where("stat_kawin = ?", 4).Count(&belum)
	} else {
		setup.DB.Model(&models.Penduduk{}).Where("user_id = ?", user.Id).Where("stat_kawin = ?", 1).Count(&kawin)
		setup.DB.Model(&models.Penduduk{}).Where("user_id = ?", user.Id).Where("stat_kawin = ?", 2).Count(&ceraihidup)
		setup.DB.Model(&models.Penduduk{}).Where("user_id = ?", user.Id).Where("stat_kawin = ?", 3).Count(&ceraimati)
		setup.DB.Model(&models.Penduduk{}).Where("user_id = ?", user.Id).Where("stat_kawin = ?", 4).Count(&belum)
	}

	c.JSON(http.StatusOK, gin.H{
		"kawin":      kawin,
		"ceraihidup": ceraihidup,
		"ceraimati":  ceraimati,
		"belum":      belum,
	})
}

func GenderCount(c *gin.Context) {
	// Mengambil data user yang sedang login
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "status": false})
		return
	}

	var user models.User
	if err := setup.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "status": false})
		return
	}

	var laki int64
	var perempuan int64

	if user.Level == "admin" {
		setup.DB.Model(&models.Penduduk{}).Where("kelamin = ?", 1).Count(&laki)
		setup.DB.Model(&models.Penduduk{}).Where("kelamin = ?", 2).Count(&perempuan)
	} else {
		setup.DB.Model(&models.Penduduk{}).Where("user_id = ?", user.Id).Where("kelamin = ?", 1).Count(&laki)
		setup.DB.Model(&models.Penduduk{}).Where("user_id = ?", user.Id).Where("kelamin = ?", 2).Count(&perempuan)
	}

	c.JSON(http.StatusOK, gin.H{
		"laki":      laki,
		"perempuan": perempuan,
	})
}

func RangeData(c *gin.Context) {
	user, ok := getCurrentUser(c)
	if !ok {
		return
	}

	// Mengambil parameter year dan month dari query
	yearInt, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter year tidak valid"})
		return
	}
	monthInt, err := strconv.Atoi(c.Query("month"))
	if err != nil || monthInt < 1 || monthInt > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter month tidak valid"})
		return
	}

	// Jumlah hari dalam bulan: tanggal 0 bulan berikutnya = hari terakhir bulan ini.
	// Cara ini otomatis benar untuk tahun kabisat (termasuk aturan abad).
	limit := time.Date(yearInt, time.Month(monthInt)+1, 0, 0, 0, 0, 0, time.UTC).Day()

	dataPerDay := make([]int64, 0, limit)
	for day := 1; day <= limit; day++ {
		var count int64
		query := setup.DB.Model(&models.Penduduk{}).
			Where("YEAR(created_at) = ?", yearInt).
			Where("MONTH(created_at) = ?", monthInt).
			Where("DAY(created_at) = ?", day)

		// User biasa hanya melihat datanya sendiri; admin melihat semua.
		if !isAdmin(user) {
			query = query.Where("user_id = ?", user.Id)
		}

		query.Count(&count)
		dataPerDay = append(dataPerDay, count)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dataPerDay,
	})
}

func AllData(c *gin.Context) {
	user, ok := getCurrentUser(c)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jumlah": computeDataCount(user),
		"status": computeAliveCount(user),
	})
}
