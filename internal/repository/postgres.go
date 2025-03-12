package repository

import (
	"30.8.1/internal/core"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(connectionString string) (*postgresRepository, error) {
	db, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}
	return &postgresRepository{db: db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) GetTasks(taskID, authorID int) ([]core.Task, error) {
	rows, err := r.db.Query(context.Background(), `
        SELECT id, opened, closed, author_id, assigned_id, title, content
        FROM tasks WHERE ($1 = 0 OR id = $1) AND ($2 = 0 OR author_id = $2)
        ORDER BY id;
    `, taskID, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []core.Task
	for rows.Next() {
		var task core.Task
		if err := rows.Scan(&task.ID, &task.Opened, &task.Closed, &task.AuthorID, &task.AssignedID, &task.Title, &task.Content); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func (r *postgresRepository) CreateTask(task core.Task) (int, error) {
	var id int
	err := r.db.QueryRow(context.Background(), `
        INSERT INTO tasks (opened, closed, author_id, assigned_id, title, content)
        VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;
    `, task.Opened, task.Closed, task.AuthorID, task.AssignedID, task.Title, task.Content).Scan(&id)
	return id, err
}

func (r *postgresRepository) CreateTasks(tasks []core.Task) ([]int, error) {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	var ids []int
	for _, task := range tasks {
		var id int
		err := tx.QueryRow(context.Background(), `
            INSERT INTO tasks (opened, closed, author_id, assigned_id, title, content)
            VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;
        `, task.Opened, task.Closed, task.AuthorID, task.AssignedID, task.Title, task.Content).Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, tx.Commit(context.Background())
}

func (r *postgresRepository) UpdateTask(task core.Task) error {
	_, err := r.db.Exec(context.Background(), `
        UPDATE tasks SET opened = $1, closed = $2, author_id = $3, assigned_id = $4, title = $5, content = $6
        WHERE id = $7;
    `, task.Opened, task.Closed, task.AuthorID, task.AssignedID, task.Title, task.Content, task.ID)
	return err
}

func (r *postgresRepository) DeleteTask(taskID int) error {
	_, err := r.db.Exec(context.Background(), `DELETE FROM tasks WHERE id = $1;`, taskID)
	return err
}
