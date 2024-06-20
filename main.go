package main

import (
	"fmt"
	"github.com/hokaccha/go-prettyjson"
	"zyxel/pkg/zyxel"
)

func main() {
	router := zyxel.NewZyxel("admin", "1mOlOkO2", "http://192.168.1.1")
	dev, err := router.DeviceList()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(dev)
	s, _ := prettyjson.Marshal(dev)
	fmt.Println(string(s))

	//input := "BDIYJLKOOQVIQUVNILLUJKPWAHUFRKAO12a6091ec5986f4c891c99d78ab7d2a2"
	//hash := sha256.New()
	//hash.Write([]byte(input))
	//hashedBytes := hash.Sum(nil)
	//hashedString := hex.EncodeToString(hashedBytes)
	//fmt.Println(hashedString)
}
