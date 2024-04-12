package utils

import "github.com/google/uuid"

func UuidGenerate() string {
    id := uuid.New()
    return id.String()
}
