package usecase

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"midtrans-forwarder/internal/domain"
)

type CallbackMidtransUseCase struct {
	serverKey string
}

func NewCallbackMidtransUseCase(serverKey string) *CallbackMidtransUseCase {
	return &CallbackMidtransUseCase{serverKey: serverKey}
}

func (u *CallbackMidtransUseCase) ValidateCallback(callback domain.MidtransCallback) error {
	signature := fmt.Sprintf("%s%s%s%s", callback.OrderID, callback.TransactionStatus, callback.GrossAmount, u.serverKey)
	hash := sha512.Sum512([]byte(signature))
	if hex.EncodeToString(hash[:]) != callback.SignatureKey {
		return fmt.Errorf("invalid signature")
	}
	return nil
}
