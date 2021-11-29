# bytesconv

string与[]byte的直接转换是通过底层数据copy实现的

```
var a = []byte("hello boy")
var b = string(a)
```

这种操作在并发量达到十万百万级别的时候会拖慢程序的处理速度，所以如果不修改数据，可以通过内存转换，避免底层数据拷贝的方式处理。
