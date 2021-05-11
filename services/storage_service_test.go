package services

import (
	"fmt"
	"path"
	"testing"
)

func TestParent(t *testing.T) {
	p := "/home/zhanlan/projects/goProjects"

	fmt.Println(path.Dir(p))
	fmt.Println(path.Base(p))
	fmt.Println(path.Clean(p))
	fmt.Println(path.Split(p))

	fmt.Println(path.Clean("/"))
}
