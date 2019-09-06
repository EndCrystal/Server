package main

import (
	"fmt"
	"io/ioutil"

	"github.com/EndCrystal/Server/token"
)

func main() {
	pub, priv := token.GenerateKeys()
	ioutil.WriteFile("key.pub", pub[:], 0444)
	ioutil.WriteFile("key.priv", priv[:], 0400)
	fmt.Println("Generated")
}
