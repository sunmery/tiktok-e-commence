// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type OrdersOrderItems struct {
	ID        int32          `json:"id"`
	OrderID   int32          `json:"orderID"`
	ProductID int32          `json:"productID"`
	Name      string         `json:"name"`
	Quantity  int32          `json:"quantity"`
	Price     pgtype.Numeric `json:"price"`
}

type OrdersOrders struct {
	ID            int32     `json:"id"`
	Owner         string    `json:"owner"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	StreetAddress string    `json:"streetAddress"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	Country       string    `json:"country"`
	ZipCode       string    `json:"zipCode"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
