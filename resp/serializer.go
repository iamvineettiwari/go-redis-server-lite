package resp

import (
	"errors"
	"fmt"
	"strconv"
)

func serializeSimpleString(input string) ([]byte, error) {
	data := []byte{}
	data = append(data, SIMPLE_STRING_PREFIX)
	data = append(data, []byte(input)...)
	data = append(data, breakPoint...)
	return data, nil
}

func serializeBulkString(input any) ([]byte, error) {
	data := []byte{}

	data = append(data, BULK_STRING_PREFIX)

	if input == nil {
		data = append(data, []byte("-1")...)
		data = append(data, breakPoint...)
		return data, nil
	}

	inputString := input.(string)
	length := strconv.Itoa(len(inputString))

	data = append(data, []byte(length)...)
	data = append(data, breakPoint...)
	data = append(data, []byte(inputString)...)
	data = append(data, breakPoint...)

	return data, nil
}

func serializeError(input string) ([]byte, error) {
	data := []byte{}
	data = append(data, ERROR_PREFIX)
	data = append(data, []byte(input)...)
	data = append(data, breakPoint...)
	return data, nil
}

func serializeInteger(input int) ([]byte, error) {
	data := []byte{}
	data = append(data, INTEGER_PREFIX)
	data = append(data, []byte(strconv.Itoa(input))...)
	data = append(data, breakPoint...)
	return data, nil
}

func serializeArray(input any) ([]byte, error) {
	data := []byte{}

	data = append(data, ARRAY_PREFIX)

	if input == nil {
		data = append(data, []byte("-1")...)
		data = append(data, breakPoint...)
		return data, nil
	}

	elements := input.([]ArrayType)
	length := strconv.Itoa(len(elements))
	data = append(data, []byte(length)...)
	data = append(data, breakPoint...)

	for _, item := range elements {
		switch item.Type {
		case SIMPLE_STRING:
			serializedData, err := serializeSimpleString(item.Value.(string))

			if err != nil {
				return []byte{}, err
			}

			data = append(data, serializedData...)
		case BULK_STRING:
			serializedData, err := serializeBulkString(item.Value)

			if err != nil {
				return []byte{}, err
			}

			data = append(data, serializedData...)
		case ERROR:
			serializedData, err := serializeError(item.Value.(string))

			if err != nil {
				return []byte{}, err
			}

			data = append(data, serializedData...)
		case INTEGER:
			number, err := getIntValue(item.Value)

			if err != nil {
				return []byte{}, err
			}

			serializedData, err := serializeInteger(number)

			if err != nil {
				return []byte{}, err
			}

			data = append(data, serializedData...)

		case ARRAY:
			serializedData, err := serializeArray(item.Value)

			if err != nil {
				return []byte{}, err
			}

			data = append(data, serializedData...)
		}
	}

	return data, nil
}

func getIntValue(data any) (int, error) {
	num, isNumber := data.(int)

	if isNumber {
		return num, nil
	}

	str, isString := data.(string)

	if isString {
		num, err := strconv.Atoi(str)

		if err != nil {
			return 0, err
		}

		return num, nil
	}

	return 0, errors.New(fmt.Sprintf("Can not convert %v to integer", data))
}

func Serialize(dataType string, data any) ([]byte, error) {
	switch dataType {
	case SIMPLE_STRING:
		str, isString := data.(string)

		if !isString {
			return []byte{}, errors.New("Can not convert data to string")
		}

		return serializeSimpleString(str)

	case BULK_STRING:
		return serializeBulkString(data)

	case ERROR:
		str, isString := data.(string)

		if !isString {
			return []byte{}, errors.New("Can not convert data to error")
		}

		return serializeError(str)

	case INTEGER:
		number, err := getIntValue(data)

		if err != nil {
			return []byte{}, err
		}
		return serializeInteger(number)
	case ARRAY:
		return serializeArray(data)
	}

	return []byte{}, nil
}
