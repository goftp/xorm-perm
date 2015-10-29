package xormperm

import (
	"os"
	"time"

	"github.com/go-xorm/xorm"
)

type Group struct {
	Id      int64
	Name    string    `xorm:"unique"`
	Created time.Time `xorm:"created"`
}

type UserGroup struct {
	UserName  string `xorm:"pk"`
	GroupName string `xorm:"pk"`
}

type Perm struct {
	Id      int64
	Path    string `xorm:"unique"`
	Owner   string
	Group   string
	Mode    os.FileMode
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

type XormPerm struct {
	engine       *xorm.Engine
	defaultOwner string
	defaultGroup string
	defaultMode  os.FileMode
}

func New(engine *xorm.Engine, owner, group string, mode os.FileMode) *XormPerm {
	return &XormPerm{engine, owner, group, mode}
}

func (db *XormPerm) GetOwner(rPath string) (string, error) {
	var perm = Perm{Path: rPath}
	has, err := db.engine.Get(&perm)
	if err != nil {
		return "", err
	}
	if !has {
		return db.defaultOwner, nil
	}

	return perm.Owner, nil
}

func (db *XormPerm) GetGroup(rPath string) (string, error) {
	var perm = Perm{Path: rPath}
	has, err := db.engine.Get(&perm)
	if err != nil {
		return "", err
	}
	if !has {
		return db.defaultGroup, nil
	}

	return perm.Group, nil
}

func (db *XormPerm) GetMode(rPath string) (os.FileMode, error) {
	var perm = Perm{Path: rPath}
	has, err := db.engine.Get(&perm)
	if err != nil {
		return os.ModeType, err
	}
	if !has {
		return db.defaultMode, nil
	}

	return perm.Mode, nil
}

func (db *XormPerm) ChOwner(rPath, owner string) error {
	var perm = Perm{Path: rPath}
	has, err := db.engine.Get(&perm)
	if err != nil {
		return err
	}
	if !has {
		perm.Owner = owner
		perm.Group = db.defaultGroup
		perm.Mode = db.defaultMode
		_, err = db.engine.Insert(&perm)
		return err
	}

	_, err = db.engine.Update(&Perm{Owner: owner}, &perm)
	return err
}

func (db *XormPerm) ChGroup(rPath, group string) error {
	var perm = Perm{Path: rPath}
	has, err := db.engine.Get(&perm)
	if err != nil {
		return err
	}
	if !has {
		perm.Owner = db.defaultOwner
		perm.Group = group
		perm.Mode = db.defaultMode
		_, err = db.engine.Insert(&perm)
		return err
	}

	_, err = db.engine.Update(&Perm{Group: group}, &perm)
	return err
}

func (db *XormPerm) ChMode(rPath string, mode os.FileMode) error {
	var perm = Perm{Path: rPath}
	has, err := db.engine.Get(&perm)
	if err != nil {
		return err
	}
	if !has {
		perm.Owner = db.defaultOwner
		perm.Group = db.defaultGroup
		perm.Mode = mode
		_, err = db.engine.Insert(&perm)
		return err
	}

	_, err = db.engine.Update(&Perm{Mode: mode}, &perm)
	return err
}
