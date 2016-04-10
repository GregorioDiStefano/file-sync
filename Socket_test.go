package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCompletePayload_1(t *testing.T) {
	payloadPrefix := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04}
	payload := append(payloadPrefix, []byte{0xff}...)
	payload = append(payload, []byte{0xff}...)
	payload = append(payload, []byte{0xff}...)
	payload = append(payload, []byte{0xff}...)

	assert.True(t, getCompletePayload(payload))
}

func TestGetCompletePayload_2(t *testing.T) {
	payloadPrefix := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04}
	payload := append(payloadPrefix, []byte{0xff}...)
	payload = append(payload, []byte{0xff}...)
	payload = append(payload, []byte{0xff}...)

	assert.False(t, getCompletePayload(payload))
}

func TestGetCompletePayload_3(t *testing.T) {
	payloadPrefix := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x02, 0xFF, 0xFF}
	var payload []byte
	for i := 0; i < 33751039; i++ {
		payload = append(payload, []byte{0xff}...)
	}

	payload = append(payloadPrefix, payload...)

	assert.True(t, getCompletePayload(payload))
}

func TestGetCompletePayload_4(t *testing.T) {
	payloadPrefix := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xFF, 0xFF, 0xFF}
	var payload []byte
	for i := 0; i < 50331647; i++ {
		payload = append(payload, []byte{0xff}...)
	}

	payload = append(payloadPrefix, payload...)

	assert.True(t, getCompletePayload(payload))
}

func TestGetCompletePayload_5(t *testing.T) {
	payloadPrefix := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xFF, 0xFF, 0xFF}
	var payload []byte
	for i := 0; i < 50331646; i++ {
		payload = append(payload, []byte{0xff}...)
	}

	payload = append(payloadPrefix, payload...)

	assert.False(t, getCompletePayload(payload))
}
