package chord

import (

	"crypto/sha256"
	"fmt"
	// "math/big"
	"math/rand"
	"net"
	"os"
	"strings"
	"net/rpc/jsonrpc"
	"errors"
	"net/rpc"
	"github.com/fatih/color"
)

type Chord int 

var D Database
var Gnode *ChordNode

func Modifyip(addr string, port string)string{
	port = "9000"
	if port == ""{
		port = "9000"
	}
	ans := strings.Split(addr,":")
	newaddr := ans[0]+":"+port
	return newaddr
}

func (c *Chord) Update(arg *Update_Arg, rep *int) (err error){
	fmt.Println("our_rpc update function has been called with ",arg.Option," ",arg.Ipaddr)
	if arg.Option == "update"{
		D.Replicas[arg.Id] = new(KVstore)
		*(D.Replicas[arg.Id]) = arg.Rep
		
	} else if arg.Option == "delete"{
		fmt.Println("Deleting address of ", arg.Ipaddr, " ", arg.Id)
		delete(D.Replicas,arg.Id)
		fmt.Println("Akash garg check ",D.Replicas[arg.Id])
	} else if arg.Option == "lelo"{
		D.Data = new(KVstore)
		*D.Data = arg.Rep
		fmt.Println("our_rpc lelo table received",((D.Data).Table))
		dst := new(ChordNode)
		color.Blue(string(len(D.Data.Table)))
		succ_list := Gnode.Our_SuccessorList()
		flag := make(map[string]bool)
		for i := 0; i < sha256.Size*8; i++ {
			if succ_list[i].Get_ipaddr()==Gnode.Get_ipaddr() {
				break
			}
			finger := succ_list[i]
			if len(finger.Get_ipaddr()) <= 10 {
				continue
			} 
			if flag[finger.Get_ipaddr()] == true {
				continue
			}
			dst.Set_ipaddr(finger.Get_ipaddr())
			flag[finger.Get_ipaddr()] = true
			Gnode.Update(*dst, arg.Rep)
			// if finger.ipaddr != "" {
			// 	if i == 0 || finger.ipaddr != prevfinger.ipaddr {
			// 		dst.ipaddr = finger.ipaddr
			// 		Gnode.Update(dst, arg.Rep)
			// 	}
			// }
			// *prevfinger = *finger
			
		}
	}else{
		fmt.Println("Unrecognized option description:",arg.Option)
		fmt.Println(arg)
		err = errors.New("Unrecognzed option : ")
		return err
	}
	return nil
}
type RGetArg struct {
	Key string
	Id [sha256.Size] byte
}
func (node *Chord) Get(key string, value *string) (err error) {
	fmt.Println("Get function called")
	addr, err := Lookup(sha256.Sum256([]byte(key)), Gnode.ipaddr)
	if addr != Gnode.Get_ipaddr() {
		fmt.Println("Getting to ",addr)
		client, err := net.Dial("tcp", Modifyip(addr, ""))
		check_err(err)
		c := jsonrpc.NewClient(client)
		err = c.Call("Chord.Get", key, &value)	
	} else {
		defer D.Lock.RUnlock()
		D.Lock.RLock()
		succ_list := Gnode.Our_SuccessorList()
		counter := rand.Uint32()%1
		for i := 0; i < sha256.Size*8; i++ {
			if succ_list[i].Get_ipaddr()==Gnode.Get_ipaddr() {
				break
			}
			finger := succ_list[i]
			if len(finger.Get_ipaddr()) <= 10 {
				continue
			}
			counter = counter+1
		}
		x := rand.Uint32() % (counter + 1)
		if x == counter {
			addr = Gnode.Get_ipaddr()
		}else {
			addr = succ_list[x].Get_ipaddr()
		}
		if addr != Gnode.Get_ipaddr() {
			color.Red("Can serve but load balancing to ",addr)
			client, err := net.Dial("tcp", Modifyip(addr, ""))
			check_err(err)
			c := jsonrpc.NewClient(client)
			t := ""
			err = c.Call("Chord.RGet", RGetArg{key,Gnode.Get_id()}, &t)
			*value = t
			check_err(err)
		} else {
			color.Green("Serving Get request for key from self ",key)
			*value, err = D.Get(key)
		}
	}
	return
}

func (*Chord) RGet(inp RGetArg, value *string) (err error){
	defer D.Lock.RUnlock()
	D.Lock.RLock()
	color.Green("Serving Loadbalanced RGet request for key: ",inp.Key)
	*value, _ = D.Replicas[inp.Id].Table[inp.Key]
	return
}
type PutArg struct {
	Key string
	Value string
}
func (node *Chord) Put(inp *PutArg, ret *int) (err error) {
	color.Cyan("Put function called :", inp.Key," ",inp.Value)
	addr, err := Lookup(sha256.Sum256([]byte(inp.Key)), Gnode.ipaddr)
	if addr != Gnode.Get_ipaddr() {
		fmt.Println("Putting to ",addr)
		client, err := net.Dial("tcp", Modifyip(addr, ""))
		check_err(err)
		c := jsonrpc.NewClient(client)
		t := 0
		err = c.Call("Chord.Put", inp, &t)	
	} else {
		defer D.Lock.Unlock()
		D.Lock.Lock()
		color.Red("Putting myself!")
		err = D.Put(inp.Key, inp.Value) // doesn't release lock till all repos have synced up
		Gnode.Update_repo_succ(*D.Data)
	}
	*ret = 0
	return
}

func (node* ChordNode)TcpServer(){
	fmt.Println("TcpServer function has been called to start tcp server.")
	chrd := new(Chord)
	rpc.Register(chrd)
	newaddr := Modifyip(node.ipaddr,"9000")

	tcpAddr, err := net.ResolveTCPAddr("tcp", newaddr) //Need some more explanation here
	if err != nil{
		fmt.Println(err)
		return
	}else{
		fmt.Println("TCPAddress Successfully Resolved")
	}
	
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil{
		fmt.Println(err)
		return
	}else{
		fmt.Println("Listening on Socket "+newaddr)
	}
	
	
	for {
		conn, err := listener.Accept()
		fmt.Print("Accepted  connection from :")
		fmt.Println(conn.RemoteAddr())
		if err != nil {
			fmt.Println("Error: ",err.Error())
			continue
		}
		go jsonrpc.ServeConn(conn)
		fmt.Println("Connection Terminated")
	}
	return
}

func check_err(err error) {
	if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}
}