package main

import (
    "flag" // command line option parser
    "fmt"
    "io"
    "math"
    "math/cmplx"
    "math/rand"
    "os"
    "reflect"
    "runtime"
    "strings"
    "time"
)

const (
    Space   = " "
    Newline = "\n"
)

type ErrNegativeSqrt float64
func (e ErrNegativeSqrt) Error() string {
    return fmt.Sprintf("cannot sqrt negative number: %.02f", float64(e))
}

func mysqrt(x float64) (float64, error) {
    if x < 0 {
        return 0, ErrNegativeSqrt(x)
    }
    return math.Sqrt(x), nil
}

func sqnewt_factory(z, x float64) func () float64 {
    sqnewt := func () float64 {
        z = z - ((z * z - x) / (2 * z))
        return z
    }
    return sqnewt
}

func TestNewtonSquareRoot() {
    const (
        TOLERANCE = 0.00000001
        MAX = 101
    )

    for v := 1; v < MAX; v++ {
        sqnewt := sqnewt_factory(float64(v), float64(v))
        fmt.Print(v)
        last := 0.0
        for {
            sqrt := sqnewt()
            fmt.Printf("\t %.04f", sqrt)
            if last == 0.0 {
                last = sqrt
            } else {
                if math.Abs(last - sqrt) < TOLERANCE {
                    // done
                    break
                } else {
                    last = sqrt
                }
            }
        }
        fmt.Printf("\t(last=%.08f, math.sqrt=%.08f)\n", last, math.Sqrt(float64(v)))
    }
}

func word_count(s string) map[string]int {
    var counts_by_word = map[string]int{}
    for _, word := range strings.Fields(s) {
        counts_by_word[word] += 1
    }
    return counts_by_word
}

func TestWordCounts() {
    counts := word_count(`When beetles fight these battles in a bottle with their paddles 
                    and the bottle's on a poodle and the poodle's eating noodles... 
                    ...they call this a muddle puddle tweetle poodle beetle noodle 
                    bottle paddle battle.`) 
    for w,c := range counts {
        fmt.Println(w, c)
    }
}

func Pic(dx, dy int) [][]uint8 {
    var ret = make([][]uint8, dy)

    for y := 0; y < dy; y++ {
        ret[y] = make([]uint8, dx)
        for x := 0; x < dx; x++ {
            ret[y][x] = uint8(x&y)
        }
    }
    return ret
}

func fibonacci() func() int {
    var prev int = 0
    var curr int = 1
    f := func() int {
        tmp := curr
        ret := curr
        curr = prev + curr
        prev = tmp
        return ret
    }

    return f
}

func TestFib() {
    f := fibonacci()
    for i := 0; i < 20; i++ {
        fmt.Println(f())
    }
}

func Cbrt(x complex128) complex128 {
    z := complex128(1)

    for i := 0; i < 10; i++ {
        z = z - ((cmplx.Pow(z, 3) - x)  / (3 * cmplx.Pow(z, 2)))
    }

    return z
}

type Abser interface {
    Abs() float64
}

type Vertex struct {
    x, y float64
}
func (v *Vertex) Abs() float64 {
    return math.Sqrt(v.x*v.x + v.y*v.y)
}

type myfloat float64
func (f myfloat) Abs() float64 {
    v := f
    if v < 0 {
        return float64(v * -1)
    }
    return float64(v)
}

type rot13Reader struct {
    r io.Reader
}

func (r *rot13Reader) Read(p []byte) (n int, err error) {
    // read from input stream
    n, err = r.r.Read(p)
    rot13 := func(p []byte) []byte {
        for i := 0; i < len(p); i++ {
            if p[i] > 64 && p[i] < 91 {
                p[i] += 13
                if p[i] > 90 {
                    p[i] = p[i] % 90 + 64
                }
            }
            if p[i] > 96 && p[i] < 123 {
                p[i] += 13
                if p[i] > 122 {
                    p[i] = p[i] % 122 + 96
                }
            }
        }
        return p
    }
    // translate chars
    p = rot13(p)
    return
}

func TestRot13() {
    r := rot13Reader{strings.NewReader("Lbh penpxrq gur pbqr!")}
    io.Copy(os.Stdout, &r)
    fmt.Print("\n")
}

//
// Tree Checker
//
type ChanInt chan int

type Tree struct {
    Left *Tree
    Value int
    Right *Tree
}

func TreeAddNode(root *Tree, value int) *Tree {
    var node *Tree

    parent := root
    for {
        if value > parent.Value {
            if parent.Right == nil {
                parent.Right = new(Tree)
                node = parent.Right
                break
            } else {
                parent = parent.Right
            }
        } else {
            if parent.Left == nil {
                parent.Left = new(Tree)
                node = parent.Left
                break
            } else {
                parent = parent.Left
            }
        }
    }
    node.Value = value
    return node
}

func TreeNew(value int) *Tree {
    const (
        SIZE int = 20
        MAX int = 100
    )
    root := new(Tree)
    root.Value = value

    // add some random nodes
    rand.Seed(time.Now().UnixNano())
    for i := 0; i < SIZE; i++ {
        n := rand.Intn(MAX)
        TreeAddNode(root, n)  
    }
    return root
}

func _TreeWalk(t *Tree, ch ChanInt) {
    // pre-order tree walk
    if t == nil {
        return
    }
    if t.Left != nil {
        _TreeWalk(t.Left, ch)
    }
    //fmt.Println(t.Value)
    ch <- t.Value
    if t.Right != nil {
        _TreeWalk(t.Right, ch)
    }
}

func TreeWalk(t *Tree, ch ChanInt) {
    _TreeWalk(t, ch)
    close(ch)
}

func TreeSame(t1, t2 *Tree) bool {
    t1_ch := make(ChanInt)
    t2_ch := make(ChanInt)

    go TreeWalk(t1, t1_ch)
    go TreeWalk(t2, t2_ch)

    // read from both channels comparing values
    for t1_val := range t1_ch {
        t2_val := <- t2_ch
        //fmt.Println(t1_val, t2_val)
        if t1_val != t2_val {
            return false
        }
    }

    return true
}

func TestTree() {
   // 
   t1 := TreeNew(1)
   t2 := TreeNew(1)

   fmt.Println("t1 == t1", TreeSame(t1, t1))
   fmt.Println("t1 == t2", TreeSame(t1, t2))
   fmt.Println("t2 == t1", TreeSame(t2, t1))
   fmt.Println("t2 == t2", TreeSame(t2, t2))
}

//
// Web Crawler
//
type Fetcher interface {
    // Fetch returns the body of URL and
    // a slice of URLs found on that page.
    Fetch(url string) (body string, urls []string, err error)
}

type request struct {
    url string
    depth int
}

type result struct {
    url   string
    body  string
    urls  []string
    depth int
    err   error
}

// crawler
type crawler struct {
    f   Fetcher
    req chan *request
    res chan *result
}

func (c *crawler) fetch(url string, depth int) {
    c.req <- &request{url, depth}
}

func (c *crawler) run() {
    for req := range c.req {
        if req.depth <= 0 {
            c.res <- &result{req.url, "", nil, req.depth,
                fmt.Errorf("depth exceeded: %s", req.url)}
            return
        }
        body, urls, err := c.f.Fetch(req.url)
        c.res <- &result{req.url, body, urls, req.depth, err}
    }
}

func Crawl(url string, depth int, f Fetcher, workers int) {
    c := &crawler{f, make(chan *request), make(chan *result)}
    for i := 0; i < workers; i++ {
        go c.run()
    }
    go c.fetch(url, depth)
    is_queued, n := map[string]bool{url: true}, 1
    for res := range c.res {
        n--
        if res != nil {
            if res.err != nil {
                fmt.Println(res.err)
            } else {
                fmt.Printf("found: %s %q [%d]\n", res.url, res.body, res.depth)
            }
            for _, url := range res.urls {
                if !is_queued[url] {
                    // queue up unfetched urls
                    go c.fetch(url, res.depth-1)
                    n++
                }
                is_queued[url] = true
            }
        }
        if n == 0 {
            // stop when we've fetch everything that was added
            return
        }
    }
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
    body string
    urls []string
}

func (f *fakeFetcher) Fetch(url string) (string, []string, error) {
    if res, ok := (*f)[url]; ok {
        return res.body, res.urls, nil
    }
    return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
func TestFetcher() {
    var fetcher = &fakeFetcher{
        "http://golang.org/": &fakeResult{
            "The Go Programming Language",
            []string{
                "http://golang.org/pkg/",
                "http://golang.org/cmd/",
            },
        },
        "http://golang.org/pkg/": &fakeResult{
            "Packages",
            []string{
                "http://golang.org/",
                "http://golang.org/cmd/",
                "http://golang.org/pkg/fmt/",
                "http://golang.org/pkg/os/",
            },
        },
        "http://golang.org/pkg/fmt/": &fakeResult{
            "Package fmt",
            []string{
                "http://golang.org/",
                "http://golang.org/pkg/",
            },
        },
        "http://golang.org/pkg/os/user/": &fakeResult{
            "Package user",
            []string{
                "http://golang.org/",
                "http://golang.org/pkg/",
                "http://golang.org/pkg/os/",
                "http://golang.org/pkg/os/user/user.go",
            },
        },
        "http://golang.org/pkg/os/": &fakeResult{
            "Package os",
            []string{
                "http://golang.org/",
                "http://golang.org/pkg/",
                "http://golang.org/pkg/os/user/",
            },
        },
    }

    Crawl("http://golang.org/", 4, fetcher, 4)
}


func get_func_name(i interface{}) string {
    return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func main() {
    flag.Parse() // Scans the arg list and sets up flags
    var s string = ""
    for i := 0; i < flag.NArg(); i++ {
        if i > 0 {
            s += Space
        }
        s += flag.Arg(i)
    }
    runner := func (f func()) {
        fmt.Printf("Testing %s ...\n", get_func_name(f))
        f()
        fmt.Printf("\n")

    }
    runner(TestNewtonSquareRoot)
    runner(TestWordCounts)
    runner(TestFib)
    runner(TestRot13)
    runner(TestTree)
    runner(TestFetcher)
}

