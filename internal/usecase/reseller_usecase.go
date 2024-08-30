package usecase

import (
    "context"

    "midtrans-forwarder/internal/domain"
    "midtrans-forwarder/internal/repository"
)

type ResellerUseCase struct {
    resellerRepo *repository.ResellerRepository
}

func NewResellerUseCase(resellerRepo *repository.ResellerRepository) *ResellerUseCase {
    return &ResellerUseCase{resellerRepo: resellerRepo}
}

func (u *ResellerUseCase) CreateReseller(ctx context.Context, reseller *domain.Reseller) error {
    
    return u.resellerRepo.CreateReseller(ctx, reseller)
}

func (u *ResellerUseCase) GetResellerByID(ctx context.Context, id int64) (*domain.Reseller, error) {
    return u.resellerRepo.GetResellerByID(ctx, id)
}

func (u *ResellerUseCase) GetAllResellers(ctx context.Context) ([]domain.Reseller, error) {
    return u.resellerRepo.GetAllResellers(ctx)
}

func (u *ResellerUseCase) UpdateReseller(ctx context.Context, reseller *domain.Reseller) error {
    return u.resellerRepo.UpdateReseller(ctx, reseller)
}

func (u *ResellerUseCase) DeleteReseller(ctx context.Context, id int64) error {
    return u.resellerRepo.DeleteReseller(ctx, id)
}