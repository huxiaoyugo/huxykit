## 新旧版本jsonMarshal性能对比

### 测试代码
```go
type Car struct {
	Name string
	Age int
	Brand string
	Address string
	Level  int
	*Wheel
}

type Wheel struct {
	Name string
	Brand string
	*Material
}

type Material struct {
	Name string
	Type string
}

var c = &Car{
		Name:    "BMW320",
		Age:     100,
		Brand:   "BMW",
		Address: "北京",
		Level:   1,
		Wheel:   &Wheel{
			Name:     "",
			Brand:    "",
			Material: &Material{
				Name: "橡胶",
				Type: "1",
			},
		},
	}

func BenchmarkJson(b *testing.B) {
	for i:=0;i<b.N;i++ {
		json.Marshal(c)
	}
}
```

### 各版本测试结果
#### 1.9.7
```go
goos: darwin
goarch: amd64
pkg: github.com/huxiaoyugo/huxykit/research
BenchmarkJson-4   	 2000000	       971 ns/op	     320 B/op	       2 allocs/op
BenchmarkJson-4   	 2000000	       898 ns/op	     320 B/op	       2 allocs/op
BenchmarkJson-4   	 2000000	       872 ns/op	     320 B/op	       2 allocs/op
BenchmarkJson-4   	 2000000	       870 ns/op	     320 B/op	       2 allocs/op
BenchmarkJson-4   	 2000000	       876 ns/op	     320 B/op	       2 allocs/op
BenchmarkJson-4   	 2000000	       903 ns/op	     320 B/op	       2 allocs/op
BenchmarkJson-4   	 2000000	       874 ns/op	     320 B/op	       2 allocs/op
BenchmarkJson-4   	 2000000	       878 ns/op	     320 B/op	       2 allocs/op
BenchmarkJson-4   	 2000000	       864 ns/op	     320 B/op	       2 allocs/op
BenchmarkJson-4   	 2000000	       876 ns/op	     320 B/op	       2 allocs/op
PASS
```

#### 1.11.13
```go
goos: darwin
goarch: amd64
pkg: github.com/huxiaoyugo/huxykit/research
BenchmarkJson-4   	 2000000	       692 ns/op	      96 B/op	       1 allocs/op
BenchmarkJson-4   	 2000000	       686 ns/op	      96 B/op	       1 allocs/op
BenchmarkJson-4   	 2000000	       694 ns/op	      96 B/op	       1 allocs/op
BenchmarkJson-4   	 2000000	       691 ns/op	      96 B/op	       1 allocs/op
BenchmarkJson-4   	 2000000	       702 ns/op	      96 B/op	       1 allocs/op
BenchmarkJson-4   	 2000000	       713 ns/op	      96 B/op	       1 allocs/op
BenchmarkJson-4   	 2000000	       682 ns/op	      96 B/op	       1 allocs/op
BenchmarkJson-4   	 2000000	       676 ns/op	      96 B/op	       1 allocs/op
BenchmarkJson-4   	 2000000	       679 ns/op	      96 B/op	       1 allocs/op
BenchmarkJson-4   	 2000000	       676 ns/op	      96 B/op	       1 allocs/op
PASS
```

#### 1.14.3
```go
BenchmarkJson
BenchmarkJson-4   	 2399896	       472 ns/op	      96 B/op	       1 allocs/op
BenchmarkJson-4   	 2612428	       455 ns/op	      96 B/op	       1 allocs/op
BenchmarkJson-4   	 2642835	       456 ns/op	      96 B/op	       1 allocs/op
BenchmarkJson-4   	 2625843	       460 ns/op	      96 B/op	       1 allocs/op
BenchmarkJson-4   	 2609326	       453 ns/op	      96 B/op	       1 allocs/op
BenchmarkJson-4   	 2644610	       454 ns/op	      96 B/op	       1 allocs/op
BenchmarkJson-4   	 2626513	       454 ns/op	      96 B/op	       1 allocs/op
BenchmarkJson-4   	 2585499	       453 ns/op	      96 B/op	       1 allocs/op
BenchmarkJson-4   	 2621565	       455 ns/op	      96 B/op	       1 allocs/op
BenchmarkJson-4   	 2637963	       456 ns/op	      96 B/op	       1 allocs/op
PASS
```

### benchstate对比分析
```go
benchstat 1.9.txt 1.14.txt 
name    old time/op    new time/op    delta
Json-4     879ns ± 3%     454ns ± 0%  -48.29%  (p=0.000 n=9+8)

name    old alloc/op   new alloc/op   delta
Json-4      320B ± 0%       96B ± 0%  -70.00%  (p=0.000 n=10+10)

name    old allocs/op  new allocs/op  delta
Json-4      2.00 ± 0%      1.00 ± 0%  -50.00%  (p=0.000 n=10+10)


benchstat 1.11.txt 1.14.txt 
name    old time/op    new time/op    delta
Json-4     689ns ± 3%     454ns ± 0%  -34.04%  (p=0.000 n=10+8)

name    old alloc/op   new alloc/op   delta
Json-4     96.0B ± 0%     96.0B ± 0%     ~     (all equal)

name    old allocs/op  new allocs/op  delta
Json-4      1.00 ± 0%      1.00 ± 0%     ~     (all equal)
```


## 1.9与1.11版本的差别

+ 对encodeState对象做了缓存池
```go
// 1.11
func Marshal(v interface{}) ([]byte, error) {
	e := newEncodeState()
	err := e.marshal(v, encOpts{escapeHTML: true})
	if err != nil {
		return nil, err
	}
	buf := append([]byte(nil), e.Bytes()...)
	e.Reset()
	encodeStatePool.Put(e)
	return buf, nil
}

// 1.9
func Marshal(v interface{}) ([]byte, error) {
	e := &encodeState{}
	err := e.marshal(v, encOpts{escapeHTML: true})
	if err != nil {
		return nil, err
	}
	return e.Bytes(), nil
}
```

+ cachedTypeFields使用了sync.Map
```go
// ===== 1.9 =======
var fieldCache struct {
	value atomic.Value // map[reflect.Type][]field
	mu    sync.Mutex   // used only by writers
}

// cachedTypeFields is like typeFields but uses a cache to avoid repeated work.
func cachedTypeFields(t reflect.Type) []field {
	m, _ := fieldCache.value.Load().(map[reflect.Type][]field)
	f := m[t]
	if f != nil {
		return f
	}

	// Compute fields without lock.
	// Might duplicate effort but won't hold other computations back.
	f = typeFields(t)
	if f == nil {
		f = []field{}
	}

	fieldCache.mu.Lock()
	m, _ = fieldCache.value.Load().(map[reflect.Type][]field)
	newM := make(map[reflect.Type][]field, len(m)+1)
	for k, v := range m {
		newM[k] = v
	}
	newM[t] = f
	fieldCache.value.Store(newM)
	fieldCache.mu.Unlock()
	return f
}

// ===== 1.11 ========

var fieldCache sync.Map // map[reflect.Type][]field

// cachedTypeFields is like typeFields but uses a cache to avoid repeated work.
func cachedTypeFields(t reflect.Type) []field {
	if f, ok := fieldCache.Load(t); ok {
		return f.([]field)
	}
	f, _ := fieldCache.LoadOrStore(t, typeFields(t))
	return f.([]field)
}
```


#### 参考
https://draveness.me/golang/docs/part4-advanced/ch09-stdlib/golang-json/