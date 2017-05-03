package list

import (
	"fmt"
	"reflect"

	"github.com/divyag9/gomodel/pkg/cache"
	"github.com/divyag9/gomodel/pkg/database"
	"github.com/divyag9/gomodel/pkg/pb"
)

//OrderNumberInfo contains information required to retrieve imageDetails by orderNumber
type OrderNumberInfo struct {
	OrderNumber  int64
	DatabaseInfo *database.Info
	CacheInfo    *cache.Info
}

//GetImageDetails retrieves image details for an ordernumber. It will first see if the details are in the cache if not gets the details from database
func (o *OrderNumberInfo) GetImageDetails() ([]*contentservice.ImageDetail, error) {
	// Retrive the imagedetails from cache if it exists
	imageDetails, err := o.CacheInfo.GetImageDetails("")
	if err != nil {
		//Making a call to the database cause the results are not in cache
		if reflect.DeepEqual(err, fmt.Errorf("memcache: cache miss")) {

			imageDetails, err = o.DatabaseInfo.GetImageDetailsByOrderNumber(o.OrderNumber)
			if err != nil {
				return nil, err
			}
			//Set the imageDetails to cache
			err = o.CacheInfo.SetImageDetails("", imageDetails)
			if err != nil {
				return nil, err
			}
			return imageDetails, nil
		}
		return nil, err
	}

	return imageDetails, nil
}
