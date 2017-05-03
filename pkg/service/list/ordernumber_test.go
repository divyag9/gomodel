package list

import (
	"testing"

	ora "gopkg.in/rana/ora.v4"

	"github.com/divyag9/gocontentservice/pkg/contentservice"
	"github.com/divyag9/gomodel/pkg/cache"
)

type FakeDatabaseInfo struct {
	Session *ora.Ses
}

type FakeInfo struct {
	OrderNumber  int64
	DatabaseInfo *FakeDatabaseInfo
	CacheInfo    *cache.Info
}

func (i *FakeDatabaseInfo) GetImageDetailsByOrderNumber(orderNumber int64) ([]*contentservice.ImageDetail, error) {
	return nil, nil
}

func TestGetImageDetails(t *testing.T) {
	fakeInfo := &FakeInfo{
		OrderNumber:  1,
		DatabaseInfo: &FakeDatabaseInfo{},
		CacheInfo:    &cache.Info{},
	}
	imageDetails, err := fakeInfo.GetImageDetails()

}
