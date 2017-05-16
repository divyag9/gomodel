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
	Memcache        *memcache.Client
	SecondsToExpiry int32
	ImageIdsGetter  contentserviceinterface.ImageIdsGetter
}

//New initializes the cache client to get Imagedetails for array of imageIds
func New(memcache *memcache.Client, secondsToExpiry int32, imageIdsGetter contentserviceinterface.ImageIdsGetter) *Client {
	return &Client{
		Memcache:        memcache,
		SecondsToExpiry: secondsToExpiry,
		ImageIdsGetter:  imageIdsGetter,
	}
}

//GetImageDetailsByImageIds retrieves ImageDetails for given imageIds
func (c *Client) GetImageDetailsByImageIds(imageIds []int64) ([]*contentservice.ImageDetail, error) {
	//Construct keys for caching
	keys := getKeys(imageIds)

	//Get ImageDetail's from cache in bulk
	imageIdsCache, imageDetails, err := getImageDetailsMulti(keys, c.Memcache)
	if err != nil {
		return nil, err
	}

	//Get imageIds whose imageDetails could not be retrieved from cache
	imageIdsDatabase := getImageIdsToRetieveFromDatabase(imageIds, imageIdsCache)

	//Get imageDetails from database, for ids that couldn't be retrieved from cache
	if len(imageIdsDatabase) > 0 {
		imageDetailsDatabase, err := c.ImageIdsGetter.GetImageDetailsByImageIds(imageIdsDatabase)
		if err != nil {
			return nil, err
		}

		//Set cache for each of the imageIds retrieved from database
		for _, v := range imageDetailsDatabase {
			key := fmt.Sprintf("Id:%q", v.ImageId)
			err = setImageDetails(key, v, c)
			if err != nil {
				return nil, err
			}
		}

		imageDetails = append(imageDetails, imageDetailsDatabase...)
		fmt.Println("from databse")
	}

	return imageDetails, nil
}

//getKeys returns array of key strings to be passed to GetMulti
func getKeys(imageIds []int64) []string {
	keys := make([]string, len(imageIds))
	for i, v := range imageIds {
		keys[i] = fmt.Sprintf("Id:%q", v)
	}

	return keys
}

//getImageIdsToRetieveFromDatabase returns array of imageids whose imagedetails could not be retrieved from cache.
//returnes array of imageIds is used to retieve from database
func getImageIdsToRetieveFromDatabase(imageIds []int64, imageIdsCache []int64) []int64 {
	m := make(map[int64]int)
	for _, imageID := range imageIdsCache {
		m[imageID] = 1
	}
	var imageIdsDatabase []int64
	for _, imageID := range imageIds {
		if m[imageID] == 0 {
			imageIdsDatabase = append(imageIdsDatabase, imageID)
		}
	}

	return imageIdsDatabase
}

//getImageDetailsMulti retrieves the ImageDetails from cache in bulk
func getImageDetailsMulti(keys []string, memcacheClient *memcache.Client) ([]int64, []*contentservice.ImageDetail, error) {
	//Get imageDetails for given keys from cache
	keysToImageDetailMap, err := memcacheClient.GetMulti(keys)
	if err != nil {
		return nil, nil, err
	}
	var imageIdsCache []int64
	imageDetails := []*contentservice.ImageDetail{}
	for _, v := range keysToImageDetailMap {
		//Deserialize ImageDetail bytes to ImageDetail struct
		imageDetailBytes := v.Value
		imageDetail := &contentservice.ImageDetail{}
		decBuf := bytes.NewBuffer(imageDetailBytes)
		err = gob.NewDecoder(decBuf).Decode(&imageDetail)
		if err != nil {
			return nil, nil, err
		}

		imageDetails = append(imageDetails, imageDetail)
		imageIdsCache = append(imageIdsCache, imageDetail.ImageId)
	}

	fmt.Println("from cache")

	return imageIdsCache, imageDetails, nil
}

//setImageDetails sets the ImageDetails to cache for a given key
func setImageDetails(key string, imageDetail *contentservice.ImageDetail, client *Client) error {
	//Serialize ImageDetails into bytes
	encBuf := new(bytes.Buffer)
	err := gob.NewEncoder(encBuf).Encode(imageDetail)
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
