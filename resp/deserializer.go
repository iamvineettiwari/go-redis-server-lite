package resp

import (
	"bytes"
	"errors"
	"strconv"
)

func decodeSimpleString(data []byte) (any, error, int) {
	splitedData := bytes.Split(data, breakPoint)
	value := string(splitedData[0])

	return value, nil, len(splitedData[0]) + len(breakPoint)
}

func decodeBulkString(data []byte) (any, error, int) {
	splitedData := bytes.Split(data, breakPoint)

	if string(splitedData[0]) == "-1" {
		return nil, nil, len(splitedData[0]) + len(breakPoint)
	}

	return string(splitedData[1]), nil, len(splitedData[1]) + len(splitedData[0]) + (len(breakPoint) * 2)
}

func decodeInteger(data []byte) (any, error, int) {
	splitedData := bytes.Split(data, breakPoint)

	value, err := strconv.Atoi(string(splitedData[0]))

	if err != nil {
		return nil, err, 0
	}

	return value, nil, len(splitedData[0]) + len(breakPoint)
}

func decodeError(data []byte) (any, error, int) {
	splitedData := bytes.Split(data, breakPoint)

	value := string(splitedData[0])

	return value, nil, len(splitedData[0]) + len(breakPoint)
}

func decodeArray(data []byte) (any, error, int) {
	firstSplit := bytes.SplitN(data, breakPoint, 2)

	length, err := strconv.Atoi(string(firstSplit[0]))

	if err != nil {
		return nil, err, 0
	}

	if length < 0 {
		return nil, nil, 0
	}

	values := []ArrayType{}

	body := firstSplit[1]
	totalRead := 0
	startPositon := 0

	for i := 0; i < length; i++ {
		elements, dataType, err, readLength := Deserialize(body[startPositon:])

		if err != nil {
			return nil, err, 0
		}

		totalRead += readLength

		values = append(values, ArrayType{
			Value: elements,
			Type:  dataType,
		})

		startPositon += readLength
	}

	return values, nil, totalRead + len(firstSplit[0]) + len(breakPoint)
}

func Deserialize(data []byte) (any, string, error, int) {
	prefix := data[0]

	var decodedData any
	var err error
	var readLength int
	var dataType string

	switch prefix {
	case SIMPLE_STRING_PREFIX:
		dataType = SIMPLE_STRING
		decodedData, err, readLength = decodeSimpleString(data[1:])

		if err != nil {
			decodedData = nil
			readLength = -1
		}

	case BULK_STRING_PREFIX:
		dataType = BULK_STRING
		decodedData, err, readLength = decodeBulkString(data[1:])

		if err != nil {
			decodedData = nil
			readLength = -1
		}

	case INTEGER_PREFIX:
		dataType = INTEGER
		decodedData, err, readLength = decodeInteger(data[1:])

		if err != nil {
			decodedData = nil
			readLength = -1
		}

	case ERROR_PREFIX:
		dataType = ERROR
		decodedData, err, readLength = decodeError(data[1:])

		if err != nil {
			decodedData = nil
			readLength = -1
		}

	case ARRAY_PREFIX:
		dataType = ARRAY
		decodedData, err, readLength = decodeArray(data[1:])

		if err != nil {
			decodedData = nil
			readLength = -1
		}

	default:
		dataType = UNSUPORTED_TYPE
		decodedData = nil
		err = errors.New("Data type is not supported")
		readLength = -1
	}

	return decodedData, dataType, err, readLength + 1
}
