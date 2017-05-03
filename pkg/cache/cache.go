package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/divyag9/gomodel/pkg/pb/github.com/divyag9/proto"
)

// Info holds all the information required for caching
type Info struct {
	MemClient       *memcache.Client
	SecondsToExpiry int32
}

// Caller contains methods for caching
type Caller interface {
	GetImageDetails(string) ([]*contentservice.ImageDetail, error)
	SetImageDetails(string, []*contentservice.ImageDetail) error
}

//GetImageDetails retrieves ImageDetails from cache for a given key
func (i *Info) GetImageDetails(key string /*may be also send the other contexts*/) ([]*contentservice.ImageDetail, error) {
	return []*contentservice.ImageDetail{&contentservice.ImageDetail{}}, nil
}

//SetImageDetails sets the ImageDetails to cache for a given key
func (i *Info) SetImageDetails(key string, imageDetails []*contentservice.ImageDetail) error {
	return nil
}
