package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/divyag9/gomodel/pkg/interface"
	"github.com/divyag9/gomodel/pkg/pb"
)

// Config holds all the information required for caching
type Config struct {
	Memcache          *memcache.Client
	SecondsToExpiry   int32
	OrderNumberGetter contentserviceinterface.OrderNumberGetter
}

//NewOrderClient initializes the cache client to get imagedetails for ordernumber
func NewOrderClient(memcache *memcache.Client, secondsToExpiry int32, orderNumberGetter contentserviceinterface.OrderNumberGetter) *Config {
	return &Config{
		Memcache:          memcache,
		SecondsToExpiry:   secondsToExpiry,
		OrderNumberGetter: orderNumberGetter,
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

//setImageDetails sets the ImageDetails to cache for a given key
func setImageDetails(key string, imageDetails []*contentservice.ImageDetail) error {
	return nil
}
