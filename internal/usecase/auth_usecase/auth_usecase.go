package auth_usecase

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"gin-real-time-talk/internal/entity"
	"gin-real-time-talk/internal/entity/interfaces"
	"gin-real-time-talk/pkg/email"
	"gin-real-time-talk/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepo     interfaces.UserRepository
	emailService *email.EmailService
}

func NewAuthUsecase(userRepo interfaces.UserRepository, emailService *email.EmailService) interfaces.AuthUsecase {
	return &authUsecase{
		userRepo:     userRepo,
		emailService: emailService,
	}
}

func (u *authUsecase) Register(email, password, firstName, lastName string) (*entity.User, error) {
	existingUser, err := u.userRepo.GetByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &entity.User{
		Email:         email,
		Password:      string(hashedPassword),
		FirstName:     firstName,
		LastName:      lastName,
		EmailVerified: false,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if u.emailService.IsConfigured() {
		if err := u.SendTwoFactorCode(email); err != nil {
			return nil, fmt.Errorf("failed to send verification code: %w", err)
		}
	}

	return user, nil
}

func (u *authUsecase) Login(email, password string) (string, string, *entity.User, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return "", "", nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", nil, errors.New("invalid email or password")
	}

	if !user.EmailVerified {
		return "", "", nil, errors.New("email not verified")
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, user, nil
}

func (u *authUsecase) SendTwoFactorCode(email string) error {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	code := generateCode()
	expiresAt := time.Now().Add(10 * time.Minute)

	user.TwoFactorCode = code
	user.TwoFactorExpiresAt = &expiresAt

	if err := u.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if u.emailService.IsConfigured() {
		if err := u.emailService.SendVerificationCode(email, code); err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
	} else {
		return fmt.Errorf("email service not configured, code: %s", code)
	}

	return nil
}

func (u *authUsecase) VerifyTwoFactorCode(email, code string) (string, string, *entity.User, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return "", "", nil, errors.New("user not found")
	}

	if user.TwoFactorCode == "" || user.TwoFactorExpiresAt == nil {
		return "", "", nil, errors.New("verification code not found or expired")
	}

	if time.Now().After(*user.TwoFactorExpiresAt) {
		return "", "", nil, errors.New("verification code expired")
	}

	if user.TwoFactorCode != code {
		return "", "", nil, errors.New("invalid verification code")
	}

	user.EmailVerified = true
	user.TwoFactorCode = ""
	user.TwoFactorExpiresAt = nil

	if err := u.userRepo.Update(user); err != nil {
		return "", "", nil, fmt.Errorf("failed to update user: %w", err)
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, user, nil
}

func (u *authUsecase) RefreshToken(refreshToken string) (string, string, *entity.User, error) {
	claims, err := jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", nil, errors.New("invalid refresh token")
	}

	user, err := u.userRepo.GetByID(claims.UserID)
	if err != nil {
		return "", "", nil, errors.New("user not found")
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := jwt.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, newRefreshToken, user, nil
}

func (u *authUsecase) ValidateAccessToken(token string) (*entity.User, error) {
	claims, err := jwt.ValidateAccessToken(token)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	user, err := u.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func generateCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("%06d", n.Int64())
}
