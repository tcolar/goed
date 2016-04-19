package assert

import (
	"testing"

	gc "gopkg.in/check.v1"
)

var _ = gc.Suite(&TestingTestSuite{})

// The state for the test suite
type TestingTestSuite struct {
}

func Test(t *testing.T) { gc.TestingT(t) }

func (s *TestingTestSuite) Test(c *gc.C) {
	True(c, 5 == 3+2)
	False(c, 5 == 3)
	Nil(c, nil)
	NotNil(c, 5)
	Eq(c, "foo", "foo")
	NotEq(c, "foo", "bar")
	o1 := []string{"a", "b"}
	o2 := []string{"a", "b"}
	DeepEq(c, o1, o2)
	o2 = append(o2, "c")
	DeepNotEq(c, o1, o2)
	IsType(c, 5, 3)
	NotIsType(c, 5, "a")
}
