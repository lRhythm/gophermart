package pg

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/lRhythm/gophermart/internal/models"
	"github.com/lib/pq"
	"time"
)

// User.

func (db *DB) UserCreateAndIDGet(login, passwordHash string) (string, error) {
	var u user
	query := `insert into users (login, password_hash) values ($1, $2) returning id`
	err := db.con.Get(&u, query, login, passwordHash)
	var errPq *pq.Error
	if err != nil && errors.As(err, &errPq) && errPq.Code == pgerrcode.UniqueViolation {
		err = models.ErrorStorageConflict
	}
	return u.ID, err
}

func (db *DB) UserIDAndPasswordHashByLoginGet(login string) (string, string, error) {
	var u user
	query := `select id, password_hash from users where login = $1`
	err := db.con.Get(&u, query, login)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		err = models.ErrorStorageNoRows
	}
	return u.ID, u.PasswordHash, err
}

func (db *DB) UserBalanceAndWithdrawnGet(userID string) (float32, float32, error) {
	var u user
	query := `select balance, withdrawn from users where id = $1`
	err := db.con.Get(&u, query, userID)
	return s2r(u.Balance), s2r(u.Withdrawn), err
}

func (db *DB) UserBalanceGet(userID string) (float32, error) {
	var u user
	query := `select balance, withdrawn from users where id = $1`
	err := db.con.Get(&u, query, userID)
	return s2r(u.Balance), err
}

func (db *DB) UserBalanceWithdraw(orderID string, sum float32, userID string, processedAt time.Time) error {
	tx, err := db.con.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `insert into withdrawals (user_id, order_id, sum, processed_at) values ($1, $2, $3, $4)`
	_, err = tx.Exec(query, userID, orderID, r2s(sum), processedAt)
	if err != nil {
		return err
	}
	query = `update users set balance = balance - $1, withdrawn = withdrawn + $2 where id = $3`
	s := r2s(sum)
	_, err = tx.Exec(query, s, s, userID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (db *DB) UserWithdrawalList(userID string) ([]models.WithdrawalDTO, error) {
	ws := make(withdrawals, 0)
	query := `select order_id, sum, processed_at from withdrawals where user_id = $1 order by processed_at`
	err := db.con.Select(&ws, query, userID)
	return ws.toDTO(), err
}

// Order.

func (db *DB) OrderUserIDGetAndInsert(orderID, userID string, uploadedAt time.Time) (bool, string, error) {
	var exists bool // Существует ли заказ.
	tx, err := db.con.Beginx()
	if err != nil {
		return exists, "", err
	}
	defer tx.Rollback()
	// Выборка user.id заказа, если заказ существует.
	var o order
	query := `select user_id from orders where id = $1`
	err = tx.Get(&o, query, orderID)
	if err != nil {
		// Если заказ не существует - вставка записи заказа.
		if errors.Is(err, sql.ErrNoRows) {
			query = `insert into orders (id, user_id, uploaded_at) values ($1, $2, $3)`
			_, err = tx.Exec(query, orderID, userID, uploadedAt)
			if err != nil {
				return exists, "", err
			}
			// Commit только при условии, что заказ создан успешно.
			if err = tx.Commit(); err != nil {
				return exists, "", err
			}
			return exists, "", nil
		}
		// Прочая ошибка.
		return exists, "", err
	}
	// Заказ существует.
	exists = true
	return exists, o.UserID, err
}

// Без транзакции.
//func (db *DB) OrderUserIDGet(orderID string) (string, error) {
//	var o order
//	query := `select user_id from orders where id = $1`
//	err := db.con.Get(&o, query, orderID)
//	if err != nil && errors.Is(err, sql.ErrNoRows) {
//		err = models.ErrorStorageNoRows
//	}
//	return o.UserID, err
//}
//
//func (db *DB) OrderInsert(orderID, userID string, uploadedAt time.Time) error {
//	query := `insert into orders (id, user_id, uploaded_at) values ($1, $2, $3)`
//	_, err := db.con.Exec(query, orderID, userID, uploadedAt)
//	return err
//}

func (db *DB) OrderList(userID string) ([]models.OrderDTO, error) {
	os := make(orders, 0)
	query := `select id, status, accrual, uploaded_at from orders where user_id = $1 order by uploaded_at`
	err := db.con.Select(&os, query, userID)
	return os.toDTO(), err
}

func (db *DB) OrdersIDForProcessing(statuses []string) ([]string, error) {
	ordersID := make([]string, 0)
	query := `update orders set in_processing = true where in_processing = false and status = any($1) returning id`
	err := db.con.Select(&ordersID, query, pq.Array(statuses))
	return ordersID, err
}

func (db *DB) OrderUnlockProcessing(orderID string) error {
	query := `update orders set in_processing = false where id = $1`
	_, err := db.con.Exec(query, orderID)
	return err
}

func (db *DB) OrderStatusAndAccrualSet(orderID string, status string, accrual float32) error {
	// Если начисления == 0, то не нужно обновлять баланс пользователя, т.е. транзакция не нужна.
	if accrual == 0 {
		query := `update orders set status = $1, in_processing = false where id = $2`
		_, err := db.con.Exec(query, status, orderID)
		return err
	}

	tx, err := db.con.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var userID string
	query := `update orders set status = $1, accrual = $2, in_processing = false where id = $3 returning user_id`
	err = tx.Get(&userID, query, status, r2s(accrual), orderID)
	if err != nil {
		return err
	}
	query = `update users set balance = balance + $1 where id = $2`
	_, err = tx.Exec(query, r2s(accrual), userID)
	if err != nil {
		return err
	}
	return tx.Commit()
}
