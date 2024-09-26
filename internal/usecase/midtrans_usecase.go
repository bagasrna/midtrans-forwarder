package usecase

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

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
	signature := fmt.Sprintf("%s%s%s%s", callback.OrderID, callback.StatusCode, callback.GrossAmount, u.serverKey)
	hash := sha512.Sum512([]byte(signature))
	if hex.EncodeToString(hash[:]) != callback.SignatureKey {
		return fmt.Errorf("invalid signature")
	}
	return nil
}

func (u *MidtransUseCase) ForwardToReseller(ctx context.Context, callback domain.MidtransCallback, rawBody []byte) error {
	// Raps Forwarder
	orderParts := strings.Split(callback.OrderID, "-")
	if len(orderParts) == 2 {
		// Jika hanya ada satu tanda "-", forward ke URL dari environment variable RAPS_URL
		rapsURL := os.Getenv("RAPS_URL") + "/api/midtrans/callback"
		if rapsURL == "" {
			return fmt.Errorf("RAPS_URL is not set in environment variables")
		}
		fmt.Printf("Forwarding to RAPS_URL at %s\n", rapsURL)
		// Implementasikan logika untuk meneruskan ke rapsURL di sini
		return u.forwardRequest(ctx, rapsURL, rawBody)
	}

	// Reseller Forwarder
	resellers, err := u.resellerRepo.GetAllResellers(ctx)
	if err != nil {
		return err
	}

	resellerCode := strings.Split(callback.OrderID, "-")[0]
	for _, reseller := range resellers {
		if reseller.Code == resellerCode {
			resellerUrl := reseller.URL + "/api/midtrans/callback"
			return u.forwardRequest(ctx, resellerUrl, rawBody)
		}
	}

	return fmt.Errorf("no matching reseller found for code %s", resellerCode)
}

func (u *MidtransUseCase) forwardRequest(ctx context.Context, url string, rawBody[]byte) error {
	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
		"User-Agent":   "Mozilla/4.0",
	}

	for attempts := 0; attempts < 3; attempts++ {
		fmt.Println(bytes.NewBuffer(rawBody))
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(rawBody))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		for key, value := range headers {
			req.Header.Set(key, value)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Failed to forward request, attempt %d: %s\n", attempts+1, err)
			time.Sleep(time.Duration(1<<attempts) * time.Second) // Exponential backoff
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Successful response
			fmt.Printf("Successfully forwarded request to %s\n", url)
			return nil
		}

		fmt.Printf("Failed to forward request, attempt %d, response status: %d\n", attempts+1, resp.StatusCode)
		time.Sleep(time.Duration(1<<attempts) * time.Second) // Exponential backoff
	}

	return fmt.Errorf("failed to forward request after 3 attempts: " + string(rawBody))
}