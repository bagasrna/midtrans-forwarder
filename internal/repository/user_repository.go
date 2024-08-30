package repository

import (
    "context"
    "database/sql"

    "midtrans-forwarder/internal/domain"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
    query := "INSERT INTO users (email, password, name) VALUES (?, ?, ?)"
    result, err := r.db.ExecContext(ctx, query, user.Email, user.Password, user.Name)
    if err != nil {
        return err
    }
    id, err := result.LastInsertId()
    if err != nil {
        return err
    }
    user.ID = id
    return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
    query := "SELECT id, email, password, name FROM users WHERE id = ?"
    row := r.db.QueryRowContext(ctx, query, id)

    var user domain.User
    err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Name)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
    query := "SELECT id, email, name FROM users"
    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []domain.User
    for rows.Next() {
        var user domain.User
        if err := rows.Scan(&user.ID, &user.Email, &user.Name); err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    return users, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
    query := "UPDATE users SET email = ?, name = ? WHERE id = ?"
    _, err := r.db.ExecContext(ctx, query, user.Email, user.Name, user.ID)
    return err
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int64) error {
    query := "DELETE FROM users WHERE id = ?"
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}