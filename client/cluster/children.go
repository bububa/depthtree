package cluster

import (
	"encoding/json"
	"fmt"
	"github.com/bububa/depthtree"
	"github.com/bububa/depthtree/client"
)

func Children(client *client.Client, db string, depth int, k int, limit int) ([]*depthtree.Cluster, error) {
	endPoint := fmt.Sprintf("/cluster/children/%s/%d/%d?limit=%d", db, depth, k, limit)
	js, err := client.Get(endPoint)
	var clusters []*depthtree.Cluster
	err = json.Unmarshal(js, &clusters)
	if err != nil {
		return nil, err
	}
	return clusters, nil
}
