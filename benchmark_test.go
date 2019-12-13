package memo

import(
	"testing"
	"fmt"
	"time"
)

var memo *Memo

func BenchmarkSet(b *testing.B) {
	memo = New(30*time.Second, 40 *time.Second)
	for n := 0; n < b.N; n++ {
		memo.Set(fmt.Sprintf("key%d", n), n)
	}
}

func BenchmarkGet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		memo.Get(fmt.Sprintf("key%d", n))		
	}
}
