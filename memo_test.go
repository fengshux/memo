package memo


import (
	"testing"
	"fmt"
	"time"
)


func TestSetAndGet(t *testing.T) {

	memo := New(30*time.Second, 40 *time.Second)
	for i := 0 ; i < 1000; i++ {
		go func ( v int){
			memo.Set(fmt.Sprintf("key%d", v), v)
			val := memo.Get(fmt.Sprintf("key%d", v))
			fmt.Println(val)	
		} (i)
	}
	time.Sleep(50 * time.Second)
	
	for i := 0 ; i < 1000; i++ {
		go func (v int){		
			val := memo.Get(fmt.Sprintf("key%d", v))
			fmt.Println(val)	
		} (i)
	}
}
