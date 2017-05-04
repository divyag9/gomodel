package cache

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/divyag9/gomodel/pkg/pb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	ora "gopkg.in/rana/ora.v4"
)

type FakeDatabaseConfig struct {
	Session *ora.Ses
}

func (f *FakeDatabaseConfig) GetImageDetailsByOrderNumber(orderNumber int64) ([]*contentservice.ImageDetail, error) {
	return []*contentservice.ImageDetail{{
		Archived:               "",
		Category:               "foo",
		ContractorId:           72494,
		DateCreated:            testDate,
		DateModified:           testDate,
		ImageUTCDate:           nil,
		ImageTakenDate:         nil,
		DeptCode:               "01",
		DescPrefix:             "foo",
		DescText:               "foo bar",
		FileSize:               180,
		ImageId:                3001240405,
		ImageFileName:          "C:\\Temp\\images600\\016\\555\\76215592-b810-48f0-a9e2-ac681ab0ea38.png",
		ImageHeight:            100,
		ImageType:              1,
		ImageRotated:           false,
		ImageWidth:             100,
		OrderNumber:            600016555,
		ReleaseDate:            testDate,
		ScanDate:               testDate,
		ThumbnailSize:          0,
		WebFileName:            "images/600/016/555/76215592-b810-48f0-a9e2-ac681ab0ea38.png",
		MimeType:               "image/png",
		GeneratedImageFilePath: "https://sbimage.sgpdev.com/images/600/016/555/76215592-b810-48f0-a9e2-ac681ab0ea38.png?d794b9a5-02ac-4b86-be81-cbbf0d22abf7",
		Guid: "",
	}}, nil
}

var cases = []struct {
	cacheConfig          *Config
	orderNumber          int64
	expectedImageDetails []*contentservice.ImageDetail
	expectedError        error
}{
	{
		cacheConfig: &Config{
			Memcache:          getMemcacheClient(),
			SecondsToExpiry:   50,
			OrderNumberGetter: &FakeDatabaseConfig{},
		},
		orderNumber: 600016555,
		expectedImageDetails: []*contentservice.ImageDetail{{
			Archived:               "",
			Category:               "foo",
			ContractorId:           72494,
			DateCreated:            testDate,
			DateModified:           testDate,
			ImageUTCDate:           nil,
			ImageTakenDate:         nil,
			DeptCode:               "01",
			DescPrefix:             "foo",
			DescText:               "foo bar",
			FileSize:               180,
			ImageId:                3001240405,
			ImageFileName:          "C:\\Temp\\images600\\016\\555\\76215592-b810-48f0-a9e2-ac681ab0ea38.png",
			ImageHeight:            100,
			ImageType:              1,
			ImageRotated:           false,
			ImageWidth:             100,
			OrderNumber:            600016555,
			ReleaseDate:            testDate,
			ScanDate:               testDate,
			ThumbnailSize:          0,
			WebFileName:            "images/600/016/555/76215592-b810-48f0-a9e2-ac681ab0ea38.png",
			MimeType:               "image/png",
			GeneratedImageFilePath: "https://sbimage.sgpdev.com/images/600/016/555/76215592-b810-48f0-a9e2-ac681ab0ea38.png?d794b9a5-02ac-4b86-be81-cbbf0d22abf7",
			Guid: "",
		}},
		expectedError: nil,
	},
}

func getMemcacheClient() *memcache.Client {
	servers := os.Getenv("MEMCACHE_SERVERS")
	memcacheServers := strings.Split(servers, ",")
	mc := memcache.New(memcacheServers...)

	return mc
}

var testDate *timestamp.Timestamp

func getDate() *timestamp.Timestamp {
	testDate, _ = ptypes.TimestampProto(time.Date(2017, 03, 14, 20, 52, 45, 0, time.UTC))
	return testDate
}

func TestGetImageDetailsByOrderNumber(t *testing.T) {
	for _, c := range cases {
		imageDetails, err := c.cacheConfig.GetImageDetailsByOrderNumber(c.orderNumber)

		if !reflect.DeepEqual(err, c.expectedError) {
			t.Errorf("Expected err to be %q but it was %q", c.expectedError, err)
		}
		if !reflect.DeepEqual(c.expectedImageDetails, imageDetails) {
			t.Errorf("Expected %q but got %q", c.expectedImageDetails, imageDetails)
		}
	}
}
