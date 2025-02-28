package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/erknas/todo-list/internal/config"
	"github.com/erknas/todo-list/internal/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPool struct {
	pool *pgxpool.Pool
}

func NewPostgresPool(ctx context.Context, cfg *config.Config) (*PostgresPool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName)

	poolCfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &PostgresPool{pool: pool}, nil
}

func (p *PostgresPool) CreateTask(ctx context.Context, task types.TaskRequest) (int, error) {
	query := `INSERT INTO tasks(title, description, status) VALUES($1, $2, $3) RETURNING id`

	var id int

	if err := p.pool.QueryRow(ctx, query, task.Title, task.Description, task.Status).Scan(&id); err != nil {
		return -1, err
	}

	return id, nil
}

func (p *PostgresPool) GetTasks(ctx context.Context) ([]types.Task, error) {
	query := `SELECT id, title, description, status, created_at, updated_at FROM tasks ORDER BY created_at`

	var tasks []types.Task

	rows, err := p.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := types.Task{}
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (p *PostgresPool) UpdateTask(ctx context.Context, id int, req types.TaskRequest) error {
	query, args := prepareUpdate(id, req)

	res, err := p.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return &TaskNotFoundError{ID: id}
	}

	return nil
}

func (p *PostgresPool) DeleteTask(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE id=$1`

	res, err := p.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return &TaskNotFoundError{ID: id}
	}

	return nil
}

func (p *PostgresPool) Close() {
	p.pool.Close()
}

func prepareUpdate(id int, req types.TaskRequest) (string, []any) {
	var (
		args   []any
		fields []string
		count  = 1
	)

	query := `UPDATE tasks SET `

	if len(req.Title) != 0 {
		fields = append(fields, fmt.Sprintf("title=$%d", count))
		args = append(args, req.Title)
		count++
	}

	if len(req.Description) != 0 {
		fields = append(fields, fmt.Sprintf("description=$%d", count))
		args = append(args, req.Description)
		count++
	}

	if len(req.Status) != 0 {
		fields = append(fields, fmt.Sprintf("status=$%d", count))
		args = append(args, req.Status)
		count++
	}

	fields = append(fields, fmt.Sprintf("updated_at=$%d", count))
	args = append(args, time.Now())
	count++

	query += fmt.Sprintf("%s WHERE id=$%d", join(fields, ", "), count)
	args = append(args, id)

	return query, args
}

func join(fields []string, sep string) string {
	var result string

	for i, element := range fields {
		if i > 0 {
			result += sep
		}
		result += element
	}

	return result
}
