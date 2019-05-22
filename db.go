package depthtree

import (
	"fmt"
	"github.com/mkideal/log"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

const EXTENSION = ".bin"

type Database struct {
	trees  map[string]*Tree
	dbPath string
	sync.RWMutex
}

func NewDatabase(dbPath string) *Database {
	return &Database{
		trees:  make(map[string]*Tree),
		dbPath: dbPath,
	}
}

func (this *Database) NewTree(name string) *Tree {
	this.RLock()
	if tree, found := this.trees[name]; found {
		this.RUnlock()
		return tree
	}
	this.RUnlock()
	tree := NewTree(name)
	this.Lock()
	this.trees[name] = tree
	this.Unlock()
	return tree
}

func (this *Database) Use(name string) *Tree {
	this.RLock()
	if tree, found := this.trees[name]; found {
		this.RUnlock()
		return tree
	}
	this.RUnlock()
	return nil
}

func (this *Database) List() []string {
	this.RLock()
	var dbs []string
	for k, _ := range this.trees {
		dbs = append(dbs, k)
	}
	this.RUnlock()
	return dbs
}

func (this *Database) Truncate(name string) {
	this.Lock()
	delete(this.trees, name)
	dbPath := this.dbPath
	this.Unlock()
	filePath := path.Join(dbPath, fmt.Sprintf("%s%s", name, EXTENSION))
	os.Remove(filePath)
}

func (this *Database) Flush() {
	var trees []*Tree
	this.RLock()
	dbPath := this.dbPath
	for _, tree := range this.trees {
		trees = append(trees, tree)
	}
	this.RUnlock()
	wg := &sync.WaitGroup{}
	for _, tree := range trees {
		filePath := path.Join(dbPath, fmt.Sprintf("%s%s", tree.Name(), EXTENSION))
		wg.Add(1)
		go func(wg *sync.WaitGroup, tree *Tree, filePath string) {
			log.Info("flushing: %s", filePath)
			err := tree.Flush(filePath)
			if err != nil {
				log.Error(err.Error())
			}
			wg.Done()
		}(wg, tree, filePath)
	}
	wg.Wait()
}

func (this *Database) Open() error {
	dbPath := this.dbPath
	log.Info("opening: %s", dbPath)
	list, err := ioutil.ReadDir(dbPath)
	if err != nil {
		return err
	}
	for _, f := range list {
		if f.IsDir() {
			continue
		}
		filename := f.Name()
		if !strings.HasSuffix(filename, EXTENSION) {
			continue
		}
		log.Info("reading db: %s", filename)
		treeName := strings.TrimSuffix(filename, EXTENSION)
		filePath := path.Join(dbPath, filename)
		tree := NewTree(treeName)
		err := tree.Reload(filePath)
		if err != nil {
			return err
		}
		this.Lock()
		this.trees[treeName] = tree
		this.Unlock()
	}
	log.Info("db ready: %s", dbPath)
	return nil
}
