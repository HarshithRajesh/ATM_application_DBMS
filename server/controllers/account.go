package controllers

import (
  "strconv"
  "time"
  "github.com/HarshithRajesh/zapster/initializers"
  "github.com/HarshithRajesh/zapster/models"
	"net/http"
	"github.com/gin-gonic/gin"
)

func CreateAccount(c *gin.Context) {
	var input struct {
		UserID      uint   `json:"user_id" binding:"required"`
		CardID      uint   `json:"card_id" binding:"required"`
		AccountType string `json:"account_type" binding:"required,oneof=Savings Current"`
	}

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create account
	account := models.Account{
		UserID:               input.UserID,
		CardID:               input.CardID,
		AccountType:          input.AccountType,
		Balance:              0.00,
		DailyWithdrawalLimit: 1000.00,
	}

	if result := initializers.DB.Create(&account); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account created successfully", "account": account})
}

func GetAccount(c *gin.Context) {
accountID, err := strconv.Atoi(c.Param("id"))
if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
    return
}

	var account models.Account
	if result := initializers.DB.Preload("User").Preload("Card").First(&account, accountID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": account})
}

func UpdateAccount(c *gin.Context) {
	accountID := c.Param("id")

	var input struct {
		AccountType          string  `json:"account_type" binding:"omitempty,oneof=Savings Current"`
		DailyWithdrawalLimit float64 `json:"daily_withdrawal_limit" binding:"omitempty"`
	}

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find account
	var account models.Account
	if result := initializers.DB.First(&account, accountID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	// Update account
	initializers.DB.Model(&account).Updates(input)
	c.JSON(http.StatusOK, gin.H{"message": "Account updated successfully", "account": account})
}
func DeleteAccount(c *gin.Context) {
	accountID := c.Param("id")

	if result := initializers.DB.Delete(&models.Account{}, accountID); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}
func ListAccounts(c *gin.Context) {
	userID := c.Query("user_id") // Assume user_id is passed as a query parameter

	var accounts []models.Account
	if result := initializers.DB.Where("user_id = ?", userID).Find(&accounts); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No accounts found for the user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accounts": accounts})
}

func CreateCard(c *gin.Context) {
	var input struct {
		UserID     uint   `json:"user_id" binding:"required"`
		CardNumber string `json:"card_number" binding:"required"`
		PinHash    string `json:"pin_hash" binding:"required"`
		ExpiryDate string `json:"expiry_date" binding:"required"`
	}

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse the ExpiryDate string into time.Time
	expiryDate, err := time.Parse("2006-01-02", input.ExpiryDate) // Assuming format "YYYY-MM-DD"
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}
var existingCard models.Card
if err := initializers.DB.Where("card_number = ?", input.CardNumber).First(&existingCard).Error; err == nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Card with this number already exists"})
    return
}
// Assuming input.UserID is passed to create a card for that specific user
var user models.User
if err := initializers.DB.First(&user, input.UserID).Error; err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
    return
}

// Create card
card := models.Card{
    UserID:     input.UserID,
    CardNumber: input.CardNumber,
    PinHash:    input.PinHash, // In production, hash the pin securely
    ExpiryDate: expiryDate,
}

if result := initializers.DB.Create(&card); result.Error != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
    return
}

// Fetch the created card with the User details
var createdCard models.Card
if err := initializers.DB.Preload("User").First(&createdCard, card.ID).Error; err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load user details"})
    return
}
createdCard.User.Password = ""

c.JSON(http.StatusOK, gin.H{"message": "Card created successfully", "card": createdCard})

}
func UpdateCardStatus(c *gin.Context) {
	var input struct {
		CardNumber string `json:"card_number" binding:"required"`
		Status     string `json:"status" binding:"required,oneof=active blocked"`
	}

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the card using card number
	var card models.Card
	if err := initializers.DB.Preload("User").Where("card_number = ?", input.CardNumber).First(&card).Error; err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Card not found"})
		return
	}

	// Update card status
	card.CardStatus = input.Status
	if result := initializers.DB.Save(&card); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
  card.User.Password = ""

	c.JSON(http.StatusOK, gin.H{"message": "Card status updated successfully", "card": card})
}
func LockCardAfterFailedAttempts(c *gin.Context) {
	const maxFailedAttempts = 3

	// Check card status and failed attempts
	var input struct {
		CardNumber string `json:"card_number" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var card models.Card
	if err := initializers.DB.Preload("User").Where("card_number = ?", input.CardNumber).First(&card).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Card not found"})
		return
	}

	// Lock the card if failed attempts exceed the limit
	if card.FailedAttempts >= maxFailedAttempts {
		card.CardStatus = "blocked"
		if err := initializers.DB.Save(&card).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Card blocked due to too many failed attempts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card is active"})
}
func DeleteCard(c *gin.Context) {
	cardNumber := c.Param("card_number")

	// Find the card using card number
	var card models.Card
	if err := initializers.DB.Where("card_number = ?", cardNumber).First(&card).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Card not found"})
		return
	}

	// Delete the card
	if err := initializers.DB.Delete(&card).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete card"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card deleted successfully"})
}

