## reflect性能分析

### 实验设计
利用go benchmark进行性能测试。
+ 实验组：利用reflect调用方法
+ 对照组：直接调用对象的方法

```go
type Target struct {
    Val int
}

func (t *Target) Run() {
    t.Val++
}

func RunDirect(t *Target) {
    t.Run()
}

func RunByReflect(v interface{}) {
    f := reflect.ValueOf(v)
    f.MethodByName("Run").Call(nil)
}

func BenchmarkRunFuncDirect(b *testing.B) {
    var t = Target{}
    for i:=0;i<b.N;i++ {
        RunDirect(&t)
    }
}

func BenchmarkRunFuncReflect(b *testing.B) {
    var t = Target{}
    for i:=0;i<b.N;i++ {
        RunByReflect(&t)
    }
}
```

### 测试命令
```go
go test -bench=. -count=1 -benchmem
```

### 测试结果
```go
BenchmarkRunDirect-4            609154983                2.00 ns/op            0 B/op          0 allocs/op
BenchmarkRunByReflect-4          2377942               510 ns/op             128 B/op          4 allocs/op
```

### 简单得出结论
+ reflect会分配4次内存
+ 性能比是直接调用的差了数百倍。