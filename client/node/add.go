package node

import (
	"fmt"
	"github.com/bububa/depthtree/client"
)

type AddRequest struct {
	Id  int64
	Pid int64
}

func Add(client *client.Client, db string, req *AddRequest) error {
	endPoint := fmt.Sprintf("/node/add/%s", db)
	_, err := client.Post(endPoint, req)
	return err
}
