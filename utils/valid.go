package utils

import (
	"context"
	"fashora-backend/services/user_service"
	"regexp"
	"strings"
)

func ValidatePhoneNumber(phone string) bool {
	re := regexp.MustCompile(`^\d{10,11}$`)
	return re.MatchString(phone)
}

func formatPhoneNumberVN(phone string) string {
	if strings.HasPrefix(phone, "0") {
		return "+84" + phone[1:]
	}
	return phone
}

func ValidatePhoneOTP(ctx context.Context, phone string) bool {
	verifiedPhones, _ := user_service.GetVerifiedPhoneNumbers(ctx)
	verifiedPhoneMap := make(map[string]bool, len(verifiedPhones))
	for _, verifiedPhone := range verifiedPhones {
		verifiedPhoneMap[verifiedPhone] = true
	}
	return verifiedPhoneMap[formatPhoneNumberVN(phone)]
}
