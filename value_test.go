package color

import (
	"os"
	"testing"
)

func BenchmarkColorFuncs(b *testing.B) {
	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()
	os.Stdout = os.NewFile(0, os.DevNull)
	for i := 0; i < b.N; i++ {
		Black("hello from %s", "black")
		Green("hello from %s", "green")
		Red("hello from %q.  i'm %d", "red", 23)
	}
}

func BenchmarkColorFuncsParallel(b *testing.B) {
	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()
	os.Stdout = os.NewFile(0, os.DevNull)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Black("hello from %s xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "black")
			Green("hello from %s yyyyyhfhhdhehehehehhskdkdkdkdkdkdkdkkkekkekekekeekekkk", "green")
			Red("hello from %q.  i'm %d", "red", 23)
			Black("more blach stuff")
			Green("hello from %s xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "green")
			Red("hello from %q.  i'm %d", "red", 23)
		}
	})
}
