package services

import (
	"avito-shop/models"
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

func GetOrCreateUser(username, password string) (*models.User, error) {
	var user models.User
	err := db.QueryRow("SELECT id, username, password, coins FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Coins)
	if err != nil {
		if err == sql.ErrNoRows {
			err = db.QueryRow("INSERT INTO users (username, password, coins) VALUES ($1, $2, $3) RETURNING id, username, password, coins",
				username, password, 1000).Scan(&user.ID, &user.Username, &user.Password, &user.Coins)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &user, nil
}

func GetUserInfo(username string) (*models.InfoResponse, error) {
	var info models.InfoResponse

	err := db.QueryRow("SELECT coins FROM users WHERE username = $1", username).Scan(&info.Coins)
	if err != nil {
		return nil, err
	}

	invRows, err := db.Query("SELECT item_type, quantity FROM inventory WHERE user_id = (SELECT id FROM users WHERE username = $1)", username)
	if err != nil {
		return nil, err
	}
	defer invRows.Close()
	for invRows.Next() {
		var item models.InventoryItem
		if err := invRows.Scan(&item.ItemType, &item.Quantity); err != nil {
			return nil, err
		}
		info.Inventory = append(info.Inventory, item)
	}

	receivedRows, err := db.Query("SELECT from_user, amount FROM coin_transactions WHERE to_user = $1 ORDER BY created_at DESC", username)
	if err != nil {
		return nil, err
	}
	defer receivedRows.Close()
	for receivedRows.Next() {
		var rec models.TransactionRecord
		if err := receivedRows.Scan(&rec.FromUser, &rec.Amount); err != nil {
			return nil, err
		}
		info.CoinHistory.Received = append(info.CoinHistory.Received, rec)
	}

	sentRows, err := db.Query("SELECT to_user, amount FROM coin_transactions WHERE from_user = $1 ORDER BY created_at DESC", username)
	if err != nil {
		return nil, err
	}
	defer sentRows.Close()
	for sentRows.Next() {
		var rec models.TransactionRecord
		if err := sentRows.Scan(&rec.ToUser, &rec.Amount); err != nil {
			return nil, err
		}
		info.CoinHistory.Sent = append(info.CoinHistory.Sent, rec)
	}

	return &info, nil
}

func SendCoins(fromUsername, toUsername string, amount int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var fromCoins int64
	err = tx.QueryRow("SELECT coins FROM users WHERE username = $1 FOR UPDATE", fromUsername).Scan(&fromCoins)
	if err != nil {
		return err
	}

	if fromCoins < amount {
		return errors.New("недостаточно монет")
	}

	_, err = tx.Exec("UPDATE users SET coins = coins - $1 WHERE username = $2", amount, fromUsername)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE users SET coins = coins + $1 WHERE username = $2", amount, toUsername)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO coin_transactions (from_user, to_user, amount) VALUES ($1, $2, $3)", fromUsername, toUsername, amount)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func BuyItem(username, item string) error {
	price, ok := MerchItems[item]
	if !ok {
		return errors.New("товар не найден")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var coins int64
	err = tx.QueryRow("SELECT coins FROM users WHERE username = $1 FOR UPDATE", username).Scan(&coins)
	if err != nil {
		return err
	}
	if coins < price {
		return errors.New("недостаточно монет для покупки")
	}

	_, err = tx.Exec("UPDATE users SET coins = coins - $1 WHERE username = $2", price, username)
	if err != nil {
		return err
	}

	var count int64
	err = tx.QueryRow("SELECT quantity FROM inventory WHERE user_id = (SELECT id FROM users WHERE username = $1) AND item_type = $2", username, item).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = tx.Exec("INSERT INTO inventory (user_id, item_type, quantity) VALUES ((SELECT id FROM users WHERE username = $1), $2, 1)", username, item)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		_, err = tx.Exec("UPDATE inventory SET quantity = quantity + 1 WHERE user_id = (SELECT id FROM users WHERE username = $1) AND item_type = $2", username, item)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
