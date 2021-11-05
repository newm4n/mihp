package helper

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"time"
)

const (
	STypeInt byte = iota
	STypeByte
	STypeInt8
	STypeInt16
	STypeInt32
	STypeInt64
	STypeUint
	STypeUint8
	STypeUint16
	STypeUint32
	STypeUint64
	STypeBool
	STypeString
	STypeDuration
	STypeTime
	STypeStringArray
)

func Read(r io.Reader) (interface{}, error) {
	typ := []byte{0}
	if _, err := r.Read(typ); err != nil {
		return nil, err
	}
	switch typ[0] {
	case STypeBool:
		return ReadBool(r)
	case STypeByte:
		uin, err := ReadUint8(r)
		if err != nil {
			return byte(0), err
		}
		return byte(uin), nil
	case STypeInt:
		return ReadInt(r)
	case STypeInt8:
		return ReadInt8(r)
	case STypeInt16:
		return ReadInt16(r)
	case STypeInt32:
		return ReadInt32(r)
	case STypeInt64:
		return ReadInt64(r)
	case STypeUint:
		return ReadUint(r)
	case STypeUint8:
		return ReadUint8(r)
	case STypeUint16:
		return ReadUint16(r)
	case STypeUint32:
		return ReadUint32(r)
	case STypeUint64:
		return ReadUint64(r)
	case STypeString:
		return ReadString(r)
	case STypeDuration:
		return ReadDuration(r)
	case STypeTime:
		return ReadTime(r)
	case STypeStringArray:
		return ReadStringArray(r)
	default:
		return nil, fmt.Errorf("unsupported context value")
	}
}

func Put(w io.Writer, i interface{}) error {
	switch i.(type) {
	case bool:
		if _, err := w.Write([]byte{STypeBool}); err != nil {
			return err
		}
		return PutBool(w, i.(bool))
	case int:
		if _, err := w.Write([]byte{STypeInt}); err != nil {
			return err
		}
		return PutInt(w, i.(int))
	case int8:
		if _, err := w.Write([]byte{STypeInt8}); err != nil {
			return err
		}
		return PutInt8(w, i.(int8))
	case int16:
		if _, err := w.Write([]byte{STypeInt16}); err != nil {
			return err
		}
		return PutInt16(w, i.(int16))
	case int32:
		if _, err := w.Write([]byte{STypeInt32}); err != nil {
			return err
		}
		return PutInt32(w, i.(int32))
	case int64:
		if _, err := w.Write([]byte{STypeInt64}); err != nil {
			return err
		}
		return PutInt64(w, i.(int64))
	case uint:
		if _, err := w.Write([]byte{STypeUint}); err != nil {
			return err
		}
		return PutUint(w, i.(uint))
	case uint8:
		if reflect.TypeOf(i).String() == "byte" {
			if _, err := w.Write([]byte{STypeByte}); err != nil {
				return err
			}
		} else {
			if _, err := w.Write([]byte{STypeUint8}); err != nil {
				return err
			}
		}
		return PutUint8(w, i.(uint8))
	case uint16:
		if _, err := w.Write([]byte{STypeUint16}); err != nil {
			return err
		}
		return PutUint16(w, i.(uint16))
	case uint32:
		if _, err := w.Write([]byte{STypeUint32}); err != nil {
			return err
		}
		return PutUint32(w, i.(uint32))
	case uint64:
		if _, err := w.Write([]byte{STypeUint64}); err != nil {
			return err
		}
		return PutUint64(w, i.(uint64))
	case string:
		if _, err := w.Write([]byte{STypeString}); err != nil {
			return err
		}
		return PutString(w, i.(string))
	case time.Time:
		if _, err := w.Write([]byte{STypeTime}); err != nil {
			return err
		}
		return PutTime(w, i.(time.Time))
	case time.Duration:
		if _, err := w.Write([]byte{STypeDuration}); err != nil {
			return err
		}
		return PutDuration(w, i.(time.Duration))
	case []string:
		if _, err := w.Write([]byte{STypeStringArray}); err != nil {
			return err
		}
		return PutStringArray(w, i.([]string))
	default:
		return fmt.Errorf("unsupported context value")
	}
}

func PutByte(w io.Writer, b byte) error {
	if _, err := w.Write([]byte{b}); err != nil {
		return err
	}
	return nil
}

func ReadByte(r io.Reader) (byte, error) {
	bytes := make([]byte, 1)
	if _, err := r.Read(bytes); err != nil {
		return 0, err
	}
	return bytes[0], nil
}

func PutInt(w io.Writer, i int) error {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(i))
	if l, err := w.Write(bytes); err != nil || l != 8 {
		return err
	}
	return nil
}

func ReadInt(r io.Reader) (int, error) {
	bytes := make([]byte, 8)
	if l, err := r.Read(bytes); err != nil || l != 8 {
		return 0, err
	}
	return int(binary.BigEndian.Uint64(bytes)), nil
}

func PutInt8(w io.Writer, i int8) error {
	_, err := w.Write([]byte{byte(i)})
	return err
}

func ReadBool(r io.Reader) (bool, error) {
	b, err := ReadByte(r)
	if err != nil {
		return false, err
	}
	return b != 0, nil
}

func PutBool(w io.Writer, b bool) error {
	by := byte(0)
	if b {
		by = 1
	}
	return PutByte(w, by)
}

func ReadInt8(r io.Reader) (int8, error) {
	bytes := make([]byte, 1)
	if _, err := r.Read(bytes); err != nil {
		return 0, err
	}
	return int8(bytes[0]), nil
}

func PutInt16(w io.Writer, i int16) error {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, uint16(i))
	if l, err := w.Write(bytes); err != nil || l != 2 {
		return err
	}
	return nil
}

func ReadInt16(r io.Reader) (int16, error) {
	bytes := make([]byte, 2)
	if l, err := r.Read(bytes); err != nil || l != 2 {
		return 0, err
	}
	return int16(binary.BigEndian.Uint16(bytes)), nil
}

func PutInt32(w io.Writer, i int32) error {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(i))
	if l, err := w.Write(bytes); err != nil || l != 4 {
		return err
	}
	return nil
}
func ReadInt32(r io.Reader) (int32, error) {
	bytes := make([]byte, 4)
	if l, err := r.Read(bytes); err != nil || l != 4 {
		return 0, err
	}
	return int32(binary.BigEndian.Uint32(bytes)), nil
}

func PutInt64(w io.Writer, i int64) error {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(i))
	if l, err := w.Write(bytes); err != nil || l != 8 {
		return err
	}
	return nil
}
func ReadInt64(r io.Reader) (int64, error) {
	bytes := make([]byte, 8)
	if l, err := r.Read(bytes); err != nil || l != 8 {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(bytes)), nil
}

func PutUint(w io.Writer, i uint) error {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(i))
	if l, err := w.Write(bytes); err != nil || l != 8 {
		return err
	}
	return nil
}

func ReadUint(r io.Reader) (uint, error) {
	bytes := make([]byte, 8)
	if l, err := r.Read(bytes); err != nil || l != 8 {
		return 0, err
	}
	return uint(binary.BigEndian.Uint64(bytes)), nil
}

func PutUint8(w io.Writer, i uint8) error {
	_, err := w.Write([]byte{byte(i)})
	return err
}

func ReadUint8(r io.Reader) (uint8, error) {
	bytes := make([]byte, 1)
	if _, err := r.Read(bytes); err != nil {
		return 0, err
	}
	return uint8(bytes[0]), nil
}

func PutUint16(w io.Writer, i uint16) error {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, i)
	if l, err := w.Write(bytes); err != nil || l != 2 {
		return err
	}
	return nil
}

func ReadUint16(r io.Reader) (uint16, error) {
	bytes := make([]byte, 2)
	if l, err := r.Read(bytes); err != nil || l != 2 {
		return 0, err
	}
	return binary.BigEndian.Uint16(bytes), nil
}

func PutUint32(w io.Writer, i uint32) error {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, i)
	if l, err := w.Write(bytes); err != nil || l != 4 {
		return err
	}
	return nil
}
func ReadUint32(r io.Reader) (uint32, error) {
	bytes := make([]byte, 4)
	if l, err := r.Read(bytes); err != nil || l != 4 {
		return 0, err
	}
	return binary.BigEndian.Uint32(bytes), nil
}

func PutUint64(w io.Writer, i uint64) error {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, i)
	if l, err := w.Write(bytes); err != nil || l != 8 {
		return err
	}
	return nil
}
func ReadUint64(r io.Reader) (uint64, error) {
	bytes := make([]byte, 8)
	if l, err := r.Read(bytes); err != nil || l != 8 {
		return 0, err
	}
	return binary.BigEndian.Uint64(bytes), nil
}
func PutString(w io.Writer, str string) error {
	if err := PutUint32(w, uint32(len(str))); err != nil {
		return err
	}
	if l, err := w.Write([]byte(str)); err != nil || l != len(str) {
		return fmt.Errorf("%w : error put string", err)
	}
	return nil
}
func ReadString(r io.Reader) (string, error) {
	if l, err := ReadUint32(r); err != nil {
		return "", err
	} else {
		larr := make([]byte, l)
		if rl, err := r.Read(larr); err != nil || rl != int(l) {
			return "", err
		}
		return string(larr), nil
	}
}
func PutDuration(w io.Writer, duration time.Duration) error {
	return PutInt64(w, int64(duration))
}

func ReadDuration(r io.Reader) (time.Duration, error) {
	if i, err := ReadInt64(r); r != nil {
		return 0, err
	} else {
		return time.Duration(i), nil
	}
}
func PutTime(w io.Writer, t time.Time) error {
	ts := t.Format(time.RFC3339)
	return PutString(w, ts)
}
func ReadTime(r io.Reader) (time.Time, error) {
	if s, err := ReadString(r); err != nil {
		return time.Now(), err
	} else {
		return time.Parse(time.RFC3339, s)
	}
}

func PutStringArray(w io.Writer, strArr []string) error {
	if err := PutUint32(w, uint32(len(strArr))); err != nil {
		return err
	}
	for _, str := range strArr {
		if err := PutUint32(w, uint32(len(str))); err != nil {
			return err
		}
		if l, err := w.Write([]byte(str)); err != nil || l != len(str) {
			return fmt.Errorf("%w : error put string", err)
		}
	}
	return nil
}
func ReadStringArray(r io.Reader) ([]string, error) {
	if arrl, err := ReadUint32(r); err != nil {
		return nil, err
	} else {
		ret := make([]string, arrl)
		for i := 0; i < int(arrl); i++ {
			if l, err := ReadUint32(r); err != nil {
				return nil, err
			} else {
				larr := make([]byte, l)
				if rl, err := r.Read(larr); err != nil || rl != int(l) {
					return nil, err
				}
				ret[i] = string(larr)
			}
		}
		return ret, nil
	}
}
