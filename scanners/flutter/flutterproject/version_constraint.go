package flutterproject

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

type VersionConstraintSource string

const (
	PubspecLockVersionSource VersionConstraintSource = "pubspec_lock"
	PubspecVersionSource     VersionConstraintSource = "pubspec_yaml"
	FVMConfigVersionSource   VersionConstraintSource = "fvm_config_json"
	ASDFConfigVersionSource  VersionConstraintSource = "tool_versions"
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

func NewVersionConstraint(version string, source VersionConstraintSource) (*VersionConstraint, error) {
	var v *semver.Version
	var c *semver.Constraints

	var vErr error
	v, vErr = semver.NewVersion(version)
	if vErr != nil {
		var cErr error
		c, cErr = semver.NewConstraint(version)
		if cErr != nil {
			return nil, fmt.Errorf("invalid version (%s): not a semantic version (%s) nor a version constraint (%s)", version, vErr, cErr)
		}
	}

	return &VersionConstraint{
		Version:    v,
		Constraint: c,
		Source:     source,
	}, nil
}
