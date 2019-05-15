package db

import (
	"fmt"
	"github.com/bububa/depthtree/client"
)

type TruncateRequest struct {
	Name string `json:"name" binding:"required"`
}

func Truncate(client *client.Client, db string) error {
	endPoint := fmt.Sprintf("/db/truncate")
	req := TruncateRequest{
		Name: db,
	}
	_, err := client.Post(endPoint, req)
	return err
}
