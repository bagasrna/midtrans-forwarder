package usecase

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"

	"midtrans-forwarder/internal/domain"
	"midtrans-forwarder/internal/repository"
)

type MidtransUseCase struct {
	resellerRepo *repository.ResellerRepository
	serverKey    string
}

func NewMidtransUseCase(resellerRepo *repository.ResellerRepository, serverKey string) *MidtransUseCase {
	return &MidtransUseCase{
		resellerRepo: resellerRepo,
		serverKey:    serverKey,
	}
}

func (u *MidtransUseCase) ValidateCallback(callback domain.MidtransCallback) error {
	signature := fmt.Sprintf("%s%s%s%s", callback.OrderID, callback.TransactionStatus, callback.GrossAmount, u.serverKey)
	hash := sha512.Sum512([]byte(signature))
	if hex.EncodeToString(hash[:]) != callback.SignatureKey {
		return fmt.Errorf("invalid signature")
	}
	return nil
}

func (u *MidtransUseCase) ForwardToReseller(ctx context.Context, callback domain.MidtransCallback) error {
	resellers, err := u.resellerRepo.GetAllResellers(ctx)
	if err != nil {
		return err
	}

	resellerCode := strings.Split(callback.OrderID, "-")[0]
	for _, reseller := range resellers {
		if reseller.Code == resellerCode {
			fmt.Printf("Forwarding to reseller %s at URL %s\n", reseller.Name, reseller.URL)
			return nil
		}
	}

	return fmt.Errorf("no matching reseller found for code %s", resellerCode)
}