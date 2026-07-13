package controllers

import (
	"backend_golang/config"
	"backend_golang/models"
	"backend_golang/setup"
	"backend_golang/utils"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
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

	var keluarga []models.Keluarga
	var total int64

	query := setup.DB.Preload("User").Order("created_at DESC")
	countQuery := setup.DB.Model(&models.Keluarga{})
	if !isAdmin(user) {
		query = query.Where("user_id = ?", user.Id)
		countQuery = countQuery.Where("user_id = ?", user.Id)
	}

	countQuery.Count(&total)

	// Pagination opsional: aktif hanya jika ?page= / ?limit= diberikan.
	if page, limit, ok := paginationParams(c); ok {
		query = query.Limit(limit).Offset((page - 1) * limit)
	}

	err := query.Find(&keluarga).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  "Gagal mengambil data: " + err.Error(),
		})
		return
	}

	if len(keluarga) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"data":    []map[string]interface{}{},
			"message": "Data kosong",
		})
		return
	}

	formattedKeluargas := make([]gin.H, len(keluarga))
	for i, keluarga := range keluarga {

		formattedKeluargas[i] = gin.H{
			"id":         keluarga.Id,
			"no_kk":      keluarga.NoKk,
			"kk_nama":    keluarga.KkNama,
			"alamat":     keluarga.Alamat,
			"rt":         keluarga.Rt,
			"rw":         keluarga.Rw,
			"kode_pos":   keluarga.KodePos,
			"status":     config.GetStatus(int(keluarga.Status)),
			"foto_kk":    keluarga.FotoKk,
			"foto_rumah": keluarga.FotoRumah,
			"latitude":   keluarga.Latitude,
			"longtitude": keluarga.Longtitude,
			"user_id":    keluarga.User.Name,
			"created_at": keluarga.CreatedAt,
			"updated_at": keluarga.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"data":  formattedKeluargas,
		"total": total,
	})
}

func Latest(c *gin.Context) {
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
	var keluarga []models.Keluarga

	if user.Level == "admin" {
		if err := setup.DB.Preload("User").Order("created_at DESC").Limit(5).Find(&keluarga).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
			return
		}
	} else {
		if err := setup.DB.Preload("User").Where("user_id = ?", user.Id).Order("created_at DESC").Limit(5).Find(&keluarga).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
			return
		}
	}

	if len(keluarga) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"data":    []map[string]interface{}{},
			"message": "Data kosong",
		})
		return
	}

	formattedKeluargas := make([]gin.H, len(keluarga))
	for i, keluarga := range keluarga {

		formattedKeluargas[i] = gin.H{
			"id":         keluarga.Id,
			"no_kk":      keluarga.NoKk,
			"kk_nama":    keluarga.KkNama,
			"alamat":     keluarga.Alamat,
			"rt":         keluarga.Rt,
			"rw":         keluarga.Rw,
			"kode_pos":   keluarga.KodePos,
			"status":     config.GetStatus(int(keluarga.Status)),
			"user_id":    keluarga.UserId,
			"foto_kk":    keluarga.FotoKk,
			"foto_rumah": keluarga.FotoRumah,
			"latitude":   keluarga.Latitude,
			"longtitude": keluarga.Longtitude,
			"created_at": keluarga.CreatedAt,
			"updated_at": keluarga.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
		"data": formattedKeluargas,
	})
}

func LatestForInput(c *gin.Context) {
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
	var keluarga []models.Keluarga

	if user.Level == "admin" {
		if err := setup.DB.Preload("User").Order("created_at DESC").Limit(1).Find(&keluarga).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
			return
		}
	} else {
		if err := setup.DB.Preload("User").Where("user_id = ?", user.Id).Order("created_at DESC").Limit(1).Find(&keluarga).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
			return
		}
	}

	if len(keluarga) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"data":    []map[string]interface{}{},
			"message": "Data kosong",
		})
		return
	}

	formattedKeluargas := make([]gin.H, len(keluarga))
	for i, keluarga := range keluarga {

		formattedKeluargas[i] = gin.H{
			"id":         keluarga.Id,
			"no_kk":      keluarga.NoKk,
			"kk_nama":    keluarga.KkNama,
			"alamat":     keluarga.Alamat,
			"rt":         keluarga.Rt,
			"rw":         keluarga.Rw,
			"kode_pos":   keluarga.KodePos,
			"status":     config.GetStatus(int(keluarga.Status)),
			"user_id":    keluarga.UserId,
			"foto_kk":    keluarga.FotoKk,
			"foto_rumah": keluarga.FotoRumah,
			"latitude":   keluarga.Latitude,
			"longtitude": keluarga.Longtitude,
			"created_at": keluarga.CreatedAt,
			"updated_at": keluarga.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
		"data": formattedKeluargas,
	})
}

func GetKeluargaByID(c *gin.Context) {
	id := c.Param("id")

	user, ok := getCurrentUser(c)
	if !ok {
		return
	}

	var keluarga models.Keluarga
	if err := setup.DB.First(&keluarga, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	if !canAccessResource(c, user, keluarga.UserId) {
		return
	}

	response := gin.H{
		"id":         keluarga.Id,
		"no_kk":      keluarga.NoKk,
		"kk_nama":    keluarga.KkNama,
		"alamat":     keluarga.Alamat,
		"rt":         keluarga.Rt,
		"rw":         keluarga.Rw,
		"kode_pos":   keluarga.KodePos,
		"status":     config.GetStatus(int(keluarga.Status)),
		"user_id":    keluarga.UserId,
		"foto_kk":    keluarga.FotoKk,
		"foto_rumah": keluarga.FotoRumah,
		"latitude":   keluarga.Latitude,
		"longtitude": keluarga.Longtitude,
		"created_at": keluarga.CreatedAt,
		"updated_at": keluarga.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func AddKeluarga(c *gin.Context) {

	user, ok := getCurrentUser(c)
	if !ok {
		return
	}

	var keluarga models.Keluarga

	noKk := c.PostForm("no_kk")
	if noKk == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No KK wajib diisi"})
		return
	}

	kkNama := c.PostForm("kk_nama")
	if kkNama == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama Kepala Keluarga harus diisi"})
		return
	}

	alamat := c.PostForm("alamat")
	if alamat == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Alamat harus diisi"})
		return
	}

	rt := c.PostForm("rt")
	if rt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "RT harus diisi"})
		return
	}

	rw := c.PostForm("rw")
	if rw == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "RW harus diisi"})
		return
	}

	kodePos := c.PostForm("kode_pos")
	if kodePos == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode Pos harus diisi"})
		return
	}

	status := c.PostForm("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status harus diisi"})
		return
	}

	// user_id diambil dari user yang sedang login, bukan dari input client
	userIdInt := int64(user.Id)

	latitude := c.PostForm("latitude")
	longtitude := c.PostForm("longtitude")

	var latitudeFloat, longtitudeFloat float64

	// no_kk & status sudah dipastikan tidak kosong di atas; di sini divalidasi formatnya.
	noKkInt, err := strconv.ParseInt(noKk, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format No KK tidak valid"})
		return
	}

	statusInt64, err := strconv.ParseInt(status, 10, 8)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format Status tidak valid"})
		return
	}
	statusInt := int8(statusInt64)

	// Koordinat bersifat opsional, tapi jika diisi harus berupa angka valid.
	if latitude != "" {
		latitudeFloat, err = strconv.ParseFloat(latitude, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format Latitude tidak valid"})
			return
		}
	}
	if longtitude != "" {
		longtitudeFloat, err = strconv.ParseFloat(longtitude, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format Longitude tidak valid"})
			return
		}
	}

	// Track file yang sudah ditulis ke disk; kalau tx gagal / return early
	// sebelum commit, defer akan menghapusnya supaya tidak jadi orphan.
	var savedPaths []string
	committed := false
	defer func() {
		if committed {
			return
		}
		for _, p := range savedPaths {
			_ = os.Remove(p)
		}
	}()

	// Foto bersifat opsional.
	if fotoKk, ferr := c.FormFile("foto_kk"); ferr == nil {
		uploadPath, verr := utils.ValidateAndBuildImagePath("public/uploads/foto-kk", fotoKk)
		if verr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": verr.Error()})
			return
		}
		if err := c.SaveUploadedFile(fotoKk, uploadPath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal menyimpan foto"})
			return
		}
		savedPaths = append(savedPaths, uploadPath)
		keluarga.FotoKk = uploadPath
	}

	if fotoRumah, ferr := c.FormFile("foto_rumah"); ferr == nil {
		uploadPath, verr := utils.ValidateAndBuildImagePath("public/uploads/foto-rumah", fotoRumah)
		if verr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": verr.Error()})
			return
		}
		if err := c.SaveUploadedFile(fotoRumah, uploadPath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal menyimpan foto"})
			return
		}
		savedPaths = append(savedPaths, uploadPath)
		keluarga.FotoRumah = uploadPath
	}

	newKeluarga := models.Keluarga{
		NoKk:       noKkInt,
		KkNama:     kkNama,
		Alamat:     alamat,
		Rt:         rt,
		Rw:         rw,
		KodePos:    kodePos,
		Status:     int8(statusInt),
		UserId:     userIdInt,
		FotoKk:     keluarga.FotoKk,
		FotoRumah:  keluarga.FotoRumah,
		Latitude:   latitudeFloat,
		Longtitude: longtitudeFloat,
	}

	tx := setup.DB.Begin()

	if err := tx.Create(&newKeluarga).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan program: " + err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data"})
		return
	}
	committed = true
	setup.DB.Preload("User").First(&newKeluarga, newKeluarga.Id)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Data berhasil ditambahkan",
		"data":    newKeluarga,
	})
}

func UpdateKeluarga(c *gin.Context) {
	id := c.Param("id")

	user, ok := getCurrentUser(c)
	if !ok {
		return
	}

	var keluarga models.Keluarga

	if err := setup.DB.First(&keluarga, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data keluarga tidak ditemukan"})
		return
	}

	if !canAccessResource(c, user, keluarga.UserId) {
		return
	}

	noKk := c.PostForm("no_kk")
	kkNama := c.PostForm("kk_nama")
	alamat := c.PostForm("alamat")
	rt := c.PostForm("rt")
	rw := c.PostForm("rw")
	kodePos := c.PostForm("kode_pos")
	status := c.PostForm("status")
	latitude := c.PostForm("latitude")
	longtitude := c.PostForm("longtitude")

	fotoKk, err := c.FormFile("foto_kk")
	if err == nil {
		uploadPath, verr := utils.ValidateAndBuildImagePath("public/uploads/foto-kk", fotoKk)
		if verr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": verr.Error()})
			return
		}
		if err := c.SaveUploadedFile(fotoKk, uploadPath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal menyimpan foto"})
			return
		}
		keluarga.FotoKk = uploadPath
	}

	FotoRumah, err := c.FormFile("foto_rumah")
	if err == nil {
		uploadPath, verr := utils.ValidateAndBuildImagePath("public/uploads/foto-rumah", FotoRumah)
		if verr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": verr.Error()})
			return
		}
		if err := c.SaveUploadedFile(FotoRumah, uploadPath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal menyimpan foto"})
			return
		}
		keluarga.FotoRumah = uploadPath
	}

	if fotoKk == nil {
		keluarga.FotoKk = ""
	}
	if FotoRumah == nil {
		keluarga.FotoRumah = ""
	}

	// Konversi nilai-nilai form ke tipe data yang sesuai
	var noKkInt int64
	var statusInt int8
	var latitudeFloat, longtitudeFloat float64

	if noKk != "" {
		noKkInt, err = strconv.ParseInt(noKk, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format No KK tidak valid"})
			return
		}
	}

	if status != "" {
		statusInt64, err := strconv.ParseInt(status, 10, 8)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format Status tidak valid"})
			return
		}
		statusInt = int8(statusInt64)
	}

	if latitude != "" {
		latitudeFloat, err = strconv.ParseFloat(latitude, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format Latitude tidak valid"})
			return
		}
	}

	if longtitude != "" {
		longtitudeFloat, err = strconv.ParseFloat(longtitude, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format Longitude tidak valid"})
			return
		}
	}

	tx := setup.DB.Begin()

	updateData := map[string]interface{}{}

	// Hanya menambahkan field yang memiliki nilai
	if noKk != "" {
		updateData["no_kk"] = noKkInt
	}
	if kkNama != "" {
		updateData["kk_nama"] = kkNama
	}
	if alamat != "" {
		updateData["alamat"] = alamat
	}
	if rt != "" {
		updateData["rt"] = rt
	}
	if rw != "" {
		updateData["rw"] = rw
	}
	if kodePos != "" {
		updateData["kode_pos"] = kodePos
	}
	if status != "" {
		updateData["status"] = statusInt
	}
	if keluarga.FotoKk != "" {
		updateData["foto_kk"] = keluarga.FotoKk
	}
	if keluarga.FotoRumah != "" {
		updateData["foto_rumah"] = keluarga.FotoRumah
	}
	if latitude != "" {
		updateData["latitude"] = latitudeFloat
	}
	if longtitude != "" {
		updateData["longtitude"] = longtitudeFloat
	}

	if err := tx.Model(&keluarga).Updates(updateData).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate data: " + err.Error()})
		return
	}

	tx.Commit()

	setup.DB.Preload("User").First(&keluarga, keluarga.Id)
	c.JSON(http.StatusOK, gin.H{
		"message": "Data keluarga berhasil diupdate",
		"data":    keluarga,
	})
}

func DeleteKeluarga(c *gin.Context) {
	id := c.Param("id")

	user, ok := getCurrentUser(c)
	if !ok {
		return
	}

	var keluarga models.Keluarga
	if err := setup.DB.First(&keluarga, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data keluarga tidak ditemukan",
		})
		return
	}

	if !canAccessResource(c, user, keluarga.UserId) {
		return
	}

	if err := setup.DB.Delete(&keluarga).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menghapus data keluarga",
		})
		return
	}

	// Hapus file foto yang terkait agar tidak menjadi sampah di disk.
	for _, p := range []string{keluarga.FotoKk, keluarga.FotoRumah} {
		if p != "" {
			_ = os.Remove(p)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data keluarga berhasil dihapus",
	})
}
