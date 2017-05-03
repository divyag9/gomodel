package list

import (
	"testing"

	ora "gopkg.in/rana/ora.v4"

	"github.com/divyag9/gomodel/pkg/cache"
	"github.com/divyag9/gomodel/pkg/pb"
)

type FakeDatabaseInfo struct {
	Session *ora.Ses
}

func (i *FakeDatabaseInfo) GetImageDetailsByOrderNumber(orderNumber int64) ([]*contentservice.ImageDetail, error) {
	return nil, nil
}

func (i *FakeDatabaseInfo) GetImageDetailsByImageIds(imageIds []int64) ([]*contentservice.ImageDetail, error) {
	return nil, nil
}

func TestGetImageDetails(t *testing.T) {
	orderNumberInfo := &OrderNumberInfo{
		OrderNumber: 1,
		Database:    &FakeDatabaseInfo{},
		Cache:       &cache.Info{},
	}
	imageDetails, err := orderNumberInfo.GetImageDetails()

}
