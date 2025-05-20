package utils

import (
    "strings"
    "crypto/sha256"
    "encoding/hex"
) 

func NormalizeFEN(fen string) string {
    parts := strings.Split(fen, " ")
    if len(parts) < 4 {
        return fen
    }
    return strings.Join(parts[:4], " ")
}

func HashFEN(normalizedFEN string) string {
    hash := sha256.Sum256([]byte(normalizedFEN))
    return hex.EncodeToString(hash[:])
}
