package utils

import (
	"github.com/sirupsen/logrus"
)

func DoOrDie(err error) {
	if err != nil {
		logrus.Fatalf("%+v", err)
	}
}

func CopySlice[A any](s []A) []A {
	newCopy := make([]A, len(s))
	copy(newCopy, s)
	return newCopy
}
