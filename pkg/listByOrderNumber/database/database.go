package database

import (
	"github.com/divyag9/gomodel/pkg/pb"
	"github.com/divyag9/gomodel/pkg/utility"
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

	//Creating imageDetails from resultSet
	imageDetails := []*contentservice.ImageDetail{}
	if resultSet.IsOpen() {
		for resultSet.Next() {
			imageDetail, err := utility.GetImageDetailFromResultSet(resultSet)
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
