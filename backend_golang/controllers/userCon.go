package controllers

import (
	"backend_golang/models"
	"backend_golang/setup"
	"backend_golang/utils"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

func GetAllUser(c *gin.Context) {
	user, ok := getCurrentUser(c)
	if !ok {
		return
	}

	// Hanya admin yang boleh melihat daftar seluruh user.
	if !isAdmin(user) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses", "status": false})
		return
	}

	var users []models.User

	// Order harus dipanggil sebelum Find agar tidak diabaikan GORM.
	if err := setup.DB.Order("name ASC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "status": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func UpdateUser(c *gin.Context) {
	userId := c.Param("id")
	var user models.User

	// Cek user yang akan diedit
	if err := setup.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	// Ambil data form
	name := c.PostForm("name")
	email := c.PostForm("email")
	profilePicture := c.PostForm("profile_picture")
	level := c.PostForm("level")

	// Update data user
	updates := map[string]interface{}{}

	if name != "" {
		updates["name"] = name
	}
	if email != "" {
		updates["email"] = email
	}
	if profilePicture != "" {
		updates["profile_picture"] = profilePicture
	}

	if level != "" {
		updates["level"] = level
	}

	// Handle foto profil jika ada
	file, err := c.FormFile("profile_picture")
	if err == nil {
		// Validasi tipe (berbasis MIME) & ukuran, lalu buat nama file acak yang aman
		uploadPath, verr := utils.ValidateAndBuildImagePath("public/uploads/profile_pictures", file)
		if verr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": verr.Error()})
			return
		}

		// Simpan file
		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan foto profil"})
			return
		}

		// Simpan path lengkap ke database
		updates["profile_picture"] = uploadPath
	}

	if err := setup.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui pengguna"})
		return
	}

	// Reload user data dengan relasi
	setup.DB.First(&user, userId)

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diperbarui", "data": user})
}

func PasswordUpdate(c *gin.Context) {
	userId := c.Param("id")
	var user models.User

	// Cek user yang akan diedit
	if err := setup.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	// Ambil data form
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")

	if password != confirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password dan konfirmasi password tidak cocok"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Update data user
	updates := map[string]interface{}{}

	if password != "" {
		updates["password"] = string(hashedPassword)
	}

	if err := setup.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui pengguna"})
		return
	}

	// Reload user data dengan relasi
	setup.DB.First(&user, userId)

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diperbarui", "data": user})
}
