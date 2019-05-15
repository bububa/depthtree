package node

import (
	"encoding/json"
	"fmt"
	"github.com/bububa/depthtree"
	"github.com/bububa/depthtree/client"
)

func Parents(client *client.Client, db string, nodeId int64) ([]*depthtree.Node, error) {
	endPoint := fmt.Sprintf("/parents/%s/%d", db, nodeId)
	js, err := client.Get(endPoint)
	var parents []*depthtree.Node
	err = json.Unmarshal(js, &parents)
	if err != nil {
		return nil, err
	}
	return parents, nil
}
