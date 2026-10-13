package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/Lykeion-org/demes/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)


type UserRepository interface {
	CreateUser(ctx context.Context, email string, password string) (*models.LykeionUser, error)
	ReadUser(ctx context.Context, uid string) (*models.LykeionUser, error)
	ReadUserByEmailAddress(ctx context.Context, email string) (*models.LykeionUser, error)
	UpdateUser(ctx context.Context, user *models.LykeionUser) (*models.LykeionUser, error)
	DeleteUser(ctx context.Context, uid string) (*models.LykeionUser, error)
	ListUsers(ctx context.Context, page int, pageSize int) ([]*models.LykeionUser,error)
	CheckPassword(ctx context.Context, email string, password string) (bool, error)
}

type userRepository struct {	
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func(r *userRepository)CreateUser(ctx context.Context, email string, password string) (*models.LykeionUser, error){
	email = normalizeEmail(email)
	if email == "" {
		return nil, status.Error(codes.InvalidArgument, "email address is required")
	}

	_, err := r.ReadUserByEmailAddress(ctx, email)
	if err == nil {
		return nil, errEmailTaken
	}
	if !errors.Is(err, errUserNotFound) {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	newUser := &models.LykeionUser{
		UID: 					uuid.New(),
		EmailAddress: 			email,
		PasswordHash: 			string(passwordHash),
		EmailAddressHash: 		hashEmail(email),
		AvatarSource: 			"/users/avatar/placeholder.webp",
		NativeLanguage: 		0,
		ReferenceLanguage: 		0,
		TargetLanguage: 		0,
		Role: 					1,
		Verified: 				false,
		CreatedAt: 				time.Now(),
		UpdatedAt: 				time.Now(),
	}
	err = r.db.
		WithContext(ctx).
		Create(newUser).
		Error
	if isUniqueViolation(err) { return nil, errEmailTaken }
	if err != nil { return nil, errors.New("error creating user")}
	
	return newUser, nil

}
func(r *userRepository)ReadUser(ctx context.Context, uid string) (*models.LykeionUser, error){
	userUUID, err := uuid.Parse(uid)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	var user models.LykeionUser
	err = r.db.
		WithContext(ctx).
		Where("uid = ?", userUUID).
		First(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("error fetching user")
	}

	return &user, nil
}

func(r *userRepository)ReadUserByEmailAddress(ctx context.Context, email string) (*models.LykeionUser, error){
	var user models.LykeionUser
	
	err := r.db.
		WithContext(ctx).
		Where("lower(email_address) = ?", normalizeEmail(email)).
		First(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errUserNotFound
		}
		return nil, errors.New("error fetching user")
	}

	return &user, nil
}

func(r *userRepository)UpdateUser(ctx context.Context, user *models.LykeionUser) (*models.LykeionUser, error){
	if user.UID == uuid.Nil {
		return nil, errors.New("user ID cannot be empty")
	}

	err := r.db.
		WithContext(ctx).
		Save(user).
		Error

	if isUniqueViolation(err) {
		return nil, errEmailTaken
	}
	if err != nil {
		return nil, errors.New("error updating user")
	}

	return user, nil
}
func(r *userRepository)DeleteUser(ctx context.Context, uid string) (*models.LykeionUser, error){
	userUUID, err := uuid.Parse(uid)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	var user models.LykeionUser
	err = r.db.
		WithContext(ctx).
		Where("uid = ?", userUUID).
		First(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("error finding user to delete")
	}

	err = r.db.
		WithContext(ctx).
		Delete(&user).
		Error

	if err != nil {
		return nil, errors.New("error deleting user")
	}

	return &user, nil
}
func(r *userRepository)ListUsers(ctx context.Context, page int, pageSize int) ([]*models.LykeionUser,error){
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var users []*models.LykeionUser
	offset := (page - 1) * pageSize

	err := r.db.
		WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Find(&users).
		Error

	if err != nil {
		return nil, errors.New("error listing users")
	}

	return users, nil
}

func(r *userRepository)CheckPassword(ctx context.Context, email string, password string) (bool, error){
	user, err := r.ReadUserByEmailAddress(ctx, email)
	if err != nil {
		return false, err
	}

	slog.Debug("user found in database","user", user)

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return false, errors.New("error checking password")
	}

	return true, nil

}

var (
	errUserNotFound = errors.New("user not found")
	errEmailTaken   = status.Error(codes.AlreadyExists, "email address is already registered")
)

// Emails are compared case-insensitively, so "A@x.com" and "a@x.com" are the same account.
func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func hashEmail(email string) string {
	hash := sha256.Sum256([]byte(email))
	return hex.EncodeToString(hash[:])
}
