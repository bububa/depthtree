package db

import (
	"encoding/json"
	"fmt"
	"github.com/bububa/depthtree"
	"github.com/bububa/depthtree/client"
)

func TopChildren(client *client.Client, db string, depth int, limit int) ([]*depthtree.Node, error) {
	endPoint := fmt.Sprintf("/db/top-children/%s/%d/%d", db, depth, limit)
	js, err := client.Get(endPoint)
	var nodes []*depthtree.Node
	err = json.Unmarshal(js, &nodes)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}
