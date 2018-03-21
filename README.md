# gopool
Common pool for goroutine and timer
# Usage of Timerpool
```
var GlobalTimerpool Timerpool
var timer=GlobalTimerpool.Get(1*time.Second)
```

# Usage of Grpool
```
var grpooler = gopool.NewGrpool(1024,1, 1)
grpooler.Schedule(func(){ do something  }())
```
