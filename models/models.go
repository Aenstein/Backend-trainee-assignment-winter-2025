package models

// Пользователь
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Coins    int64  `json:"coins"`
}

// Запрос на аутентификацию.
type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Ответ при аутентификации.
type AuthResponse struct {
	Token string `json:"token"`
}

// Запрос на перевод монет.
type SendCoinRequest struct {
	ToUser string `json:"toUser" binding:"required"`
	Amount int64  `json:"amount" binding:"required"`
}

// Элемент инвентаря.
type InventoryItem struct {
	ItemType string `json:"type"`
	Quantity int64  `json:"quantity"`
}

// Запись перевода монет.
type TransactionRecord struct {
	FromUser string `json:"fromUser,omitempty"`
	ToUser   string `json:"toUser,omitempty"`
	Amount   int64  `json:"amount"`
}

// Ответ эндпоинта /api/info.
type InfoResponse struct {
	Coins       int64             `json:"coins"`
	Inventory   []InventoryItem   `json:"inventory"`
	CoinHistory struct {
		Received []TransactionRecord `json:"received"`
		Sent     []TransactionRecord `json:"sent"`
	} `json:"coinHistory"`
}
