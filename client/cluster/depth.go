package cluster

import (
	"encoding/json"
	"fmt"
	"github.com/bububa/depthtree"
	"github.com/bububa/depthtree/client"
)

func Depth(client *client.Client, db string, k int) ([]*depthtree.Cluster, error) {
	endPoint := fmt.Sprintf("/cluster/depth/%s/%d", db, k)
	js, err := client.Get(endPoint)
	var clusters []*depthtree.Cluster
	err = json.Unmarshal(js, &clusters)
	if err != nil {
		return nil, err
	}
	return clusters, nil
}
