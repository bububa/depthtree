package node

import (
	"fmt"
	"github.com/bububa/depthtree/client"
)

func BatchAdd(client *client.Client, db string, req []AddRequest) error {
	endPoint := fmt.Sprintf("/node/batch-add/%s", db)
	_, err := client.Post(endPoint, req)
	return err
}
