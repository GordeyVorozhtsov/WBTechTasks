package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func getTime() (time.Time, error) {
	time, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return time, err
}
func main() {
	timeFromNTP, err := getTime()
	fmt.Println(timeFromNTP, err)

}
