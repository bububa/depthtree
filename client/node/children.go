package node

import (
	"encoding/json"
	"fmt"
	"github.com/bububa/depthtree"
	"github.com/bububa/depthtree/client"
)

func Children(client *client.Client, db string, nodeId int64, depth int) ([]*depthtree.Node, error) {
	endPoint := fmt.Sprintf("/children/%s/%d/%d", db, nodeId, depth)
	js, err := client.Get(endPoint)
	var children []*depthtree.Node
	err = json.Unmarshal(js, &children)
	if err != nil {
		return nil, err
	}
	return children, nil
}
