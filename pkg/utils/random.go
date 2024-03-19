// Copyright (c) Microsoft Corporation. All rights reserved.

package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().Unix()))
}

func GenerateRandomString(length int) string {
	// avoid symbols which are not shell-friendly and some other things that aren't friendly
	const charSet = "abcdefghjkmnopqrstuvwxyzABCDEFGHJKMNOPQRSTUVWXYZ234567890:@%^"
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(b)
}

func GenerateRandomStringFromCharSet(length int, charSet string) string {
	// avoid symbols which are not shell-friendly and some other things that aren't friendly
	if len(charSet) == 0 {
		return ""
	}
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(b)
}
