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

func (r *PostgresRepository) GetPostById(ctx context.Context, id string) (*models.Post, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, post_content, user_id, created_at FROM posts WHERE id = $1", id)
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println("failed to close rows")
		}
	}()

	if err != nil {
		return nil, err
	}

	post := &models.Post{}
	for rows.Next() {
		err := rows.Scan(&post.Id, &post.PostContent, &post.UserId, &post.CreatedAt)
		if err != nil {
			return nil, err
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}

		return post, nil
	}

	return post, nil
}

func (r *PostgresRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	_, err := r.db.ExecContext(ctx, "UPDATE posts SET post_content = $1 WHERE id = $2 and user_id = $3", post.PostContent, post.Id, post.UserId)
	return err
}

func (r *PostgresRepository) DeletePost(ctx context.Context, id string, userId string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM posts WHERE id = $1 AND user_id = $2", id, userId)
	return err
}

func (r *PostgresRepository) ListPosts(ctx context.Context, limit uint64, page uint64) ([]*models.Post, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, post_content, user_id, created_at FROM posts ORDER BY created_at DESC LIMIT $1 OFFSET $2", limit, page*limit)
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println("failed to close rows")
		}
	}()

	if err != nil {
		return nil, err
	}

	posts := make([]*models.Post, 0)
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(&post.Id, &post.PostContent, &post.UserId, &post.CreatedAt)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
