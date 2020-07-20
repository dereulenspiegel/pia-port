package main

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandString(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

const numberLetterBytes = "1234567890"
const (
	numberLetterIdxBits = 4
	numberLetterIdxMask = 1<<numberLetterIdxBits - 1
	numberLetterIdxMax  = 63 / numberLetterIdxBits
)

func RandNumberString(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), numberLetterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), numberLetterIdxMax
		}
		if idx := int(cache & numberLetterIdxMask); idx < len(numberLetterBytes) {
			b[i] = numberLetterBytes[idx]
			i--
		}
		cache >>= numberLetterIdxBits
		remain--
	}

	return string(b)
}
