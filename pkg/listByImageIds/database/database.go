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

type numbertable struct {
	nums []int64
}

//New initializes the Client struct and returns it
func New(session *ora.Ses) *Client {
	return &Client{
		Session: session,
	}
}

//GetImageDetailsByImageIds retrieves the imageDetails for given imageIds
func (c *Client) GetImageDetailsByImageIds(imageIds []int64) ([]*contentservice.ImageDetail, error) {
	//Prepare the query
	prepStatement, err := c.Session.Prep("CALL CONTENTSERVICE.RETRIEVEIMAGEDETAILLISTFROMIDS(:iImageDetailIdTable,:genRefCursor)")
	if err != nil {
		return nil, err
	}
	defer prepStatement.Close()
	//env := &ora.Env{}
	//orastruct := &ora.OraOCINum
	a := make([]ora.Float64, len(imageIds))
	for i := 0; i < len(imageIds); i++ {
		a[i] = ora.Float64{Value: float64(imageIds[i])}
		//ora.Int64{}
		//a[i] = ora.OraOCINum(imageIds[i])
		//env.OCINumberFromInt(&a[i], imageIds[i], 1)
	}
	stmtSliceIns, err := c.Session.Prep("NUMBER_TABLE(3001271328)")
	defer stmtSliceIns.Close()
	if err != nil {
		return nil, err
	}
	// rowsAffected, err := stmtSliceIns.Exe(imageIds)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println(rowsAffected)

	//Retrieve the resultSet
	numtab := &numbertable{nums: imageIds}
	resultSet := &ora.Rset{}
	_, err = prepStatement.Exe(numtab, resultSet)
	if err != nil {
		return nil, err
	}

	//Create the imageDetails array
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
