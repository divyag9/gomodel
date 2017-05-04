package contentserviceinterface

import "github.com/divyag9/gomodel/pkg/pb"

// OrderNumberGetter contains methods to retrieve imageDetails for ordernumber
type OrderNumberGetter interface {
	GetImageDetailsByOrderNumber(orderNumber int64) ([]*contentservice.ImageDetail, error)
}

// ImageIdsGetter contains methods to retrieve imageDetails for ImageIds
type ImageIdsGetter interface {
	GetImageDetailsByImageIds(imageIds []int64) ([]*contentservice.ImageDetail, error)
}
