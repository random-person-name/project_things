package main

import (
	"flag"
	"fmt"
	"github.com/cbocovic/chord"
	"io"
	//"runtime"
	"crypto/sha256"
	"time"
)
//ChordApp is an interface for applications to run on top of a Chord DHT.
// type ChordApp interface {

// 	//Notify will alert the application of changes in the ChordNode's predecessor
// 	Notify(id [sha256.Size]byte, me [sha256.Size]byte, addr string)

// 	//Message will forward a message that was received through the DHT to the application
// 	Message(data []byte) []byte
// }
type word int
func (*word) Notify(id [sha256.Size]byte, me [sha256.Size]byte, addr string) {
	fmt.Printf("Notify called for addr ", addr)
}
func (*word) Message(data []byte) [] byte {
	fmt.Printf("message called!")
	return data
}

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())

	var startaddr string

	//set up flags
	numPtr := flag.Int("num", 5, "the size of the DHT you wish to test")
	startPtr := flag.Int("start", 1, "ipaddr to start from")

	flag.Parse()
	num := *numPtr
	start := *startPtr

	low := (1 + start) % 256
	middle := ((1 + start) / 256) % 256
	high := ((1 + start) / (256 * 256)) % 256
	startaddr = fmt.Sprintf("127.%d.%d.%d:8888", high, middle, low)

	fmt.Printf("Joining %d server starting at %s!\n", 1, startaddr)

	list := make([]*chord.ChordNode, num) //num)
	if start == 1 {

		me := new(chord.ChordNode)
		me = chord.Create(startaddr)
		list[0] = me
	} else {
		me := new(chord.ChordNode)
		me, _ = chord.Join(startaddr, "127.0.0.2:8888")
		list[0] = me
	}

	for i := 1; i < num; i++ {
		//join node to network or start a new network
		time.Sleep(time.Second)
		node := new(chord.ChordNode)
		low := (1 + start + i) % 256
		middle := ((1 + start + i) / 256) % 256
		high := ((1 + start + i) / (256 * 256)) % 256
		addr := fmt.Sprintf("127.%d.%d.%d:8888", high, middle, low)
		fmt.Printf("Joining %d server starting at %s!\n", 1, addr)
		node, _ = chord.Join(addr, startaddr)
		node.Register(byte(0) ,new(word))
		list[i] = node
		fmt.Printf("Joined server: %s.\n", addr)
	}
	//block until receive input
Loop:
	for {
		var cmd string
		var index int
		_, err := fmt.Scan(&cmd)
		switch {
		case cmd == "info":
			//print out successors and predecessors
			fmt.Printf("Node\t\t Successor\t\t Predecessor\n")
			// for _, node := range list {
			// 	// fmt.Printf("%s\n", node.Info()) //wrf
			// }
		case cmd == "fingers":
			//print out finger table
			fmt.Printf("Enter index of desired node: ")
			fmt.Scan(&index)
			if index >= 0 && index < len(list) {
				node := list[index]
				fmt.Printf("\n%s", node.ShowFingers())
			}
		case cmd == "succ":
			//print out successor list
			fmt.Printf("Enter index of desired node: ")
			fmt.Scan(&index)
			if index >= 0 && index < len(list) {
				node := list[index]
				fmt.Printf("\n%s", node.ShowSucc())
			}
		case err == io.EOF:
			break Loop
		}

	}
	for _, node := range list {
		node.Finalize()
	}

}
