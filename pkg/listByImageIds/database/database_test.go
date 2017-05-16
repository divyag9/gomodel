package database

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/divyag9/gomodel/pkg/pb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	ora "gopkg.in/rana/ora.v4"
)

var imageDetailDatabaseCases = []struct {
	ImageIds             []int64
	Client               *Client
	expectedImageDetails []*contentservice.ImageDetail
	expectedError        error
}{
	{
		ImageIds: []int64{3001271328},
		Client: &Client{
			Session: getDbSession(),
		},
		expectedImageDetails: []*contentservice.ImageDetail{{
			Archived:      "N",
			Category:      "Inspection",
			ContractorId:  0,
			DateCreated:   getTimestamp(time.Date(2017, 05, 10, 16, 47, 25, 0, time.Local)),
			DateModified:  getTimestamp(time.Date(2017, 05, 10, 16, 47, 25, 0, time.Local)),
			DeptCode:      "01",
			DescPrefix:    "Condition",
			DescText:      "Auto",
			FileSize:      180,
			ImageId:       3001271328,
			ImageFileName: "C:\\Temp\\images600\\016\\556\\b1dd3ec3-61be-4ecf-adee-22685054c953.png",
			ImageHeight:   428,
			ImageType:     5,
			ImageWidth:    640,
			OrderNumber:   600016556,
			ReleaseDate:   getTimestamp(time.Date(2017, 05, 10, 10, 10, 10, 0, time.Local)),
			ScanDate:      getTimestamp(time.Date(2017, 05, 10, 16, 47, 25, 0, time.Local)),
			ThumbnailSize: 0,
			WebFileName:   "images/600/016/556/b1dd3ec3-61be-4ecf-adee-22685054c953.png",
			//MimeType:               "image/png",
			//GeneratedImageFilePath: "https://sbimage.sgpdev.com/images/600/016/555/76215592-b810-48f0-a9e2-ac681ab0ea38.png?d794b9a5-02ac-4b86-be81-cbbf0d22abf7",
			Guid: "",
		}},
		expectedError: nil,
	},
}

func getDbSession() *ora.Ses {
	dsn := os.Getenv("GO_OCI8_LIB_CONNECT_STRING")
	_, _, ses, err := ora.NewEnvSrvSes(dsn)
	if err != nil {
		fmt.Println(err)
	}

	return ses
}

func getTimestamp(date time.Time) *timestamp.Timestamp {
	dateTimestamp, _ := ptypes.TimestampProto(date)

	return dateTimestamp
}
func TestGetImageDetailsByImageIds(t *testing.T) {
	for _, c := range imageDetailDatabaseCases {
		imageDetails, err := c.Client.GetImageDetailsByImageIds(c.ImageIds)
		defer c.Client.Session.Close()

		if !reflect.DeepEqual(err, c.expectedError) {
			t.Errorf("Expected err to be %q but it was %q", c.expectedError, err)
		}
		if !reflect.DeepEqual(c.expectedImageDetails, imageDetails) {
			t.Errorf("Expected %q but got %q", c.expectedImageDetails, imageDetails)
		}
		//Look into testing numFreeConnections
	}
}
