package entity

import "time"

type User struct {
	ID        string    `db:"id" json:"ID"`
	FullName  string    `db:"full_name" json:"fullName"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type UserEmail struct {
	ID        string    `db:"id" json:"ID"`
	UserID    string    `db:"user_id" json:"UserID"`
	Email     string    `db:"email" json:"email"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type UserPassword struct {
	ID        string    `db:"id" json:"ID"`
	UserID    string    `db:"user_id" json:"UserID"`
	Password  string    `db:"hash" json:"password"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
