package list

import (
	"github.com/divyag9/gomodel/pkg/cache"
	"github.com/divyag9/gomodel/pkg/database"
)

//OrderNumberInfo contains information required to retrieve imageDetails by orderNumber
type OrderNumberInfo struct {
	OrderNumber int64
	Database    database.Caller
	Cache       cache.Caller
}

//GetImageDetails retrieves image details for an ordernumber. It will first see if the details are in the cache if not gets the details from database
func (o *OrderNumberInfo) GetImageDetailsByOrder(Orders []OrderNumbers) ([]*contentservice.ImageDetail, error) {
	// Retrive the imagedetails from cache if it exists
	var results []*contentservice.ImageDetail
	for _, v := range Orders {
		imageDetails, err = o.Database.GetImageDetailsByOrderNumber(o.OrderNumber)
		if err != nil {
			continue
		}
		results = append(results, imageDetails)
	}
	return results, nil
}
