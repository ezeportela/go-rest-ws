package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/ezeportela/go-rest-ws/models"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db}, nil
}

func (r *PostgresRepository) Close() {
	r.db.Close()
}

func (r *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)", user.Id, user.Email, user.Password)
	return err
}

func (r *PostgresRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, email FROM users WHERE id = $1", id)
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println("failed to close rows")
		}
	}()

	if err != nil {
		return nil, err
	}

	user := &models.User{}
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Email)
		if err != nil {
			return nil, err
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}

		return user, nil
	}

	return user, nil
}

func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, email, password FROM users WHERE email = $1", email)
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println("failed to close rows")
		}
	}()

	if err != nil {
		return nil, err
	}

	user := &models.User{}
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}

		return user, nil
	}

	return user, nil
}

func (r *PostgresRepository) InsertPost(ctx context.Context, post *models.Post) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO posts (id, post_content, user_id) VALUES ($1, $2, $3)", post.Id, post.PostContent, post.UserId)
	return err
}
