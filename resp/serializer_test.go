package resp

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleStringSerialization(t *testing.T) {
	str := "hello"

	data, err := Serialize(SIMPLE_STRING, str)

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, string(data), "+hello\r\n")
}

func TestBulkStringSerialization(t *testing.T) {
	str := "Vineet Tiwari"

	data, err := Serialize(BULK_STRING, str)

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, string(data), "$13\r\nVineet Tiwari\r\n")
}

func TestIntegerSerialization(t *testing.T) {
	str := -100

	data, err := Serialize(INTEGER, str)

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, string(data), ":-100\r\n")
}

func TestErrorSerialization(t *testing.T) {
	str := "Something went wrong"

	data, err := Serialize(ERROR, str)

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, string(data), "-Something went wrong\r\n")
}

func TestArraySerialization(t *testing.T) {
	data := []ArrayType{
		{Value: "Hello", Type: SIMPLE_STRING},
		{Value: "GET KEYS", Type: BULK_STRING},
		{Value: []ArrayType{
			{Value: 1, Type: INTEGER},
		}, Type: ARRAY},
	}

	response, err := Serialize(ARRAY, data)

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, string(response), "*3\r\n+Hello\r\n$8\r\nGET KEYS\r\n*1\r\n:1\r\n")
}
