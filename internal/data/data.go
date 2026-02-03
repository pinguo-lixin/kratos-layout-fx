package data

import (
	"github.com/go-kratos/kratos-layout/internal/conf"
)

// Data .
type Data struct {
	// TODO wrapped database client
}

// NewData .
func NewData(c *conf.Data) (*Data, error) {
	return &Data{}, nil
}
