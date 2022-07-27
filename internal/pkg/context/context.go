package context

/* ------ Sample config file ------

registries:
  - name: jc-registry
    server: http://192.168.0.168:5000
    insecure-skip-tls-verify: false
    user:
    password:
  - name: jc-registry-2
    server: https://192.168.0.168:5000
    insecure-skip-tls-verify: true
    user: noname
    password: fakepass
current: jc-registry

 ------ */

import (
	"github.com/ulfox/dby/db"
)

// NewStore returns a new YAML data store.
func NewStore(filepath string) (*db.Storage, error) {
	return db.NewStorageFactory(filepath)
}
