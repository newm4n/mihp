package helper

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Properties map[string]string

func (props Properties) String() string {
	var buff = &bytes.Buffer{}
	buff.WriteString("\n")
	table := tablewriter.NewWriter(buff)
	table.SetHeader([]string{"NO", "KEY", "VALUE"})

	keys := make([]string, 0)
	for k := range props {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for k, v := range keys {
		table.Append([]string{strconv.Itoa(k + 1), v, props[v]})
	}
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
	return buff.String()
}

func NewProperties() Properties {
	return make(Properties)
}

func (prop Properties) FromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	dataBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	return prop.FromBytes(dataBytes)
}

func (prop Properties) ToFile(filePath string, overwrite bool) error {
	byteToWrite := prop.ToBytes()
	if _, err := os.Stat(filePath); err == nil {
		if overwrite {
			os.Remove(filePath)
		} else {
			return fmt.Errorf("file %s already exist", filePath)
		}
	}
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("%w : file %s cannot be crated", err, filePath)
	}
	defer file.Close()
	_, err = file.Write(byteToWrite)
	return err
}

func (prop Properties) FromBytes(data []byte) error {
	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		line := scanner.Text()
		key := strings.TrimSpace(line[:strings.Index(line, "=")])
		val := strings.TrimSpace(line[strings.Index(line, "=")+1:])
		prop[key] = val
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (prop Properties) ToBytes() []byte {
	buff := &bytes.Buffer{}
	for k, v := range prop {
		buff.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}
	return buff.Bytes()
}
