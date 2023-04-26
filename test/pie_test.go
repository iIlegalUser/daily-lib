package test

import (
	"fmt"
	"github.com/elliotchance/pie/v2"
	"strings"
	"testing"
)

func TestPie(t *testing.T) {
	name := pie.Of([]string{"Bob", "Sally", "John", "Jane"}).
		FilterNot(func(name string) bool {
			return strings.HasPrefix(name, "J")
		}).
		Map(strings.ToUpper).
		Last()

	fmt.Println(name) // "SALLY"
}

func TestEach(t *testing.T) {
	pie.Of([]string{"Bob", "Sally", "John", "Jane"}).
		Each(func(s string) {
			fmt.Println(s)
		})
}

func TestDiff(t *testing.T) {
	added, removed := pie.Diff([]string{"a", "b", "c"}, []string{"b", "c", "d"})
	fmt.Println(added, removed) // [d] [a]
}

func TestIntersect(t *testing.T) {
	ss2 := pie.Intersect([]string{"a", "b", "c"}, []string{"b", "c", "d"})
	fmt.Println(ss2) // [c b]
}
