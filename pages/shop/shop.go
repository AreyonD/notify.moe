package shop

import (
	"net/http"
	"sort"
	"sync"

	"github.com/animenotifier/arn"

	"github.com/aerogo/aero"
	"github.com/animenotifier/notify.moe/components"
	"github.com/animenotifier/notify.moe/utils"
)

var itemBuyMutex sync.Mutex

// Get shop page.
func Get(ctx *aero.Context) string {
	user := utils.GetUser(ctx)

	if user == nil {
		return ctx.Error(http.StatusUnauthorized, "Not logged in", nil)
	}

	items, err := arn.AllItems()

	if err != nil {
		return ctx.Error(http.StatusInternalServerError, "Error fetching shop item data", err)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Order < items[j].Order
	})

	return ctx.HTML(components.Shop(user, items))
}

// BuyItem ...
func BuyItem(ctx *aero.Context) string {
	// Lock via mutex to prevent race conditions
	itemBuyMutex.Lock()
	defer itemBuyMutex.Unlock()

	// Logged in user
	user := utils.GetUser(ctx)

	if user == nil {
		return ctx.Error(http.StatusUnauthorized, "Not logged in", nil)
	}

	// Item ID and quantity
	itemID := ctx.Get("item")
	quantity, err := ctx.GetInt("quantity")

	if err != nil || quantity == 0 {
		return ctx.Error(http.StatusBadRequest, "Invalid item quantity", err)
	}

	item, err := arn.GetItem(itemID)

	if err != nil {
		return ctx.Error(http.StatusInternalServerError, "Error fetching item data", err)
	}

	// Calculate total price and subtract balance
	totalPrice := int(item.Price) * quantity

	if user.Balance < totalPrice {
		return ctx.Error(http.StatusBadRequest, "Not enough gems", nil)
	}

	user.Balance -= totalPrice
	err = user.Save()

	if err != nil {
		return ctx.Error(http.StatusInternalServerError, "Error saving user data", err)
	}

	// Add item to user inventory
	inventory := user.Inventory()
	inventory.AddItem(itemID, uint(quantity))
	err = inventory.Save()

	if err != nil {
		return ctx.Error(http.StatusInternalServerError, "Error saving inventory", err)
	}

	return "ok"
}