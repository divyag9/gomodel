package database

import (
	"github.com/divyag9/gomodel/pkg/pb"
	ora "gopkg.in/rana/ora.v4"
)

//Config contains information required for the database operations
type Config struct {
	Session *ora.Ses
}

//NewConfig initializes the Config struct and returns it
func NewConfig(session *ora.Ses) *Config {
	return &Config{
		Session: session,
	}
}

//GetImageDetailsByOrderNumber retrieves the imageDetails for a given orderNumber
func (c *Config) GetImageDetailsByOrderNumber(orderNumber int64) ([]*contentservice.ImageDetail, error) {
	return []*contentservice.ImageDetail{&contentservice.ImageDetail{}}, nil
}

//GetImageDetailsByImageIds retrieves the imageDetails for a given set og imageIds
func (c *Config) GetImageDetailsByImageIds(imageIds []int64) ([]*contentservice.ImageDetail, error) {
	return []*contentservice.ImageDetail{&contentservice.ImageDetail{}}, nil
}
