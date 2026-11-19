package builtins

import (
	"fmt"

	"github.com/ep0ll/nixlang-go/value"
)

// Constant identifies a non-callable value exposed by the Nix builtins set.
type Constant uint8

const (
	ConstantInvalid Constant = iota
	ConstantTrue
	ConstantFalse
	ConstantNull
	ConstantCurrentTime
	ConstantCurrentSystem
	ConstantNixVersion
	ConstantStoreDir
	ConstantLangVersion
	ConstantNixPath
)

func (c Constant) String() string {
	switch c {
	case ConstantTrue:
		return "true"
	case ConstantFalse:
		return "false"
	case ConstantNull:
		return "null"
	case ConstantCurrentTime:
		return "currentTime"
	case ConstantCurrentSystem:
		return "currentSystem"
	case ConstantNixVersion:
		return "nixVersion"
	case ConstantStoreDir:
		return "storeDir"
	case ConstantLangVersion:
		return "langVersion"
	case ConstantNixPath:
		return "nixPath"
	default:
		return ""
	}
}

// GetConstant returns an owned Nix value for a known non-callable builtin
// constant. The caller must close the returned value.
func (b *Builtins) GetConstant(constant Constant) (*value.Value, error) {
	name := constant.String()
	if name == "" {
		return nil, fmt.Errorf("nix: invalid builtin constant %d", constant)
	}
	return b.Get(name)
}

// True returns the value of builtins.true.
func (b *Builtins) True() (bool, error) {
	v, err := b.GetConstant(ConstantTrue)
	if err != nil {
		return false, err
	}
	defer v.Close()
	return v.Bool()
}

// False returns the value of builtins.false.
func (b *Builtins) False() (bool, error) {
	v, err := b.GetConstant(ConstantFalse)
	if err != nil {
		return false, err
	}
	defer v.Close()
	return v.Bool()
}

// Null returns an owned value for builtins.null. The caller must close it.
func (b *Builtins) Null() (*value.Value, error) {
	return b.GetConstant(ConstantNull)
}

// CurrentTime returns builtins.currentTime as Unix seconds. It is unavailable
// when the evaluator is configured for pure evaluation.
func (b *Builtins) CurrentTime() (int64, error) {
	v, err := b.GetConstant(ConstantCurrentTime)
	if err != nil {
		return 0, err
	}
	defer v.Close()
	return v.Int()
}

// CurrentSystem returns builtins.currentSystem. It is unavailable when the
// evaluator is configured for pure evaluation.
func (b *Builtins) CurrentSystem() (string, error) {
	v, err := b.GetConstant(ConstantCurrentSystem)
	if err != nil {
		return "", err
	}
	defer v.Close()
	return v.String()
}

// NixVersion returns builtins.nixVersion.
func (b *Builtins) NixVersion() (string, error) {
	v, err := b.GetConstant(ConstantNixVersion)
	if err != nil {
		return "", err
	}
	defer v.Close()
	return v.String()
}

// StoreDir returns builtins.storeDir.
func (b *Builtins) StoreDir() (string, error) {
	v, err := b.GetConstant(ConstantStoreDir)
	if err != nil {
		return "", err
	}
	defer v.Close()
	return v.String()
}

// LangVersion returns builtins.langVersion.
func (b *Builtins) LangVersion() (int64, error) {
	v, err := b.GetConstant(ConstantLangVersion)
	if err != nil {
		return 0, err
	}
	defer v.Close()
	return v.Int()
}

// NixPath returns an owned value for builtins.nixPath. The caller must close it.
func (b *Builtins) NixPath() (*value.Value, error) {
	return b.GetConstant(ConstantNixPath)
}
