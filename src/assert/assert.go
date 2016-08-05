package assert

import (
	"fmt"

	"github.com/Sirupsen/logrus"
)

// Nil panic if the test is not nil
func Nil(test interface{}, params ...interface{}) {
	if test != nil {
		f := logrus.Fields{}
		for i := range params {
			f[fmt.Sprintf("param%d", i)] = params[i]
		}

		if e, ok := test.(error); ok {
			logrus.WithFields(f).Panic(e)
		}
		logrus.WithFields(f).Panic("must be nil, but its not")
	}
}

// NotNil panic if the test is nil
func NotNil(test interface{}, params ...interface{}) {
	if test == nil {
		f := logrus.Fields{}
		for i := range params {
			f[fmt.Sprintf("param%d", i)] = params[i]
		}
		logrus.WithFields(f).Panic("must not be nil, but it is")
	}
}

// True check if the value is true and panic if its not
func True(test bool, message string, params ...interface{}) {
	if !test {
		f := logrus.Fields{}
		for i := range params {
			f[fmt.Sprintf("param%d", i)] = params[i]
		}
		logrus.WithFields(f).Panic("must be true, but its not")
	}
}

// False check if the test is false
func False(test bool, message string, params ...interface{}) {
	True(!test, message, params...)

}

// Empty check if the string is empty and panic if not
func Empty(test string, message string, params ...interface{}) {
	True(test == "", message, params...)
}
