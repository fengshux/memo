package memo


import (
	"testing"
	"fmt"
)


func TestSetAndGet(t *testing.T) {


	memo := Default()

	for i := 0 ; i < 1000; i++ {
		go func (){
			memo.Set("key", i)
			val := memo.Get("key")
			fmt.Println(val)
		} ()
	}
}
