package controllers

import (
	"backend_golang/models"
	"backend_golang/setup"
	"strconv"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DataCount(c *gin.Context) {
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

	var pendudukCount int64
	var keluargaCount int64

	if user.Level == "admin" {
		setup.DB.Model(&models.Penduduk{}).Count(&pendudukCount)
		setup.DB.Model(&models.Keluarga{}).Count(&keluargaCount)
	} else {
		setup.DB.Model(&models.Penduduk{}).Where("user_id = ?", user.Id).Count(&pendudukCount)
		setup.DB.Model(&models.Keluarga{}).Where("user_id = ?", user.Id).Count(&keluargaCount)
	}

	c.Keys["response"] = gin.H{"pendudukCount": pendudukCount, "keluargaCount": keluargaCount}
}

func AliveCount(c *gin.Context) {
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

	var alivepend int64
	var nopen int64

	if user.Level == "admin" {
		setup.DB.Model(&models.Penduduk{}).Where("status = ?", 1).Count(&alivepend)
		setup.DB.Model(&models.Penduduk{}).Where("status = ?", 2).Count(&nopen)
	} else {
		setup.DB.Model(&models.Penduduk{}).Where("user_id = ?", user.Id).Where("status = ?", 1).Count(&alivepend)
		setup.DB.Model(&models.Penduduk{}).Where("user_id = ?", user.Id).Where("status = ?", 2).Count(&nopen)
	}

	c.Keys["response"] = gin.H{"alivepend": alivepend, "nopen": nopen}
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
		"kawin": kawin,
		"ceraihidup": ceraihidup,
		"ceraimati": ceraimati,
		"belum": belum,
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
		"laki": laki,
		"perempuan": perempuan,
	})
}

func RangeData(c *gin.Context) {
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

	// Mengambil parameter year dan month dari query
	year := c.Query("year")
	month := c.Query("month")
	// Konversi string ke int
	yearInt, _ := strconv.Atoi(year)
	monthInt, _ := strconv.Atoi(month)

	// Menentukan jumlah hari dalam bulan 
	var limit int
	if monthInt == 2 && yearInt%4 == 0 {
		limit = 29
	} else if monthInt == 2 && yearInt%4 == 1 {
		limit = 28
	} else if monthInt%2 == 1 {
		limit = 31
	} else if monthInt%2 == 0 {
		limit = 30
	}

	dataPerDay := make([]int64, 0)

	// Query berdasarkan role user
	if user.Level == "user" {
		for day := 1; day <= limit; day++ {
			var count int64
			setup.DB.Model(&models.Penduduk{}).
				Where("user_id = ?", user.Id).
				Where("YEAR(created_at) = ?", year).
				Where("MONTH(created_at) = ?", month).
				Where("DAY(created_at) = ?", day).
				Count(&count)
			
			dataPerDay = append(dataPerDay, count)
		}
	} else if user.Level == "admin" || user.Level == "superAdmin" {
		for day := 1; day <= limit; day++ {
			var count int64
			setup.DB.Model(&models.Penduduk{}).
				Where("YEAR(created_at) = ?", year).
				Where("MONTH(created_at) = ?", month).
				Where("DAY(created_at) = ?", day).
				Count(&count)
			
			dataPerDay = append(dataPerDay, count)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dataPerDay,
	})
}

func AllData(c *gin.Context) {
	// Mengambil data jumlah dari fungsi DataCount
	var jumlahResponse gin.H
	DataCount(c)
	if c.Writer.Status() != http.StatusOK {
		return // Jika terjadi error di DataCount, fungsi akan berhenti
	}
	c.Writer.Header().Del("Content-Type") // Reset header untuk mencegah multiple response
	jumlahResponse = c.Keys["response"].(gin.H)

	// Mengambil data status dari fungsi aliveCount
	var statusResponse gin.H
	AliveCount(c)
	if c.Writer.Status() != http.StatusOK {
		return // Jika terjadi error di aliveCount, fungsi akan berhenti
	}
	c.Writer.Header().Del("Content-Type") // Reset header untuk mencegah multiple response
	statusResponse = c.Keys["response"].(gin.H)

	// Menggabungkan semua data dan mengirim response
	c.JSON(http.StatusOK, gin.H{
		"jumlah": jumlahResponse,
		"status": statusResponse,
	})
}

