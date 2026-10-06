package user

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(u User) (int64, error) {
	query := `
		INSERT INTO users (
			username,
			password_hash,
			full_name,
			email,
			role,
			active,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		u.Username,
		u.PasswordHash,
		u.FullName,
		u.Email,
		u.Role,
		u.Active,
		u.CreatedAt,
		u.UpdatedAt,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) GetByID(id int) (User, error) {
	query := `
		SELECT
			id,
			username,
			password_hash,
			full_name,
			email,
			role,
			active,
			created_at,
			updated_at
		FROM users
		WHERE id = ?
	`

	var u User
	if err := r.db.QueryRow(query, id).Scan(
		&u.ID,
		&u.Username,
		&u.PasswordHash,
		&u.FullName,
		&u.Email,
		&u.Role,
		&u.Active,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return User{}, err
	}

	return u, nil
}

func (r *Repository) GetByUsername(username string) (User, error) {
	query := `
		SELECT
			id,
			username,
			password_hash,
			full_name,
			email,
			role,
			active,
			created_at,
			updated_at
		FROM users
		WHERE username = ?
	`

	var u User
	if err := r.db.QueryRow(query, username).Scan(
		&u.ID,
		&u.Username,
		&u.PasswordHash,
		&u.FullName,
		&u.Email,
		&u.Role,
		&u.Active,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return User{}, err
	}

	return u, nil
}

func (r *Repository) List() ([]User, error) {
	query := `
		SELECT
			id,
			username,
			password_hash,
			full_name,
			email,
			role,
			active,
			created_at,
			updated_at
		FROM users
		ORDER BY id ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.PasswordHash,
			&u.FullName,
			&u.Email,
			&u.Role,
			&u.Active,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) Update(u User) error {
	query := `
		UPDATE users
		SET
			username = ?,
			password_hash = ?,
			full_name = ?,
			email = ?,
			role = ?,
			active = ?,
			updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		u.Username,
		u.PasswordHash,
		u.FullName,
		u.Email,
		u.Role,
		u.Active,
		u.UpdatedAt,
		u.ID,
	)
	return err
}
