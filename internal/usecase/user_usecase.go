package usecase

import (
    "context"

    "midtrans-forwarder/internal/domain"
    "midtrans-forwarder/internal/repository"
    "golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
    userRepo *repository.UserRepository
}

func NewUserUseCase(userRepo *repository.UserRepository) *UserUseCase {
    return &UserUseCase{userRepo: userRepo}
}

func (u *UserUseCase) CreateUser(ctx context.Context, user *domain.User) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.Password = string(hashedPassword)
    return u.userRepo.CreateUser(ctx, user)
}

func (u *UserUseCase) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
    return u.userRepo.GetUserByID(ctx, id)
}

func (u *UserUseCase) GetAllUsers(ctx context.Context) ([]domain.User, error) {
    return u.userRepo.GetAllUsers(ctx)
}

func (u *UserUseCase) UpdateUser(ctx context.Context, user *domain.User) error {
    return u.userRepo.UpdateUser(ctx, user)
}

func (u *UserUseCase) DeleteUser(ctx context.Context, id int64) error {
    return u.userRepo.DeleteUser(ctx, id)
}