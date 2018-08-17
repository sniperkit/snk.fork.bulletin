/*
Sniperkit-Bot
- Status: analyzed
*/

package error

import (
	"log"
)

func CheckError(e error) {
	if e != nil {
		log.Fatalf("error: %v\n", e)
	}
}
