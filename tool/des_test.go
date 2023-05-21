package tool

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestDes(t *testing.T) {
	e := DesToken{}
	token, success := e.Encrypt(strconv.FormatUint(uint64(1), 10) + "," + strconv.FormatInt(time.Now().Unix(), 10))

	fmt.Println(token, success)

	ori, success := e.Decrypt(token)
	fmt.Println(ori, success)
}

func TestStr(t *testing.T) {
	s := " catiga "
	v := strings.Trim(s, " ")
	fmt.Println(v)
}
