package usecase

import (
	"context"
	"fmt"
	"strings"

	"midtrans-forwarder/internal/domain"
	"midtrans-forwarder/internal/repository"
)

type ForwardMidtransUseCase struct {
	resellerRepo *repository.ResellerRepository
}

func NewForwardMidtransUseCase(resellerRepo *repository.ResellerRepository) *ForwardMidtransUseCase {
	return &ForwardMidtransUseCase{resellerRepo: resellerRepo}
}

func (u *ForwardMidtransUseCase) ForwardToReseller(ctx context.Context, callback domain.MidtransCallback) error {
	resellers, err := u.resellerRepo.GetAllResellers(ctx)
	if err != nil {
		return err
	}

	resellerCode := strings.Split(callback.OrderID, "-")[0]
	for _, reseller := range resellers {
		if reseller.Code == resellerCode {
			// Forward the callback to the reseller's URL
			// Note: Implement the actual HTTP forwarding logic here
			fmt.Printf("Forwarding to reseller %s at URL %s\n", reseller.Name, reseller.URL)
			return nil
		}
	}

	return fmt.Errorf("no matching reseller found for code %s", resellerCode)
}