package database

import (
	"strconv"
	"time"

	"github.com/divyag9/gomodel/pkg/pb"
	"github.com/golang/protobuf/ptypes"
	google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
	ora "gopkg.in/rana/ora.v4"
)

//Client contains information required for the database operations
type Client struct {
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
	prepStatement, err := c.Session.Prep("CALL CONTENTSERVICE.RETRIEVEIMAGEDETAILLIST(:iordernumber,:genRefCursor)")
	if err != nil {
		return nil, err
	}
	defer prepStatement.Close()

	//Retrieve the resultSet
	resultSet := &ora.Rset{}
	_, err = prepStatement.Exe(orderNumber, resultSet)
	if err != nil {
		return nil, err
	}

	//Create the imageDetails array
	imageDetails := []*contentservice.ImageDetail{}
	if resultSet.IsOpen() {
		for resultSet.Next() {
			imageDetail, err := getImageDetailFromResultSet(resultSet)
			if err != nil {
				return nil, err
			}
			imageDetails = append(imageDetails, imageDetail)
		}
		if err := resultSet.Err(); err != nil {
			return nil, err
		}
	}

	return imageDetails, nil
}

//getImageDetailFromResultSet creates the ImageDetail struct from the database result set
func getImageDetailFromResultSet(resultSet *ora.Rset) (*contentservice.ImageDetail, error) {
	imageDetail := &contentservice.ImageDetail{}
	for i, v := range resultSet.Columns {
		switch v.Name {
		case "ID":
			value, err := toInt64FromOCINum(resultSet.Row[i])
			if err != nil {
				return nil, err
			}
			imageDetail.ImageId = value
		case "CONTRACTORID":
			value, err := toInt64FromOCINum(resultSet.Row[i])
			if err != nil {
				return nil, err
			}
			imageDetail.ContractorId = value
		case "ORDERNUMBER":
			value, err := toInt64FromOCINum(resultSet.Row[i])
			if err != nil {
				return nil, err
			}
			imageDetail.OrderNumber = value
		case "SCANDATE":
			scanDate, err := toTimestamp(resultSet.Row[i])
			if err != nil {
				return nil, err
			}
			imageDetail.ScanDate = scanDate
		case "RELEASEDATE":
			releaseDate, err := toTimestamp(resultSet.Row[i])
			if err != nil {
				return nil, err
			}
			imageDetail.ReleaseDate = releaseDate
		case "IMAGEFILENAME":
			if resultSet.Row[i] != nil {
				imageDetail.ImageFileName = toString(resultSet.Row[i])
			}
		case "WEBFILENAME":
			if resultSet.Row[i] != nil {
				imageDetail.WebFileName = toString(resultSet.Row[i])
			}
		case "IMAGETYPE":
			if resultSet.Row[i] != nil {
				imageDetail.ImageType = toInt32FromNumber(resultSet.Row[i])
			}
		case "IMAGEWIDTH":
			if resultSet.Row[i] != nil {
				imageDetail.ImageWidth = toInt32FromNumber(resultSet.Row[i])
			}
		case "IMAGEHEIGHT":
			if resultSet.Row[i] != nil {
				imageDetail.ImageHeight = toInt32FromNumber(resultSet.Row[i])
			}
		case "DEPTCODE":
			if resultSet.Row[i] != nil {
				imageDetail.DeptCode = toString(resultSet.Row[i])
			}
		case "ARCHIVED":
			if resultSet.Row[i] != nil {
				imageDetail.Archived = toString(resultSet.Row[i])
			}
		case "FILESIZE":
			value, err := toInt32FromOCINum(resultSet.Row[i])
			if err != nil {
				return nil, err
			}
			imageDetail.FileSize = value
		case "THUMBNAILSIZE":
			value, err := toInt32FromOCINum(resultSet.Row[i])
			if err != nil {
				return nil, err
			}
			imageDetail.ThumbnailSize = value
		case "DATECREATED":
			dateCreated, err := toTimestamp(resultSet.Row[i])
			if err != nil {
				return nil, err
			}
			imageDetail.DateCreated = dateCreated
		case "DATEMODIFIED":
			dateModefied, err := toTimestamp(resultSet.Row[i])
			if err != nil {
				return nil, err
			}
			imageDetail.DateModified = dateModefied
		case "DESCPREFIX":
			imageDetail.DescPrefix = toString(resultSet.Row[i])
		case "DESCTEXT":
			imageDetail.DescText = toString(resultSet.Row[i])
		case "CATEGORY":
			imageDetail.Category = toString(resultSet.Row[i])
		case "GUID":
			imageDetail.Guid = toString(resultSet.Row[i])
		}
	}
	//Generate these
	//imageDetail.MimeType,
	//imageDetail.GeneratedImageFilePath,
	//imagerotated
	return imageDetail, nil
}

//Convert resultset value of columns that are of type string
func toString(columnValue interface{}) string {
	var result string
	if columnValue != nil {
		result = columnValue.(string)
	}

	return result
}

//Convert resultset value of columns that are of type Date
func toTimestamp(columnValue interface{}) (*google_protobuf.Timestamp, error) {
	var result *google_protobuf.Timestamp
	if columnValue != nil {
		timestampValue, err := ptypes.TimestampProto(columnValue.(time.Time))
		if err != nil {
			return nil, err
		}
		result = timestampValue
	}
	return result, nil
}

//Convert resultset value of columns that are of type OCINum
func toInt64FromOCINum(columnValue interface{}) (int64, error) {
	var result int64
	if columnValue != nil {
		parseValue, err := strconv.ParseInt(columnValue.(ora.OCINum).String(), 10, 64)
		if err != nil {
			return 0, err
		}
		result = parseValue
	}

	return result, nil
}

//Convert resultset value of columns that are of type OCINum
func toInt32FromOCINum(columnValue interface{}) (int32, error) {
	var result int32
	if columnValue != nil {
		parseValue, err := strconv.ParseInt(columnValue.(ora.OCINum).String(), 10, 32)
		if err != nil {
			return 0, err
		}
		result = int32(parseValue)
	}

	return result, nil
}

//Convert resultset value of columns that are of type Number
func toInt32FromNumber(columnValue interface{}) int32 {
	var result int32
	if columnValue != nil {
		result = int32(columnValue.(int64))
	}

	return result
}
