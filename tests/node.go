package main

import (
	// "github.com/fatih/color"
	"sync"
	"fmt"
	"github.com/cbocovic/chord"
	"io"
	"math/big"	//"runtime"
	"os"
	"crypto/sha256"
	// "bufio"
)
var node *chord.ChordNode
var prev_pred  chord.Finger
var flag bool
type word int
func (*word) Notify(id [sha256.Size]byte, me [sha256.Size]byte, addr string) {
	fmt.Printf(">>>Notify called for addr ", addr)
	// fmt.Printf("",(id), me, node.Our_Predecessor_id())

	//Checking if new node or old node deleted
	added := false
	if(flag == false){
		fmt.Println("New node Added - nil")
		added = true
	}else{
		curId := me
		predId:= prev_pred.Get_id()
		newId := id

		if( chord.InRange(newId, predId, curId)){
			fmt.Println("New node Added: ",prev_pred.Get_ipaddr())
			added = true
		}else{
			fmt.Println("Node Deleted: ",prev_pred.Get_ipaddr())
			added = false
		}
	}

	newTable := map[string]string {}
	sendTable := map[string]string{}

	if chord.D.Data.Table == nil {
		chord.D.Data.Table = make(map[string] string)
	}
	if(added == true){
		fmt.Println("nagaraju added")
		for k := range chord.D.Data.Table{
			keyId := sha256.Sum256([]byte(k))
			//Checking if equal
			kInt := new(big.Int)
			nInt := new(big.Int)
			kInt.SetBytes(keyId[:sha256.Size])
			nInt.SetBytes(me[:sha256.Size])

			if(kInt.Cmp(nInt)==0 || chord.InRange(keyId, id, me)){
				newTable[k] = chord.D.Data.Table[k]
			}else{
				sendTable[k] = chord.D.Data.Table[k]
			}
		}
		fmt.Println("sendTable: ",sendTable, "newtable :", newTable)
		// D.Lock.Lock()
		chord.D.Data.Table = newTable
		// D.Lock.Unlock()
		dst := new(chord.ChordNode)
		dst.Set_ipaddr(addr)
		node.Lelo(*dst, chord.KVstore{sendTable})
		node.Update_repo_succ(chord.KVstore{newTable})
	}else{
		newTable = chord.D.Data.Table
		fmt.Printf("1")
		for k := range chord.D.Replicas[prev_pred.Get_id()].Table{
			newTable[k] = chord.D.Replicas[prev_pred.Get_id()].Table[k]
		}
		fmt.Printf("2")
		// nod delete repo
		// D.Lock.Lock()
		chord.D.Data.Table = newTable
		// D.Lock.Unlock()
		node.Remove_repo_succ(prev_pred.Get_id())
		node.Update_repo_succ(chord.KVstore{newTable})
	}

	chord.D.Data.Table = newTable
	prev_pred = *(node.Our_Predecessor())
	fmt.Println("*******************b4 exiting",prev_pred.Get_ipaddr())
	flag= true
	return

}
func (*word) Message(data []byte) [] byte {
	fmt.Printf("message called!")
	return data
}

func main() {
	flag = false
	chord.D.Replicas = make(map[[sha256.Size]byte]*chord.KVstore)
	chord.D.Data = new(chord.KVstore)
	// t := new()
	chord.D.Data.Table = map[string] string{}
	chord.D.Lock = new(sync.RWMutex)

	if len(os.Args) == 2 {
		node = new(chord.ChordNode)
		node = chord.Create(os.Args[1]+":8888")
	}
	if len(os.Args) > 2 {
		var err error
		node, err = chord.Join(os.Args[1]+":8888", os.Args[2]+":8888")
		check_err(err)
		fmt.Printf("Joined server: %s:8888\n", os.Args[2])
		// use(node)
	}
	fmt.Printf("I am %s:8888!\n", os.Args[1])
	node.Register(byte(0) ,new(word))
	use(node)
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
		switch {
			case cmd == "info":
				//print out successors and predecessors
				// fmt.Printf("Node\t\t Successor\t\t Predecessor\n")
				fmt.Println("Prev_pred : ",prev_pred.Get_ipaddr())
				node.Info()
				
				// for _, node := range list {
				// 	// fmt.Printf("%s\n", node.Info()) //wrf
				// }
			case cmd == "fingers":
				//print out finger table
				// fmt.Printf("Enter index of desired node: ")
				// fmt.Scan(&index)
				// if index >= 0 && index < len(list) {
				// 	node := list[index]
					fmt.Printf("\n%s", node.ShowFingers())
				// }
			case cmd == "succ":
				//print out successor list
				// fmt.Printf("Enter index of desired node: ")
				// fmt.Scan(&index)
				// if index >= 0 && index < len(list) {
				// 	node := list[index]
					// fmt.Printf("\n%s", 
				node.ShowSuccessor()
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
