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


func TestGet(t *testing.T) {
	cache := New(30*time.Second, 40 *time.Second)

	val := cache.Get("key")
	fmt.Println("val", val)
}


func TestexpireRoundAndShred(t *testing.T) {
	round, shred := expireRoundAndShred(1000)

	if round != 1 || shred != 1000 {
		t.Fail("expireRoundAndShred error")
	}
}
