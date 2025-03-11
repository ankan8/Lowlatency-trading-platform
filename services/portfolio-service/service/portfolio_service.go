package service

import (
    "errors"
    "fmt"

    "github.com/ankan8/swapsync/backend/services/portfolio-service/models"
    "github.com/ankan8/swapsync/backend/services/portfolio-service/repository"
)

// GetPortfolio retrieves the portfolio document for a given user from the DB.
func GetPortfolio(userID string) (*models.Portfolio, error) {
    return repository.GetPortfolioByUserID(userID)
}

// UpdateHoldings updates the user's holdings for a given symbol and quantity.
// We allow negative quantity for SELL, but disallow quantity=0, and require price>0.
// If you want to ensure the user cannot go below zero shares, you'd add an additional check here.
func UpdateHoldings(userID, symbol string, quantity, price float64) error {
    // Disallow quantity=0
    if quantity == 0 {
        return errors.New("invalid quantity=0")
    }
    // Disallow price <= 0
    if price <= 0 {
        return errors.New("invalid price <= 0")
    }

    // If you want to ensure final holdings >= 0, you'd fetch existing holdings and check here:
    // existingHoldings, _ := repository.GetHoldingsForSymbol(userID, symbol)
    // if existingHoldings + quantity < 0 { return errors.New("not enough shares") }

    // Otherwise, just pass it along to the repository.
    err := repository.UpdateHoldings(userID, symbol, quantity, price)
    if err != nil {
        return err
    }
    fmt.Printf("Updated holdings for user %s, symbol %s\n", userID, symbol)
    return nil
}
