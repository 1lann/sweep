package sweep

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
)

// Possible errors from ReadDecode. I/O errors can occur as well.
var (
	ErrUnexpectedType = errors.New("sweep: unexpected type")
)

// ReadDecode decodes the packet into the provided pointer to a struct.
func (d *Device) ReadDecode(dst interface{}) error {
	data, err := d.reader.ReadBytes('\n')
	if err != nil {
		return err
	}
	rd := bytes.NewReader(data[:len(data)-1])

	return rawReadDecode(rd, dst)
}

func rawReadDecode(rd io.Reader, dst interface{}) error {
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr ||
		dstValue.Elem().Kind() != reflect.Struct {
		return errors.New("packet: dst must be a pointer to a struct")
	}

	fields := dstValue.Elem().NumField()
	for i := 0; i < fields; i++ {
		field := dstValue.Elem().Field(i)

		switch field.Kind() {
		case reflect.Array:
			if field.Index(0).Kind() != reflect.Uint8 {
				return ErrUnexpectedType
			}

			for i := 0; i < field.Len(); i++ {
				value, err := readByte(rd)
				if err != nil {
					return err
				}
				field.Index(i).Set(reflect.ValueOf(value))
			}

		case reflect.Uint8:
			value, err := readByte(rd)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(value))
		case reflect.Int8:
			value, err := readByte(rd)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(int8(value)))
		case reflect.Uint16:
			value, err := readUint16(rd)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(value))
		case reflect.Int16:
			value, err := readUint16(rd)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(int16(value)))
		default:
			return ErrUnexpectedType
		}
	}

	return nil
}

func readByte(rd io.Reader) (byte, error) {
	var result byte
	err := binary.Read(rd, binary.LittleEndian, &result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func readUint16(rd io.Reader) (uint16, error) {
	var result uint16
	err := binary.Read(rd, binary.LittleEndian, &result)
	if err != nil {
		return 0, err
	}

	return result, nil
}
