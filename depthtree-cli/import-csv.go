package main

import (
	"encoding/csv"
	"github.com/bububa/depthtree/client"
	"github.com/bububa/depthtree/client/node"
	"github.com/mkideal/log"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

func ImportCsv(c *client.Client, filePath string) error {
	if !path.IsAbs(filePath) {
		wd, err := os.Getwd()
		if err != nil {
			log.Error(err.Error())
			return err
		}
		filePath = path.Join(wd, filePath)
	}
	f, err := os.Open(filePath)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer f.Close()
	_, filename := path.Split(filePath)
	ext := path.Ext(filename)
	dbname := strings.TrimSuffix(filename, ext)
	log.Info("db: %s", dbname)
	reader := csv.NewReader(f)
	var nodes []node.AddRequest
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Error(err.Error())
			return err
		}
		if len(record) < 2 {
			log.Warn("invalid record")
			continue
		}
		pid, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		nid, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		nodes = append(nodes, node.AddRequest{
			Pid: pid,
			Id:  nid,
		})
		if len(nodes) >= 1000 {
			err := node.BatchAdd(c, dbname, nodes)
			if err != nil {
				log.Error(err.Error())
			}
			nodes = []node.AddRequest{}
		}
	}
	if len(nodes) > 0 {
		err := node.BatchAdd(c, dbname, nodes)
		if err != nil {
			log.Error(err.Error())
		}
		nodes = []node.AddRequest{}
	}
	return nil
}
