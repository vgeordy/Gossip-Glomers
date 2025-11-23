/*

https://fly.io/dist-sys/2/

Test:

./maelstrom test -w unique-ids --bin ~/go/bin/maelstrom-unique-ids \
  --time-limit 30 --rate 1000 --node-count 3 \
  --availability total --nemesis partition

*/

package main

import (
	"crypto/rand"
	"fmt"

	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

/*
https://www.rfc-editor.org/rfc/rfc4122

# UIUD - universally unique identifier

128 bit number (16 bytes)

represented as 32 hexadecimal digits split into five groups

uses random number generator (v4)

f81d4fae-7dec-11d0-a765-00a0c91e6bf6

// time_low-t_mid-thv-
// f81d4fae-7dec-11d0-a765-00a0c91e6bf6
UUID                   = time-low "-" time-mid "-"

	                     time-high-and-version "-"
	                     clock-seq-and-reserved
	                     clock-seq-low "-" node
	time-low               = 4hexOctet
	time-mid               = 2hexOctet
	time-high-and-version  = 2hexOctet
	clock-seq-and-reserved = hexOctet
	clock-seq-low          = hexOctet
	node                   = 6hexOctet
	hexOctet               = hexDigit hexDigit
	hexDigit =
	      "0" / "1" / "2" / "3" / "4" / "5" / "6" / "7" / "8" / "9" /
	      "a" / "b" / "c" / "d" / "e" / "f" /
	      "A" / "B" / "C" / "D" / "E" / "F"

// output example: f81d4fae-7dec-11d0-a765-00a0c91e6bf6
*/
func UUIDv4() string {

	// Set all the other bits to randomly (or pseudo-randomly) chosen values.
	// https://stackoverflow.com/questions/32349807/how-can-i-generate-a-random-int-using-the-crypto-rand-package

	id := make([]byte, 16)

	_, err := rand.Read(id)
	if err != nil {
		log.Fatalf("Failed to read random bytes: %v", err)
	}

	// fmt.Printf("Generated secure bytes (hex): %x\n", id)

	//	Set the two most significant bits (bits 6 and 7) of the
	//	clock_seq_hi_and_reserved to zero and one, respectively.

	// ex. a765 where a7 represents clock-seq-and-reserved, where a and 7 are two seperate hexDigits and
	// 65 represents clock-seq-low, where 6 and 5 are two seperate hexDigits
	// hexOctet is two hexadigits (subtle but important)

	// msb of the would be a7 portion, and we have to set that to zero (7) and one (a) -> 10
	// clear top 2 bits

	// index:   0  1  2  3   4  5   6  7   8   9   10 11 12 13 14 15
	//  value:  xx xx xx xx  xx xx  xx xx  xx  xx  xx xx xx xx xx xx

	// a7 in binary: 10100111

	// clear top 2 bits -> 00100111
	id[8] &= 0x3F

	// set top 2 bits to 10 -> 10100111
	id[8] |= 0x80

	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number from Section 4.1.3. (0     1     0     0)

	// reset
	id[6] &= 0x0F

	// shifted by 4 to align and then OR'ed to set 0100 to the msb
	id[6] |= 4 << 4

	// fmt.Printf("%x\n", id)

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		id[0:4],
		id[4:6],
		id[6:8],
		id[8:10],
		id[10:16],
	)
}

func main() {

	// creates a new node
	n := maelstrom.NewNode()

	n.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// only difference from https://fly.io/dist-sys/1/
		reply := map[string]any{
			"type": "generate_ok",
			"id":   UUIDv4(),
		}

		return n.Reply(msg, reply)

	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}

}
