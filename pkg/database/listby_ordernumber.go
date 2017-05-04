package database

import (
	"fmt"
	"time"

	"github.com/divyag9/gomodel/pkg/pb"
	"github.com/golang/protobuf/ptypes"
	ora "gopkg.in/rana/ora.v4"
)

//Client contains information required for the database operations
type Client struct {
	Session *ora.Ses
}

//FakeDatabaseClient for testing
type FakeDatabaseClient struct {
	Session *ora.Ses
}

//New initializes the Client struct and returns it
func New(session *ora.Ses) *Client {
	return &Client{
		Session: session,
	}
}

//GetImageDetailsByOrderNumber retrieves the imageDetails for a given orderNumber
func (c *Client) GetImageDetailsByOrderNumber(orderNumber int64) ([]*contentservice.ImageDetail, error) {
	//Prepare the query
	prepStatement, err := c.Session.Prep("CALL CONTENTSERVICE.RETRIEVEDEPARTMENTS(:1)")
	if err != nil {
		return nil, err
	}
	defer prepStatement.Close()

	//Retrieve the resultSet
	resultSet := &ora.Rset{}
	rownum, err := prepStatement.Exe(resultSet)
	if err != nil {
		return nil, err
	}

	//Trying to print the values retured from ref_cursor
	fmt.Println(rownum)
	fmt.Println(resultSet.Row)
	for _, v := range resultSet.Columns {
		fmt.Println(v.Name)
	}

	return []*contentservice.ImageDetail{}, nil
}

//GetImageDetailsByOrderNumber fake implementation for testing
func (f *FakeDatabaseClient) GetImageDetailsByOrderNumber(orderNumber int64) ([]*contentservice.ImageDetail, error) {
	testDate, _ := ptypes.TimestampProto(time.Date(2017, 03, 14, 20, 52, 45, 0, time.UTC))

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
