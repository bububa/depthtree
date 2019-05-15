package db

import (
	"encoding/json"
	"fmt"
	"github.com/bububa/depthtree"
	"github.com/bububa/depthtree/client"
)

func Roots(client *client.Client, db string) ([]*depthtree.Node, error) {
	endPoint := fmt.Sprintf("/db/roots/%s", db)
	js, err := client.Get(endPoint)
	var roots []*depthtree.Node
	err = json.Unmarshal(js, &roots)
	if err != nil {
		return nil, err
	}
	return roots, nil
}
