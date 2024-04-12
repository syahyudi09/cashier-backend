package usecase

import (
	"cashier/formatter"
	"cashier/model"
	"cashier/repository"
	"cashier/utils"
	"cashier/utils/token"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidUUIDFormat = errors.New("INVALID_ID_FORMAT")

type UserUsecase interface {
	RegisterUser(register model.RegisterUserInput) error
	LoginUser(input model.LoginUserInput) (model.UserFormatter, error)
	EmailExits(email string) (bool, error)
	GetAllUser(ctx context.Context, page, pageSize int, search string) ([]*formatter.UserFormatter, int, error)
	UpdateUser(userID string, updatedUser model.UpdateUserInput) error
	CheckEmailForUpdate(email string, userID string) (bool, error)
	FindById(userID string) (model.UserModel, error)
	Delete(userID string) error
}

type userUsecase struct {
	repository repository.UserRepository
}

func (u *userUsecase) RegisterUser(register model.RegisterUserInput) error {
	user := model.UserModel{}
	user.Id = utils.UuidGenerate()
	user.Fullname = register.Fullname
	user.Email = register.Email
	user.Role = register.Role
	user.Status = "Active"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("err %w", err)
	}
	user.Password = string(hashedPassword)
	return u.repository.Create(user)
}

func (u *userUsecase) LoginUser(input model.LoginUserInput) (model.UserFormatter, error) {
	email := input.Email
	password := input.Password

	user, err := u.repository.FindByEmail(email)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return model.UserFormatter{}, fmt.Errorf("invalid password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println(err)
		return model.UserFormatter{}, fmt.Errorf("invalid password")
	}

	accessToken, err := token.GenerateToken(user.Id, string(user.Role))
	if err != nil {
		return model.UserFormatter{}, fmt.Errorf("failed to generate token: %w", err)
	}

	refreshToken, err := token.GenerateRefreshToken(user.Id, string(user.Role))
	if err != nil {
		return model.UserFormatter{}, fmt.Errorf("failed to generate token: %w", err)
	}

	formatter := model.UserFormatter{
		ID:           string(user.Id),
		Fullname:     user.Fullname,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserRole:     user.Role,
	}

	return formatter, nil
}

func (u *userUsecase) EmailExits(Email string) (bool, error) {
	return u.repository.CheckEmail(Email)
}

func (u *userUsecase) GetAllUser(ctx context.Context, page, pageSize int, search string) ([]*formatter.UserFormatter, int, error) {
	users, totalDocs, err := u.repository.GetAll(ctx, page, pageSize, search)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching users: %v", err)
	}

	var userFormatters []*formatter.UserFormatter
	for _, user := range users {
		userFormatters = append(userFormatters, &formatter.UserFormatter{
			Id:       user.Id,
			Fullname: user.Fullname,
			Email:    user.Email,
			Role:     string(user.Role),
			Status:   string(user.Status),
		})
	}

	return userFormatters, totalDocs, nil
}

func (u *userUsecase) UpdateUser(userID string, updatedUser model.UpdateUserInput) error {
	user, err := u.repository.FindByID(userID)
	if err != nil {
		return fmt.Errorf("error fetching user: %w", err)
	}

	if _, err := uuid.Parse(userID); err != nil {
		return ErrInvalidUUIDFormat
	}

	if updatedUser.Fullname != "" {
		user.Fullname = updatedUser.Fullname
	}

	if updatedUser.Email != "" {
		user.Email = updatedUser.Email
	}

	if updatedUser.Role != "" {
		user.Role = updatedUser.Role
	}
	if updatedUser.Status != "" {
		user.Status = updatedUser.Status
	}

	user.UpdatedAt = time.Now()

	if updatedUser.Password != "" {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return fmt.Errorf("err %w", err)
		}
		user.Password = string(hashedPassword)
	}

	// Panggil repository untuk melakukan update
	if err := u.repository.Update(userID, user); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

func (u *userUsecase) CheckEmailForUpdate(email string, userID string) (bool, error) {
	exists, err := u.repository.CheckEmailForUpdate(email, userID)
	if err != nil {
		return false, fmt.Errorf("error checking email for update: %w", err)
	}
	return exists, nil
}

func (u *userUsecase) FindById(userID string) (model.UserModel, error) {
	return u.repository.FindByID(userID)
}

func (u *userUsecase) Delete(userID string) error {
	if _, err := uuid.Parse(userID); err != nil {
		return ErrInvalidUUIDFormat
	}
	return u.repository.Delete(userID)
}

func NewUserUsecase(r repository.UserRepository) UserUsecase {
	return &userUsecase{
		repository: r,
	}
}
