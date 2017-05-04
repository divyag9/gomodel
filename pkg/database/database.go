package database

import (
	"fmt"

	"github.com/divyag9/gomodel/pkg/pb"
	ora "gopkg.in/rana/ora.v4"
)

//Config contains information required for the database operations
type Config struct {
	Session *ora.Ses
}

//NewConfig initializes the Config struct and returns it
func NewConfig(session *ora.Ses) *Config {
	return &Config{
		Session: session,
	}
}

//GetImageDetailsByOrderNumber retrieves the imageDetails for a given orderNumber
func (c *Config) GetImageDetailsByOrderNumber(orderNumber int64) ([]*contentservice.ImageDetail, error) {
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

//GetImageDetailsByImageIds retrieves the imageDetails for a given set og imageIds
func (c *Config) GetImageDetailsByImageIds(imageIds []int64) ([]*contentservice.ImageDetail, error) {
	//To be implemented
	return []*contentservice.ImageDetail{&contentservice.ImageDetail{}}, nil
}
