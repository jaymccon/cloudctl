package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jaymccon/cloudctl/providers/aws"
	bolt "go.etcd.io/bbolt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// TODO: cross platform paths
	cacheDir          = "~/.cloudctl/cache/"
	cacheFilename     = "bbolt.db"
	CacheROMode       = "RO"
	CacheRWMode       = "RW"
	defaultBucketName = "cloudctlDefault"
)

type Cache struct {
	Mode  string
	cache *bolt.DB
}

func NewCache(mode string) (*Cache, error) {
	cachePath, err := absPath(cacheDir)
	if err != nil {
		fmt.Printf("ERROR: finding cache dir: %q\n", err.Error())
		return nil, err
	}
	if mode == CacheRWMode {
		err = mkCacheDir(cachePath)
		if err != nil {
			fmt.Printf("ERROR: creating cache dir: %q\n", err.Error())
			return nil, err
		}
	}
	db, err := bolt.Open(*cachePath+cacheFilename, 0666, &bolt.Options{ReadOnly: isRO(mode), Timeout: 1 * time.Second})
	if err != nil {
		fmt.Printf("ERROR: opening bolt cache: %q\n", err.Error())
		return nil, err
	}
	c := Cache{
		Mode:  mode,
		cache: db,
	}
	if mode == CacheRWMode {
		err = c.createBucket()
		if err != nil {
			fmt.Printf("ERROR: creating bolt cache bucket: %q\n", err.Error())
			return nil, err
		}
	}
	return &c, nil
}

func UpdateCache() error {
	// TODO: create lock to prevent concurrent upgrade operations from racing
	fmt.Println("Downloading schema files...")
	c, err := NewCache(CacheRWMode)
	if err != nil {
		fmt.Printf("ERROR: opening cache: %q\n", err.Error())
		return err
	}
	defer func(cache *bolt.DB) {
		err := cache.Close()
		if err != nil {
			fmt.Printf("ERROR: closing cache error: %s\n", err.Error())
		}
	}(c.cache)

	schemas, err := aws.FetchSchemas()
	if err != nil {
		fmt.Printf("ERROR: fetching schemas: %q\n", err.Error())
		return err
	}
	fmt.Println("Saving schema json files to disk...")
	schemaPath, err := absPath(schemaDir)
	if err != nil {
		return err
	}
	err = mkCacheDir(schemaPath)
	if err != nil {
		return err
	}
	var schemaList []string
	for name, schema := range *schemas {
		schemaList = append(schemaList, name)
		fname := *schemaPath + strings.Replace(name, "::", "_", -1) + ".json"
		err = ioutil.WriteFile(fname, schema, 0644)
		if err != nil {
			return err
		}
	}
	fmt.Println("Parsing schema files...")
	parsedSchemas, err := ParseSchemas()
	if err != nil {
		return err
	}
	fmt.Println("Resetting cache...")
	err = c.deleteBucket()
	if err != nil {
		return err
	}
	err = c.createBucket()
	if err != nil {
		return err
	}
	err = c.PutSchemas(*parsedSchemas)
	if err != nil {
		return err
	}
	err = c.PutLists(map[string][]string{"__schemas__": schemaList})
	if err != nil {
		return err
	}
	return nil
}

func (c Cache) createBucket() error {
	err := c.cache.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultBucketName))
		if b == nil {
			_, err := tx.CreateBucket([]byte(defaultBucketName))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("create bucket: %s", err)
	}
	return nil
}

func (c Cache) deleteBucket() error {
	err := c.cache.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultBucketName))
		if b == nil {
			return nil
		}
		return tx.DeleteBucket([]byte(defaultBucketName))
	})
	if err != nil {
		return fmt.Errorf("delete bucket: %s", err)
	}
	return nil
}

func (c Cache) Put(cacheMap map[string]string) error {
	err := c.cache.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultBucketName))
		for k, v := range cacheMap {
			err := b.Put([]byte(k), []byte(v))
			if err != nil {
				return fmt.Errorf("writing to cache: %s", err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("writing to cache: %s", err)
	}
	return nil
}

func (c Cache) Get(key string) (*string, error) {
	var value *string
	err := c.cache.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultBucketName))
		v := b.Get([]byte(key))
		if v == nil {
			return errors.New(fmt.Sprintf("failed to fetch %q from cache, key does not exist", key))
		}
		strV := string(v)
		value = &strV
		return nil
	})
	return value, err
}

func (c Cache) GetList(key string) (*[]string, error) {
	var listValue []string
	value, err := c.Get(key)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(*value), &listValue)
	return &listValue, err
}

func (c Cache) PutLists(cacheMap map[string][]string) error {
	marshalledMap := map[string]string{}
	for k, v := range cacheMap {
		jsonV, err := json.Marshal(v)
		if err != nil {
			return err
		}
		marshalledMap[k] = string(jsonV)
	}
	err := c.Put(marshalledMap)
	return err
}

func (c Cache) GetSchema(key string) (*CfnSchema, error) {
	var schema CfnSchema
	value, err := c.Get(key)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(*value), &schema)
	return &schema, err
}

func (c Cache) PutSchemas(cacheMap map[string]CfnSchema) error {
	marshalledMap := map[string]string{}
	for k, v := range cacheMap {
		jsonV, err := v.ToJsonString()
		if err != nil {
			return err
		}
		marshalledMap[k] = *jsonV
	}
	err := c.Put(marshalledMap)
	return err
}

func absPath(path string) (*string, error) {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = homeDir + strings.TrimLeft(path, "~")
	}
	path, err := filepath.Abs(path)
	path = path + "/"
	if err != nil {
		return nil, err
	}
	return &path, nil
}

func isRO(mode string) bool {
	if mode == CacheROMode {
		return true
	}
	return false
}

func mkCacheDir(cachePath *string) error {
	if _, err := os.Stat(*cachePath); os.IsNotExist(err) {
		err := os.MkdirAll(*cachePath, 0755)
		if err != nil {
			fmt.Printf("ERROR: creating cache directory: %q\n", err.Error())
			return err
		}
	}
	return nil
}

func GetSchemas() (*map[string]map[string]map[string]CfnSchema, error) {
	c, err := NewCache(CacheROMode)
	if err != nil {
		return nil, err
	}
	schemaList, err := c.GetList("__schemas__")
	if err != nil {
		return nil, err
	}
	schemas := map[string]map[string]map[string]CfnSchema{}
	for _, name := range *schemaList {
		schema, err := c.GetSchema(name)
		if err != nil {
			fmt.Printf("ERROR: failed to get schema for %q: %q\n", name, err.Error())
		}
		provider, service, resource, err := splitName(name)
		if err != nil {
			return nil, err
		}
		if schema != nil {
			if !keyExists(schemas, *provider) {
				schemas[*provider] = map[string]map[string]CfnSchema{}
			}
			if !keyExists(schemas[*provider], *service) {
				schemas[*provider][*service] = map[string]CfnSchema{}
			}
			schemas[*provider][*service][*resource] = *schema
		}
	}
	return &schemas, nil
}

func GetSchema(typeName string) (*CfnSchema, error) {
	c, err := NewCache(CacheROMode)
	if err != nil {
		return nil, err
	}
	schema, err := c.GetSchema(typeName)
	if err != nil {
		fmt.Printf("ERROR: failed to get schema for %q: %q\n", typeName, err.Error())
	}
	return schema, err
}

func splitName(name string) (*string, *string, *string, error) {
	sep := "_"
	if strings.Contains(name, "::") {
		sep = "::"
	}
	splitName := strings.Split(name, sep)
	if len(splitName) != 3 {
		return nil, nil, nil, errors.New(fmt.Sprintf("resource names should split into exactly 3 elements, %q separated by %q splits into %q", name, sep, len(splitName)))
	}
	return &splitName[0], &splitName[1], &splitName[2], nil
}

func keyExists(inputMap interface{}, key string) bool {
	var m1 map[string]map[string]map[string]CfnSchema
	var m2 map[string]map[string]CfnSchema
	switch g := inputMap.(type) {
	case map[string]map[string]map[string]CfnSchema:
		m1 = inputMap.(map[string]map[string]map[string]CfnSchema)
	case map[string]map[string]CfnSchema:
		m2 = inputMap.(map[string]map[string]CfnSchema)
	default:
		fmt.Printf("ERROR cannot check keys for %q %q\n", key, g)
	}
	if m1 != nil {
		for k, _ := range m1 {
			if k == key {
				return true
			}
		}
	}
	if m2 != nil {
		for k, _ := range m2 {
			if k == key {
				return true
			}
		}
	}
	return false
}

func IsUpdatable(input interface{}) bool {
	return checkSchema(input, "update")
}

func IsConfigurable(input interface{}) bool {
	return checkSchema(input, "configure")
}

func checkSchema(input interface{}, checkType string) bool {
	switch input.(type) {
	case CfnSchema:
		schema := input.(CfnSchema)
		if checkType == "update" {
			return schema.IsUpdatable()
		} else if checkType == "configure" {
			return schema.IsConfigurable()
		}
	case map[string]CfnSchema:
		for _, i := range input.(map[string]CfnSchema) {
			if checkType == "update" {
				if i.IsUpdatable() {
					return true
				}
			} else if checkType == "configure" {
				if i.IsConfigurable() {
					return true
				}
			}
		}
	case map[string]map[string]CfnSchema:
		for _, s := range input.(map[string]map[string]CfnSchema) {
			for _, i := range s {
				if checkType == "update" {
					if i.IsUpdatable() {
						return true
					}
				} else if checkType == "configure" {
					if i.IsConfigurable() {
						return true
					}
				}
			}
		}
	}
	return false
}
