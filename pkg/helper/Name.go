package helper

import (
	"bufio"
	"embed"
	"fmt"
	"math/rand"
)

var (
	//go:embed Names.txt
	res embed.FS
)

func RandomName() string {
	names := make([]string, 0)
	f, err := res.Open("Names.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		names = append(names, scanner.Text())
	}
	n1 := names[rand.Intn(len(names))]
	n2 := names[rand.Intn(len(names))]
	return fmt.Sprintf("%s%s", n1, n2)
}
