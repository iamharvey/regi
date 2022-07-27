package data

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/ulfox/dby/db"
	"os"
)

const (
	keyCurrent          = "current"
	keyRegistries       = "registries"
	keyRegistryName     = "name"
	keyRegistryServer   = "server"
	keyRegistrySkip     = "insecure-skip-tls-verify"
	keyRegistryUser     = "user"
	keyRegistryPassword = "password"
)

// DB defines a YAML file base data storage.
type DB struct {
	*db.Storage
}

// NewDB returns a new YAML data storage.
func NewDB() (*DB, error) {
	db, err := db.NewStorageFactory(defaultLocation())
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// Registry defines a registry entry.
type Registry struct {
	Name                  string
	Server                string
	InsecureSkipTLSVerify bool
	User                  string
	Password              string
}

// Current returns the current registry setting.
func (db *DB) Current() (*Registry, error) {
	// Get current context.
	keyPath, err := db.GetPath(keyCurrent)
	if err != nil {
		return nil, err
	}
	current := keyPath.(string)

	// Get all the contexts.
	keyPath, err = db.GetPath(keyRegistries)
	if err != nil {
		return nil, err
	}

	if keyPath == nil {
		return nil, errors.New(
			"registries is not specified, there must be something wrong about setting current context")
	}
	registries := keyPath.([]interface{})

	for _, v := range registries {
		r := v.(map[interface{}]interface{})
		if r["name"] != nil && r["name"].(string) == current {
			return packUp(r)
		}
	}

	return nil, errors.New("unable to find context, there must be something wrong about adding context")
}

// List returns all the registries.
func (db *DB) List() ([]*Registry, error) {
	var registries []*Registry

	// Get all the contexts.
	keyPath, err := db.GetPath(keyRegistries)
	if err != nil {
		return nil, err
	}

	if keyPath == nil {
		return nil, errors.New(
			"registries is not specified, there must be something wrong about setting current context")
	}
	list := keyPath.([]interface{})

	for _, v := range list {
		r := v.(map[interface{}]interface{})
		reg, err := packUp(r)
		if err != nil {
			return nil, err
		}
		registries = append(registries, reg)
	}

	return registries, nil
}

// Get a registry entry.
func (db *DB) Get(name string) (*Registry, error) {
	// Get all the contexts.
	keyPath, err := db.GetPath(keyRegistries)
	if err != nil {
		return nil, err
	}

	if keyPath == nil {
		return nil, errors.New(
			"registries is not specified, there must be something wrong about setting current context")
	}
	registries := keyPath.([]interface{})

	for _, v := range registries {
		r := v.(map[interface{}]interface{})
		if r["name"] != nil && r["name"].(string) == name {
			return packUp(r)
		}
	}

	return nil, nil
}

// Set current context.
func (db *DB) Set(name string) error {
	return db.Upsert("current", name)
}

// Add new registry to the context list.
func (db *DB) Add(name, server, user, password string, skip bool) (bool, error) {
	keyPath, err := db.GetPath(keyRegistries)

	var registries []interface{}
	if keyPath != nil {
		registries = keyPath.([]interface{})
	}

	found := false
	for _, v := range registries {
		r := v.(map[interface{}]interface{})
		if r["name"] == name {
			found = true
			break
		}
	}

	if !found {
		registries = append(registries, map[string]interface{}{
			keyRegistryName:     name,
			keyRegistryServer:   server,
			keyRegistrySkip:     skip,
			keyRegistryUser:     user,
			keyRegistryPassword: password,
		})

		err = db.Upsert(keyRegistries, registries)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

// packUp converts a map format registry to struct Registry.
func packUp(reg map[interface{}]interface{}) (*Registry, error) {
	r := Registry{}
	name := reg["name"]
	server := reg["server"]
	skip := reg["insecure-skip-tls-verify"]
	user := reg["user"]
	pass := reg["password"]

	if name == nil {
		return nil, errors.New("name is not specified")
	}
	r.Name = name.(string)

	if server == nil {
		return nil, errors.New("server is not specified")
	}
	r.Server = server.(string)

	if skip != nil {
		r.InsecureSkipTLSVerify = skip.(bool)
	}

	if user != nil {
		r.User = user.(string)
	}

	if pass != nil {
		r.Password = pass.(string)
	}

	return &r, nil
}

func defaultLocation() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	return fmt.Sprintf("%s/.regi/regi.yaml", home)
}
