package models

type Version struct {
	No   string
	Name string
}

func (v Version) Id() string {
	return v.No + "-" + v.Name
}

func NewVersion(no string, name string) *Version {
	return &Version{
		no,
		name,
	}
}
