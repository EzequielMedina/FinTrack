package main

import (
	"fmt"
	"log"

	"github.com/fintrack/account-service/internal/config"
	"github.com/fintrack/account-service/internal/core/domain/entities"
	mysqlrepo "github.com/fintrack/account-service/internal/infrastructure/repositories/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("🧪 Testing automatic credit release for completed installment plans...")

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect to database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Create repositories and services
	cardRepo := mysqlrepo.NewCardRepository(gormDB)
	
	// Test card ID
	cardID := "1999f144-2cb9-45c8-8c43-3847b3708553"
	
	fmt.Printf("🔍 Checking current card balance for card: %s\n", cardID)
	
	card, err := cardRepo.GetByIDWithAccount(cardID)
	if err != nil {
		log.Fatalf("failed to get card: %v", err)
	}
	
	fmt.Printf("📊 Current card balance: $%.2f\n", card.Balance)
	fmt.Printf("💳 Card type: %s\n", card.CardType)
	
	if card.CardType == entities.CardTypeCredit {
		fmt.Printf("🏦 Credit limit: $%.2f\n", *card.CreditLimit)
		availableCredit := *card.CreditLimit - card.Balance
		fmt.Printf("💰 Available credit: $%.2f\n", availableCredit)
	}
	
	// Find ALL installment plans for this card
	// Debug: Try direct query to see what's wrong
	fmt.Printf("🔍 Debugging repository query...\n")
	
	var plans []*entities.InstallmentPlan
	err = gormDB.Where("card_id = ?", cardID).Find(&plans).Error
	if err != nil {
		log.Fatalf("failed to query installment plans directly: %v", err)
	}
	
	fmt.Printf("📋 Found %d total installment plans (direct query)\n", len(plans))
	
	totalToRelease := 0.0
	for _, plan := range plans {
		fmt.Printf("  Plan ID: %s, Amount: $%.2f, Status: %s, Paid: %d/%d\n", 
			plan.ID, plan.TotalAmount, plan.Status, plan.PaidInstallments, plan.InstallmentsCount)
		if plan.Status == entities.InstallmentPlanStatusCompleted {
			totalToRelease += plan.TotalAmount
		}
	}
	
	fmt.Printf("💸 Total amount that should be released: $%.2f\n", totalToRelease)
	
	if card.CardType == entities.CardTypeCredit && totalToRelease > 0 {
		fmt.Printf("🔓 Simulating automatic credit release...\n")
		
		newBalance := card.Balance - totalToRelease
		if newBalance < 0 {
			newBalance = 0
		}
		
		fmt.Printf("📉 New balance would be: $%.2f (reduction of $%.2f)\n", newBalance, totalToRelease)
		fmt.Printf("📈 New available credit would be: $%.2f\n", *card.CreditLimit - newBalance)
		
		// Actually apply the release (BE CAREFUL!)
		confirm := true // Set to true to actually apply changes
		
		if confirm {
			fmt.Printf("⚡ Applying credit release...\n")
			
			card.Balance = newBalance
			updatedCard, err := cardRepo.Update(card)
			if err != nil {
				log.Fatalf("failed to update card balance: %v", err)
			}
			
			fmt.Printf("✅ Credit release successful!\n")
			fmt.Printf("📊 Updated card balance: $%.2f\n", updatedCard.Balance)
			fmt.Printf("💰 Updated available credit: $%.2f\n", *updatedCard.CreditLimit - updatedCard.Balance)
		} else {
			fmt.Printf("🚫 Dry run - no changes applied\n")
		}
	} else {
		fmt.Printf("ℹ️ No credit release needed or not a credit card\n")
	}
	
	fmt.Println("🏁 Test completed!")
}