package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRemoteFilesInfo_1(t *testing.T) {
	json := "[{\"Status\":\"complete\",\"FileName\":\"test\"}, {\"Status\":\"complete\",\"FileName\":\"test2\"}]"
	rt := getRemoteFilesInfo([]byte(json))

	assert.Equal(t, "test", rt[0].FileName)
	assert.Equal(t, "complete", rt[0].Status)
	assert.Equal(t, "test2", rt[1].FileName)
	assert.Equal(t, "complete", rt[1].Status)
}
