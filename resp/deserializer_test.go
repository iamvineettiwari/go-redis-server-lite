package resp

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleString(t *testing.T) {
	str := "+PING\r\n"

	data, dataType, err, _ := Deserialize([]byte(str))

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, data, "PING")
	assert.Equal(t, dataType, SIMPLE_STRING)

	str = "+hello world\r\n"

	data, dataType, err, _ = Deserialize([]byte(str))

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, data, "hello world")
	assert.Equal(t, dataType, SIMPLE_STRING)
}

func TestBulkString(t *testing.T) {
	str := "$0\r\n\r\n"

	data, dataType, err, _ := Deserialize([]byte(str))

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, data, "")
	assert.Equal(t, dataType, BULK_STRING)

	str = "$13\r\nVineet Tiwari\r\n"

	data, dataType, err, _ = Deserialize([]byte(str))

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, data, "Vineet Tiwari")
	assert.Equal(t, dataType, BULK_STRING)
}

func TestBulkStringNull(t *testing.T) {
	str := "$-1\r\n"

	data, dataType, err, _ := Deserialize([]byte(str))

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, data, nil)
	assert.Equal(t, dataType, BULK_STRING)
}

func TestInteger(t *testing.T) {
	str := ":-100\r\n"

	data, dataType, err, _ := Deserialize([]byte(str))

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, data, -100)
	assert.Equal(t, dataType, INTEGER)
}

func TestError(t *testing.T) {
	str := "-Error message\r\n"

	data, dataType, err, _ := Deserialize([]byte(str))

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, data, "Error message")
	assert.Equal(t, dataType, ERROR)
}

func TestArray(t *testing.T) {
	str := "*1\r\n$4\r\nping\r\n"

	data, dataType, err, _ := Deserialize([]byte(str))

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, data, []ArrayType{
		{Value: "ping", Type: BULK_STRING},
	})
	assert.Equal(t, dataType, ARRAY)

	str = "*2\r\n$4\r\necho\r\n$11\r\nhello world\r\n"

	data, dataType, err, _ = Deserialize([]byte(str))

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, data, []ArrayType{
		{Value: "echo", Type: BULK_STRING},
		{Value: "hello world", Type: BULK_STRING},
	})
	assert.Equal(t, dataType, ARRAY)
}
