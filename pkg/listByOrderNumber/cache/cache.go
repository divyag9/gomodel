package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/divyag9/gomodel/pkg/interface"
	"github.com/divyag9/gomodel/pkg/pb"
)

// Client holds all the information required for caching
type Client struct {
	Memcache          *memcache.Client
	SecondsToExpiry   int32
	OrderNumberGetter contentserviceinterface.OrderNumberGetter
}

//New initializes the cache client to get Imagedetails for ordernumber
func New(memcache *memcache.Client, secondsToExpiry int32, orderNumberGetter contentserviceinterface.OrderNumberGetter) *Client {
	return &Client{
		Memcache:          memcache,
		SecondsToExpiry:   secondsToExpiry,
		OrderNumberGetter: orderNumberGetter,
	}
}

//GetImageDetailsByOrderNumber retrieves the ImageDetails for a given orderNumber
func (c *Client) GetImageDetailsByOrderNumber(orderNumber int64) ([]*contentservice.ImageDetail, error) {
	//Construct the key for caching
	key := fmt.Sprintf("OrderNumber:%q", orderNumber)

	//Get ImageDetails from cache for a given key
	imageDetails, err := getImageDetails(key, c.Memcache)
	if err != nil {
		//Making a call to the database cause the results are not in cache
		imageDetails, err = c.OrderNumberGetter.GetImageDetailsByOrderNumber(orderNumber)
		if err != nil {
			return nil, err
		}

		//Set ImageDetails retrieved from database to cache
		err = setImageDetails(key, imageDetails, c)
		if err != nil {
			return nil, err
		}
	}

	return imageDetails, nil
}

//getImageDetails retrieves the ImageDetails from cache for a given key
func getImageDetails(key string, memcacheClient *memcache.Client) ([]*contentservice.ImageDetail, error) {
	//Retieve value for a key from memcache
	item, err := memcacheClient.Get(key)
	if err != nil {
		return nil, err
	}
	imageDetailBytes := item.Value

	//Deserialize ImageDetail bytes to ImageDetail struct
	imageDetails := []*contentservice.ImageDetail{}
	decBuf := bytes.NewBuffer(imageDetailBytes)
	err = gob.NewDecoder(decBuf).Decode(&imageDetails)
	if err != nil {
		return nil, err
	}

	return imageDetails, nil
}

//setImageDetails writes the ImageDetails to cache for a given key
func setImageDetails(key string, imageDetails []*contentservice.ImageDetail, client *Client) error {
	//Serialize ImageDetails into bytes
	encBuf := new(bytes.Buffer)
	err := gob.NewEncoder(encBuf).Encode(imageDetails)
	if err != nil {
		return err
	}

	//Setting ImageDetail bytes to cache
	err = client.Memcache.Set(&memcache.Item{Key: key, Value: encBuf.Bytes(), Expiration: client.SecondsToExpiry})
	if err != nil {
		return err
	}

	return nil
}
