package repository

import (
	"cashier/model"
	"context"
	"database/sql"
	"fmt"
)

type UserRepository interface {
	Create(newUser model.UserModel) error
	FindByEmail(Email string) (model.UserModel, error)
	CheckEmail(email string) (bool, error)
	GetAll(ctx context.Context, page, pageSize int, search string) ([]*model.UserModel, int, error)
	Update(ID string, updatedUser model.UserModel) error
	FindByID(ID string) (model.UserModel, error)
	CheckEmailForUpdate(email, userID string) (bool, error)
	Delete(ID string) error 
}

type userRepository struct {
	db *sql.DB
}

func (u *userRepository) Create(newUser model.UserModel) error {
	insertQuery := "INSERT INTO users(id, fullname, email, password, role, status, created_at, updated_at) VALUES($1,$2,$3,$4,$5,$6,$7, $8)"

	_, err := u.db.Exec(insertQuery, newUser.Id, newUser.Fullname, newUser.Email, newUser.Password, newUser.Role, newUser.Status, newUser.CreatedAt, newUser.UpdatedAt)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("err on userRepository.Register: %v", err)
	}
	return nil
}

func (u *userRepository) FindByEmail(Email string) (model.UserModel, error) {
	getQuery := "SELECT id, fullname, email, password ,role from users WHERE email = $1"

	row := u.db.QueryRow(getQuery, Email)

	var user model.UserModel
	err := row.Scan(&user.Id, &user.Fullname, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return user, fmt.Errorf("error fetching user by email: %w", err)
	}
	return user, nil
}

func (u *userRepository) CheckEmail(email string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)"

	var exists bool
	err := u.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error fetching user by email: %w", err)
	}

	return exists, nil
}

func (u *userRepository) FindByID(ID string) (model.UserModel, error) {
	getQuery := "SELECT id, fullname, email, password ,role, status from users WHERE id = $1"

	row := u.db.QueryRow(getQuery, ID)

	var user model.UserModel
	err := row.Scan(&user.Id, &user.Fullname, &user.Email, &user.Password, &user.Role, &user.Status)
	if err != nil {
		return user, fmt.Errorf("error fetching user by Id: %w", err)
	}
	return user, nil
}

func (u *userRepository) GetAll(ctx context.Context, page, pageSize int, search string) ([]*model.UserModel, int, error) {
	offset := (page - 1) * pageSize
	getQuery := `
		SELECT id, fullname, email, password, role, status, created_at, updated_at 
		FROM users 
		WHERE fullname LIKE '%' || $1 || '%'
		ORDER BY id 
		LIMIT $2 
		OFFSET $3`

	rows, err := u.db.QueryContext(ctx, getQuery, search, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error on campaignRepoImpl.FindAll: %w", err)
	}
	defer rows.Close()

	var arrUser []*model.UserModel

	for rows.Next() {
		user := &model.UserModel{}
		if err := rows.Scan(
			&user.Id, &user.Fullname, &user.Email, &user.Password, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("error scanning campaign row: %w", err)
		}
		arrUser = append(arrUser, user)
	}

	countQuery := "SELECT COUNT(id) FROM users" // Menyesuaikan query COUNT dengan offset
	var totalDocs int
	if err := u.db.QueryRowContext(ctx, countQuery).Scan(&totalDocs); err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %w", err)
	}

	return arrUser, totalDocs, nil
}


func (u *userRepository) Update(userID string, updatedUser model.UserModel) error {
	updateQuery := "UPDATE users SET fullname = $1, email = $2, password = $3, role = $4, status = $5, updated_at = $6 WHERE id = $7"

	_, err := u.db.Exec(updateQuery, updatedUser.Fullname, updatedUser.Email, updatedUser.Password, updatedUser.Role, updatedUser.Status, updatedUser.UpdatedAt, userID)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("error on UserRepository.Update: %v", err)
	}
	return nil
}

func (u *userRepository) CheckEmailForUpdate(email, userID string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1 AND id != $2)"

	var exists bool
	err := u.db.QueryRow(query, email, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking existing email for update: %w", err)
	}

	return exists, nil
}

func (u *userRepository) Delete(ID string) error {
	deleteQuery := "DELETE FROM users WHERE id = $1"

	_, err := u.db.Query(deleteQuery, ID)
	if err != nil {
		return fmt.Errorf("error fetching user by Id: %w", err)
	}
	return nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
