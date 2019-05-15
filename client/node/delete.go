package node

import (
	"fmt"
	"github.com/bububa/depthtree/client"
)

type DeleteRequest struct {
	Id int64 `json:"id"`
}

func Delete(client *client.Client, db string, nodeId int64) error {
	endPoint := fmt.Sprintf("/node/delete/%s", db)
	req := DeleteRequest{
		Id: nodeId,
	}
	_, err := client.Post(endPoint, req)
	return err
}
