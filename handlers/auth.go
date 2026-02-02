package handlers

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// PasswordLookup returns the bcrypt hash stored in DB for the admin password.
type PasswordLookup func() (hash string, err error)

// PasswordUpdate stores a new bcrypt hash in DB.
type PasswordUpdate func(hash string) error

type AuthHandler struct {
	Secret         string
	GetPasswordHash PasswordLookup
	SetPasswordHash PasswordUpdate
}

type loginRequest struct {
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	hash, err := h.GetPasswordHash()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin": true,
		"exp":   time.Now().Add(72 * time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(h.Secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenStr})
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "new password cannot be empty"})
		return
	}
	hash, err := h.GetPasswordHash()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "old password is incorrect"})
		return
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	if err := h.SetPasswordHash(string(newHash)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password changed"})
}

// InitAdminPassword ensures an admin password hash exists in DB.
// Returns PasswordLookup and PasswordUpdate functions.
func InitAdminPassword(db *gorm.DB) (PasswordLookup, PasswordUpdate) {
	type Setting struct {
		Key   string `gorm:"primarykey"`
		Value string
	}

	lookup := func() (string, error) {
		var s Setting
		if err := db.Where("key = ?", "admin_password_hash").First(&s).Error; err != nil {
			return "", err
		}
		return s.Value, nil
	}

	update := func(hash string) error {
		return db.Where("key = ?", "admin_password_hash").Assign(Setting{Value: hash}).FirstOrCreate(&Setting{Key: "admin_password_hash"}).Error
	}

	// Check if hash already exists
	var s Setting
	err := db.Where("key = ?", "admin_password_hash").First(&s).Error
	if err == nil && s.Value != "" {
		// Password already in DB, nothing to do
		return lookup, update
	}

	// Determine initial password
	var plainPassword string
	if envPw := os.Getenv("ADMIN_PASSWORD"); envPw != "" {
		plainPassword = envPw
	} else {
		plainPassword = generateRandomPassword(12)
		fmt.Println("========================================")
		fmt.Printf("🔑 Initial admin password: %s\n", plainPassword)
		fmt.Println("========================================")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash admin password: %v", err)
	}
	if err := update(string(hash)); err != nil {
		log.Fatalf("Failed to save admin password hash: %v", err)
	}

	return lookup, update
}

func generateRandomPassword(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			log.Fatalf("Failed to generate random password: %v", err)
		}
		result[i] = chars[n.Int64()]
	}
	return string(result)
}
