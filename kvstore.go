package chord

import "sync"
import "errors"
import "crypto/sha256"
import "fmt"


type KVstore struct{
	Table map[string] string
}

type Database struct{
	Lock *sync.RWMutex
	Data *KVstore
	Replicas map[[sha256.Size]byte] *KVstore
}

func Initialize(){
	// D = new(Database)
	D.Replicas = make(map[[sha256.Size]byte]*KVstore)
	D.Data = new(KVstore)
	t := new(map[string] string)
	D.Data.Table = *t
	D.Lock = new(sync.RWMutex)
	// return d
}

func (d *Database) Get(key string) (ret string,err  error){
	//Todo: write the get functionality here
	// defer d.Lock.Unlock()
	// d.Lock.Lock()
	val,ok := ((d.Data.Table))[key]
	if ok!=true{
		err = errors.New("Key not found")
		fmt.Println("Get: key ",key ," could not be found")
		return val, err
	}
	err = nil
		fmt.Println("Get: key ",key ," found")
	return val,err
}

func (d Database) Put(key string, value string)(err error){
	//Todo; write the insert functionality
	// defer d.Lock.Unlock()
	// d.Lock.Lock()
	fmt.Println("reached here")
	fmt.Println(d.Data.Table)
	// d.Data.Table = make(map[string] string)
	if d.Data.Table == nil {
		fmt.Println("Empty table detected!")
		d.Data.Table = make(map[string] string)
	}
	((d.Data.Table))[key] = value
	fmt.Println("Put: added key ",key," into the database")
	err = nil
	return
}