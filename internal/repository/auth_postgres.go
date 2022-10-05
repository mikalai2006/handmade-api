package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mikalai2006/handmade/internal/domain"
)

type AuthPostgres struct {
	db *sqlx.DB
}


func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db:db}
}

func (r *AuthPostgres) CreateUser(user domain.User) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash, vk_id, google_id, github_id, apple_id) values ($1,$2,$3,$4,$5,$6,$7) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password, user.VkId, user.GoogleId, user.GithubId, user.AppleId)
	if err := row.Scan(&id); err != nil {
		return 0, nil
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}

func (r *AuthPostgres) Logout(username string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("UPDATE username FROM %s WHERE username=$1", usersTable)
	err := r.db.Get(&user, query, username)

	return user, err
}

func (r *AuthPostgres) GetByCredentials(ctx context.Context, user domain.User) (domain.User, error)  {
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, user.Username, user.Password)
	return user, err
}