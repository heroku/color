package color

import (
	"bytes"
	"os"
	"testing"
)

func TestColor(t *testing.T) {
	tt := []struct{
		text string
		code Attribute
		want string
	}{
		{"white", FgWhite, "" },
	}

	for _, tc := range tt {
		t.Run(tc.text, func(t *testing.T){
			var buff bytes.Buffer
			v, _ := New(&buff, tc.code)
			_, _ = v.Print(tc.text)
			got := buff.String()
			t.Log(got)
			if got != tc.want {
				t.Logf("got  %q", got)
				t.Logf("want %q", tc.want)
				t.Fatal()
			}
		})
	}
}

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
