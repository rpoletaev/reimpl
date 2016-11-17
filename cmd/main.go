package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rpoletaev/reimpl"
)

func main() {
	con, err := reimpl.Dial("localhost", ":2020")
	if err != nil {
		println(err)
	}

	defer con.Close()

	for {
		fmt.Print(">")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		splCmd := strings.Split(scanner.Text(), " ")

		if len(splCmd) == 0 && splCmd[0] == "" {
			continue
		}

		prs := parseParams(splCmd[1:])
		resp, err := con.Cmd(splCmd[0], prs...)

		if err != nil {
			fmt.Println(err)
		}

		switch r := resp.(type) {
		case string:
			println(r)
			break
		case []byte:
			println(string(r))
			break
		case []interface{}:
			for _, item := range r {
				println(string(item.([]byte)))
			}
			break
		case int:
			println("(integer) ", r)
			break
		case int64:
			println("(integer) ", r)
			break
		case error:
			println("ERROR: ", r.Error())
			break
		case nil:
			println("nil")
		default:
			fmt.Printf("Unexpected type: %t\n", r)
			break
		}
	}
}

func parseParams(spl []string) []interface{} {
	if len(spl) == 0 {
		return nil
	}

	prs := make([]interface{}, len(spl))
	for i, param := range spl {
		p, err := strconv.ParseInt(param, 10, 64)
		if err == nil {
			prs[i] = p
			continue
		}

		if strings.TrimSpace(param) == "" {
			continue
		}

		prs[i] = param
	}

	return prs
}
