package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/tsaqiffatih/auth-service/models"
	"github.com/tsaqiffatih/auth-service/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var validate *validator.Validate

// Initialize validator instance
func init() {
	validate = validator.New()
}

func RegisterHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "fail",
				"error": map[string]interface{}{
					"code":    http.StatusBadRequest,
					"message": "Invalid request payload",
				},
			})
			return
		}

		// Validation using validator from model
		err = validate.Struct(user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "fail",
				"error": map[string]interface{}{
					"code":    http.StatusBadRequest,
					"message": err.Error(),
				},
			})
			return
		}

		//cheking is user already exist?
		var existingUser models.User
		result := db.Where("email = ?", user.Email).First(&existingUser)
		if result.RowsAffected > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "fail",
				"error": map[string]interface{}{
					"code":    http.StatusBadRequest,
					"message": "Email already exists",
				},
			})
			return
		}

		//Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "fail",
				"error": map[string]interface{}{
					"code":    http.StatusInternalServerError,
					"message": "Error hashing password",
				},
			})
			return
		}
		user.Password = string(hashedPassword)

		// Save user to DB
		result = db.Create(&user)
		if result.Error != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "fail",
				"error": map[string]interface{}{
					"code":    http.StatusBadRequest,
					"message": "User already exists",
				},
			})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "User registered successfully",
		})
	}
}

func LoginHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var creds struct {
			Email    string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required"`
		}

		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "fail",
				"error": map[string]interface{}{
					"code":    http.StatusBadRequest,
					"message": "Invalid request payload",
				},
			})
			return
		}

		// validation input login
		err = validate.Struct(creds)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "fail",
				"error": map[string]interface{}{
					"code":    http.StatusBadRequest,
					"message": err.Error(),
				},
			})
			return
		}

		// find user with existing users email
		var user models.User
		result := db.Where("email = ?", creds.Email).First(&user)
		if result.Error != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "fail",
				"error": map[string]interface{}{
					"code":    http.StatusUnauthorized,
					"message": "User not found",
				},
			})
			return
		}

		//Check password using bcrypt
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "fail",
				"error": map[string]interface{}{
					"code":    http.StatusUnauthorized,
					"message": "Invalid password",
				},
			})
			return
		}

		// Encrypt user.ID before generating JWT
		encryptedID, err := utils.EncryptID(user.ID)

		//Generate JWT token
		tokenString, err := utils.GenerateJWT(encryptedID, user.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "fail",
				"error": map[string]interface{}{
					"code":    http.StatusInternalServerError,
					"message": "Error generating token",
				},
			})
			return
		}

		//Optionally store token in DB
		db.Create(&models.Token{UserID: user.ID, Token: tokenString})

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Login successful",
			"payload": map[string]string{
				"token": tokenString,
			},
		})
	}
}

func LogoutHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		token := r.Header.Get("Authorization")[7:] //Strip "Bearer "
		result := db.Where("token = ?", token).Delete(&models.Token{})
		if result.Error != nil || result.RowsAffected == 0 {

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "fail",
				"error": map[string]interface{}{
					"code":    http.StatusUnauthorized,
					"message": "Token not found or already invalidated",
				},
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Logged out successfully",
		})
	}
}
