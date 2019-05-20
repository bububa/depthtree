package db

import (
	"encoding/json"
	"fmt"
	"github.com/bububa/depthtree/client"
)

func List(client *client.Client) ([]string, error) {
	endPoint := fmt.Sprintf("/db/list")
	js, err := client.Get(endPoint)
	var dbs []string
	err = json.Unmarshal(js, &dbs)
	if err != nil {
		return nil, err
	}
	return dbs, err
}
