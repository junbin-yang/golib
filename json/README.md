# json

在golang的一些低版本中原生`encoding/json`库序列化的效率较低，但在高版本的golang中原生json库已经优化。两者性能已经差别不大，但原生库还无法满足需要指定tag的场景。使用号称最快的jsoniter包可直接兼容替换。

