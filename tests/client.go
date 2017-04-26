package main

import (
	// "flag"
	"fmt"
	"github.com/cbocovic/chord"
	"io"
	"math/big"	//"runtime"
	"os"
	"net/rpc/jsonrpc"
	"net"
	"crypto/sha256"
	// "bufio"
)
var node *chord.ChordNode
var prev_pred *chord.Finger
type word int
func (*word) Notify(id [sha256.Size]byte, me [sha256.Size]byte, addr string) {
	fmt.Printf("Notify called for addr ", addr)
	// fmt.Printf("",(id), me, node.Our_Predecessor_id())

	//Checking if new node or old node deleted
	added := false
	if(prev_pred == nil){
		fmt.Println("New node Added")
		prev_pred = node.Our_Predecessor()
		added = true
	}else{
		curId := me
		predId:= prev_pred.Get_id()
		newId := id

		if( chord.InRange(newId, predId, curId)){
			fmt.Println("New node Added")
			added = true
		}else{
			fmt.Println("Node Deleted")
			added = false
		}
	}

	var newTable  map[string]string
	var sendTable map[string]string
	if(added == true){
		for k := range chord.D.Data.Table{
			keyId := sha256.Sum256([]byte(k))
			//Checking if equal
			kInt := new(big.Int)
			nInt := new(big.Int)
			kInt.SetBytes(keyId[:sha256.Size])
			nInt.SetBytes(me[:sha256.Size])

			if(kInt.Cmp(nInt)==1 || chord.InRange(keyId, id, me)){
				newTable[k] = chord.D.Data.Table[k]
			}else{
				sendTable[k] = chord.D.Data.Table[k]
			}
		}
		dst := new(chord.ChordNode)
		dst.Set_ipaddr(addr)
		node.Lelo(*dst, chord.KVstore{sendTable})
		node.Update_repo_succ(chord.KVstore{newTable})
	}else{
		
		newTable = chord.D.Data.Table

		for k := range chord.D.Replicas[prev_pred.Get_id()].Table{
			newTable[k] = chord.D.Replicas[prev_pred.Get_id()].Table[k]
		}
		node.Remove_repo_succ(prev_pred.Get_id())
		node.Update_repo_succ(chord.KVstore{newTable})
	}

	chord.D.Data.Table = newTable
	prev_pred = node.Our_Predecessor()
	return

}
func (*word) Message(data []byte) [] byte {
	fmt.Printf("message called!")
	return data
}

func main() {
	prev_pred = nil

	// if len(os.Args) == 2 {
	// 	node = new(chord.ChordNode)
	// 	node = chord.Create(os.Args[1]+":8888")
	// }
	// if len(os.Args) > 2 {
	// 	var err error
	// 	node, err = chord.Join(os.Args[1]+":8888", os.Args[2]+":8888")
	// 	check_err(err)
	// 	fmt.Printf("Joined server: %s.\n", os.Args[2])
	// 	// use(node)
	// }
	fmt.Printf("I am contacting %s:9000!\n", os.Args[1])
	// node.Register(byte(0) ,new(word))
	// use(node)
	// list[i] = node
	//block until receive input
    // reader := bufio.NewReader(os.Stdin)
Loop:
	for {
		var cmd string
		// var index int
	    fmt.Print("Enter Command: ")
	    // _, err := reader.ReadString('\n')
		_, err := fmt.Scanf("%s",&cmd)
		// fmt.Println("here")
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}
		value := ""
		key := ""		
		switch {
			case cmd == "g":
				//print out successors and predecessors
				// fmt.Printf("Node\t\t Successor\t\t Predecessor\n")
				client, err := net.Dial("tcp", os.Args[1]+":9000")
				check_err(err)
				c := jsonrpc.NewClient(client)
				fmt.Scan(&key)
				err = c.Call("Chord.Get", key, &value)
				check_err(err)
				fmt.Println("Value got :", value)
				// for _, node := range list {
				// 	// fmt.Printf("%s\n", node.Info()) //wrf
				// }
			case cmd == "p":
				client, err := net.Dial("tcp", os.Args[1]+":9000")
				check_err(err)
				c := jsonrpc.NewClient(client)
				fmt.Scan(&key)
				fmt.Scan(&value)
				var t int
				fmt.Println("key, value", key, value)
				err = c.Call("Chord.Put", &chord.PutArg{key, value}, &t)
				// err = c.Call("Chord.Put", &chord.PutArg{key, value}, &t)
				check_err(err)
				//print out finger table
				// fmt.Printf("Enter index of desired node: ")
				// fmt.Scan(&index)
				// if index >= 0 && index < len(list) {
				// 	node := list[index]
					// fmt.Printf("\n%s", node.ShowFingers())
				// }
			case cmd == "succ":
				//print out successor list
				// fmt.Printf("Enter index of desired node: ")
				// fmt.Scan(&index)
				// if index >= 0 && index < len(list) {
				// 	node := list[index]
					fmt.Printf("\n%s", node.ShowSuccessor())
					// fmt.Printf("\n%s", node.ShowSucc())
				// }
			case err == io.EOF:
				fmt.Printf("hello")
				break Loop
		}

	}
	// for _, node := range list {
		node.Finalize()
	// }
}
func check_err(err error) {
	if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}
}
func use(vals ...interface{}) {
    for _, val := range vals {
        _ = val
    }
}
