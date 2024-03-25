package users

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	mrand "math/rand"
	"strings"
)

const (
	lowerCharSet = "abcdefghijklmnopqrstuvwxyz"
	upperCharSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberSet    = "0123456789"
)

func GeneratePassword(passwordRequirements *PasswordRequirements) (string, error) {
	return generatePassword(rand.Reader, passwordRequirements)
}

func GenerateEncryptionKey(requiremements *EncryptionKeyRequirements) ([]byte, error) {
	if requiremements == nil || requiremements.KeySize == nil {
		return []byte{}, fmt.Errorf("password minimumLength cannot be nil")
	}

	result := make([]byte, *requiremements.KeySize)
	_, err := rand.Read(result)
	if err != nil {
		return nil, fmt.Errorf("failed to generate encryption key: %v", err)
	}
	return result, nil
}

// GeneratePassword generates a random password with the given length and minimum number of special characters, numbers, and uppercase letters.
// If any of the minimums are zero, then the password will not contain any of those characters.  All minimums will be met even if the password
// length must be increased to do so.
func generatePassword(reader io.Reader, passwordRequirements *PasswordRequirements) (string, error) {
	var allCharSet, password strings.Builder
	var minSpecialChar, minNum, minUpperCase int32

	if passwordRequirements == nil || passwordRequirements.MinLength == nil {
		return "", fmt.Errorf("password minimumLength cannot be nil")
	}

	specs := make(map[string]int32)

	// Set special character
	if passwordRequirements.MinSpecialChar != nil && *passwordRequirements.MinSpecialChar > 0 {
		minSpecialChar = *passwordRequirements.MinSpecialChar
		if passwordRequirements.SpecialCharList == nil || len(*passwordRequirements.SpecialCharList) == 0 {
			return "", fmt.Errorf("special character list cannot be empty when minSpecialChar is greater than 0")
		}
		allCharSet.WriteString(*passwordRequirements.SpecialCharList)
		specs[*passwordRequirements.SpecialCharList] = minSpecialChar
	}

	// Set numeric
	if passwordRequirements.MinNumeric != nil && *passwordRequirements.MinNumeric > 0 {
		minNum = *passwordRequirements.MinNumeric
		allCharSet.WriteString(numberSet)
		specs[numberSet] = minNum
	}

	// Set uppercase
	if passwordRequirements.MinUpperCase != nil && *passwordRequirements.MinUpperCase > 0 {
		minUpperCase = *passwordRequirements.MinUpperCase
		allCharSet.WriteString(upperCharSet)
		specs[upperCharSet] = minUpperCase
	}

	// Always add at least one lower case
	minLowerCase := int32(1)
	allCharSet.WriteString(lowerCharSet)
	specs[lowerCharSet] = minLowerCase

	// Calculate remaining length
	remainingLength := *passwordRequirements.MinLength - minSpecialChar - minNum - minUpperCase - minLowerCase
	if remainingLength > 0 {
		specs[allCharSet.String()] = remainingLength
	}

	// Build password
	for k, v := range specs {
		for i := int32(0); i < v; i++ {
			random, err := rand.Int(reader, big.NewInt(int64(len(k))))
			if err != nil {
				return "", err
			}
			password.WriteString(string(k[random.Int64()]))
		}
	}

	// Shuffle password
	inRune := []rune(password.String())
	mrand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune), nil
}
