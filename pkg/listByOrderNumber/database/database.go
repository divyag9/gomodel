package database

import (
	"strconv"
	"time"

	"github.com/divyag9/gomodel/pkg/pb"
	"github.com/golang/protobuf/ptypes"
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
			value, err := strconv.ParseInt(resultSet.Row[i].(ora.OCINum).String(), 10, 64)
			if err != nil {
				return nil, err
			}
			imageDetail.ImageId = value
		case "CONTRACTORID":
			if resultSet.Row[i] != nil {
				value, err := strconv.ParseInt(resultSet.Row[i].(ora.OCINum).String(), 10, 64)
				if err != nil {
					return nil, err
				}
				imageDetail.ContractorId = value
			}
		case "ORDERNUMBER":
			if resultSet.Row[i] != nil {
				value, err := strconv.ParseInt(resultSet.Row[i].(ora.OCINum).String(), 10, 64)
				if err != nil {
					return nil, err
				}
				imageDetail.OrderNumber = value
			}
		case "SCANDATE":
			if resultSet.Row[i] != nil {
				scanDate, err := ptypes.TimestampProto(resultSet.Row[i].(time.Time))
				if err != nil {
					return nil, err
				}
				imageDetail.ScanDate = scanDate
			}
		case "RELEASEDATE":
			if resultSet.Row[i] != nil {
				releaseDate, err := ptypes.TimestampProto(resultSet.Row[i].(time.Time))
				if err != nil {
					return nil, err
				}
				imageDetail.ReleaseDate = releaseDate
			}
		case "IMAGEFILENAME":
			if resultSet.Row[i] != nil {
				imageDetail.ImageFileName = resultSet.Row[i].(string)
			}
		case "WEBFILENAME":
			if resultSet.Row[i] != nil {
				imageDetail.WebFileName = resultSet.Row[i].(string)
			}
		case "IMAGETYPE":
			if resultSet.Row[i] != nil {
				imageDetail.ImageType = int32(resultSet.Row[i].(int64))
			}
		case "IMAGEWIDTH":
			if resultSet.Row[i] != nil {
				imageDetail.ImageWidth = int32(resultSet.Row[i].(int64))
			}
		case "IMAGEHEIGHT":
			if resultSet.Row[i] != nil {
				imageDetail.ImageHeight = int32(resultSet.Row[i].(int64))
			}
		case "DEPTCODE":
			if resultSet.Row[i] != nil {
				imageDetail.DeptCode = resultSet.Row[i].(string)
			}
		case "ARCHIVED":
			if resultSet.Row[i] != nil {
				imageDetail.Archived = resultSet.Row[i].(string)
			}
		case "FILESIZE":
			if resultSet.Row[i] != nil {
				value, err := strconv.ParseInt(resultSet.Row[i].(ora.OCINum).String(), 10, 32)
				if err != nil {
					return nil, err
				}
				imageDetail.FileSize = int32(value)
			}
		case "THUMBNAILSIZE":
			if resultSet.Row[i] != nil {
				value, err := strconv.ParseInt(resultSet.Row[i].(ora.OCINum).String(), 10, 32)
				if err != nil {
					return nil, err
				}
				imageDetail.ThumbnailSize = int32(value)
			}
		case "DATECREATED":
			if resultSet.Row[i] != nil {
				dateCreated, err := ptypes.TimestampProto(resultSet.Row[i].(time.Time))
				if err != nil {
					return nil, err
				}
				imageDetail.DateCreated = dateCreated
			}
		case "DATEMODIFIED":
			if resultSet.Row[i] != nil {
				dateModefied, err := ptypes.TimestampProto(resultSet.Row[i].(time.Time))
				if err != nil {
					return nil, err
				}
				imageDetail.DateModified = dateModefied
			}
		case "DESCPREFIX":
			if resultSet.Row[i] != nil {
				imageDetail.DescPrefix = resultSet.Row[i].(string)
			}
		case "DESCTEXT":
			if resultSet.Row[i] != nil {
				imageDetail.DescText = resultSet.Row[i].(string)
			}
		case "CATEGORY":
			if resultSet.Row[i] != nil {
				imageDetail.Category = resultSet.Row[i].(string)
			}
		case "GUID":
			if resultSet.Row[i] != nil {
				imageDetail.Guid = resultSet.Row[i].(string)
			}
		}
	}
	//Generate these
	//imageDetail.MimeType,
	//imageDetail.GeneratedImageFilePath,
	//imagerotated
	return imageDetail, nil
}
