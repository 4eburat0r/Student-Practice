package hash

import "golang.org/x/crypto/bcrypt"

type Hasher interface {
    Hash(password string) (string, error)
    Verify(hashedPassword, password string) error
}

type BcryptHasher struct {
    cost int
}

func NewBcryptHasher() *BcryptHasher {
    return &BcryptHasher{
        cost: bcrypt.DefaultCost,
    }
}

func (h *BcryptHasher) Hash(password string) (string, error) {
    hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
    if err != nil {
        return "", err
    }
    return string(hashedBytes), nil
}

func (h *BcryptHasher) Verify(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
