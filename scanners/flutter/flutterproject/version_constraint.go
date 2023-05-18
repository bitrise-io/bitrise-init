package flutterproject

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

type VersionConstraintSource string

const (
	PubspecLockVersionSource VersionConstraintSource = "pubspec.lock"
	PubspecVersionSource     VersionConstraintSource = "pubspec.yaml"
	FVMConfigVersionSource   VersionConstraintSource = "fvm_config.json"
)

/*
VersionConstraint stores either an exact version or a version constraint.
Version is a valid semantic version, constraint supports the Caret and the traditional syntax.

Caret syntax
- ^1.2.3 = >=1.2.3 <2.0.0
- ^0.1.2 = >=0.1.2 <0.2.0 (prior to a 1.0.0 release the minor versions acts as the API stability level)

Traditional syntax
- any (any version)
- 1.2.3
- >=1.2.3
- >1.2.3
- <=1.2.3
- <1.2.3
*/
type VersionConstraint struct {
	Version    *semver.Version
	Constraint *semver.Constraints
	Source     VersionConstraintSource
}

func NewVersionConstraint(constraint string, source VersionConstraintSource) (*VersionConstraint, error) {
	var v *semver.Version
	var c *semver.Constraints
	var err error

	v, err = semver.NewVersion(constraint)
	if err != nil {
		c, err = semver.NewConstraint(constraint)
		if err != nil {
			return nil, err
		}
	}

	return &VersionConstraint{
		Version:    v,
		Constraint: c,
		Source:     source,
	}, nil
}

func (v VersionConstraint) String() string {
	if v.Version != nil {
		return fmt.Sprintf("%s (%s)", v.Version.String(), v.Source)
	}
	if v.Constraint != nil {
		return fmt.Sprintf("%s (%s)", v.Constraint.String(), v.Source)
	}
	return ""
}
