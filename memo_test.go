package memo


import (
	"testing"
	"fmt"
	"time"
)


func TestSetAndGet(t *testing.T) {

	memo := Default()
	for i := 0 ; i < 1000; i++ {
		go func (){
			memo.Set(fmt.Sprintf("key%d", i), i)
			val := memo.Get(fmt.Sprintf("key%d", i))
			fmt.Println(val)	
		} ()
	}
	time.Sleep(2 * time.Minute)
	
	for i := 0 ; i < 1000; i++ {
		go func (){		
			val := memo.Get(fmt.Sprintf("key%d", i))
			fmt.Println(val)	
		} ()
	}
}
