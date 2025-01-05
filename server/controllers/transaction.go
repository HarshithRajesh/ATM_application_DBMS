package controllers
import (
  "time"
  "github.com/HarshithRajesh/zapster/initializers"
  "github.com/HarshithRajesh/zapster/models"
	"net/http"
	"github.com/gin-gonic/gin"
  "fmt"
)
func CashWithdrawal(c *gin.Context) {
    const dailyWithdrawalLimit = 10000.00

    var input struct {
        AccountID uint    `json:"account_id" binding:"required"`
        Amount    float64 `json:"amount" binding:"required,gt=0"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Find account
    var account models.Account
    if err := initializers.DB.First(&account, input.AccountID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found"})
        return
    }

    // Check sufficient balance
    if account.Balance < input.Amount {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
        return
    }

    // Check daily limit
    var dailyLimit models.DailyLimit
    today := time.Now().Format("2006-01-02")
    if err := initializers.DB.Where("account_id = ? AND date = ?", input.AccountID, today).First(&dailyLimit).Error; err != nil {
        dailyLimit = models.DailyLimit{AccountID: input.AccountID, Date: time.Now()}
    }

    if dailyLimit.TotalWithdrawal+input.Amount > dailyWithdrawalLimit {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Daily withdrawal limit exceeded"})
        return
    }

    // Deduct balance and update daily limit
    account.Balance -= input.Amount
    if err := initializers.DB.Save(&account).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account balance"})
        return
    }

    dailyLimit.TotalWithdrawal += input.Amount
    dailyLimit.TotalTransactions++
    if err := initializers.DB.Save(&dailyLimit).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update daily limit"})
        return
    }

    // Record transaction
    transaction := models.Transaction{
        AccountID:       input.AccountID,
        TransactionType: "withdrawal",
        Amount:          input.Amount,
        Status:          "success",
        ReferenceNumber: fmt.Sprintf("TXN-%d", time.Now().Unix()),
    }

    if err := initializers.DB.Create(&transaction).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record transaction"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful", "transaction": transaction})
}
func CashDeposit(c *gin.Context) {
    var input struct {
        AccountID uint    `json:"account_id" binding:"required"`
        Amount    float64 `json:"amount" binding:"required,gt=0"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Find account
    var account models.Account
    if err := initializers.DB.First(&account, input.AccountID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found"})
        return
    }

    // Update balance
    account.Balance += input.Amount
    if err := initializers.DB.Save(&account).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account balance"})
        return
    }

    // Record transaction
    transaction := models.Transaction{
        AccountID:       input.AccountID,
        TransactionType: "deposit",
        Amount:          input.Amount,
        Status:          "success",
        ReferenceNumber: fmt.Sprintf("TXN-%d", time.Now().Unix()),
    }

    if err := initializers.DB.Create(&transaction).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record transaction"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Deposit successful", "transaction": transaction})
}
func TransactionHistory(c *gin.Context) {
    accountID := c.Query("account_id")

    // Fetch transactions
    var transactions []models.Transaction
    if err := initializers.DB.Where("account_id = ?", accountID).Order("created_at desc").Find(&transactions).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction history"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

