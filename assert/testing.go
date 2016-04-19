package assert

import (
	"reflect"

	gc "gopkg.in/check.v1"
)

func assert(c *gc.C, checker gc.Checker, expected interface{}, obtained interface{}) {
	c.Assert(obtained, checker, expected)
}

func True(c *gc.C, b bool) {
	assert(c, gc.Equals, true, b)
}

func False(c *gc.C, b bool) {
	assert(c, gc.Equals, false, b)
}

func Nil(c *gc.C, obtained interface{}) {
	c.Assert(obtained, gc.IsNil)
}

func NotNil(c *gc.C, obtained interface{}) {
	c.Assert(obtained, gc.Not(gc.IsNil))
}

func Eq(c *gc.C, obtained, expected interface{}) {
	assert(c, gc.Equals, expected, obtained)
}

func NotEq(c *gc.C, obtained, expected interface{}) {
	assert(c, gc.Not(gc.Equals), expected, obtained)
}

func DeepEq(c *gc.C, obtained, expected interface{}) {
	assert(c, gc.DeepEquals, expected, obtained)
}

func DeepNotEq(c *gc.C, obtained, expected interface{}) {
	assert(c, gc.Not(gc.DeepEquals), expected, obtained)
}

func IsType(c *gc.C, obtained, expected interface{}) {
	Eq(c, reflect.TypeOf(expected), reflect.TypeOf(obtained))
}

func NotIsType(c *gc.C, obtained, expected interface{}) {
	NotEq(c, reflect.TypeOf(expected), reflect.TypeOf(obtained))
}
