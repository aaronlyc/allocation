## 1 go monkey介绍
Monkey是Golang的一个猴子补丁（`monkeypatching`）框架，在运行时通过汇编语句重写可执行文件，将待打桩函数或方法的实现跳转到桩实现，原理和热补丁类似。如果读者想进一步了解Monkey的工作原理，请阅读博客：`http://bouk.co/blog/monkey-patching-in-go/`。
通过Monkey，我们可以解决函数或方法的打桩问题，但Monkey**不是线程安全**的。
不要将Monkey用于并发的测试中。

## 2 安装
```bash
go get github.com/aaronlyc/monkey
```

## 3 使用场景
* 基本场景：为一个函数打桩
* 基本场景：为一个过程打桩
* 基本场景：为一个方法打桩
* 复合场景：由任意相同或不同的基本场景组合而成
* 特殊场景：桩中桩

#### 1 函数/过程打桩
Exec是一个操作函数，实现很简单，代码如下所示：
```go
func Exec(cmd string, args ...string) (string, error) {
    cmdpath, err := exec.LookPath(cmd)
    if err != nil {
        fmt.Errorf("exec.LookPath err: %v, cmd: %s", err, cmd)
        return "", infra.ErrExecLookPathFailed
    }

    var output []byte
    output, err = exec.Command(cmdpath, args...).CombinedOutput()
    if err != nil {
        fmt.Errorf("exec.Command.CombinedOutput err: %v, cmd: %s", err, cmd)
        return "", infra.ErrExecCombinedOutputFailed
    }
    fmt.Println("CMD[", cmdpath, "]ARGS[", args, "]OUT[", string(output), "]")
    return string(output), nil
}
```

Exec函数的实现中调用了库函数exec.LoopPath和exec.Command，因此Exec函数的返回值和运行时的底层环境密切相关。在UT中，如果被测函数调用了Exec函数，则应根据用例的场景对Exec函数打桩。
Monkey的API非常简单和直接，我们直接看打桩代码：

```go
import (
    "testing"
    . "github.com/smartystreets/goconvey/convey"
    . "github.com/bouk/monkey"
    "infra/osencap"
)

func TestExec(t *testing.T) {
    Convey("test has digit", t, func() {
        Convey("for succ", func() {
            guard := Patch(osencap.Exec, func(cmd string, args ...string) (string, error) {
                return outputExpect, nil
            })
            defer guard.Unpatch()
            output, err := osencap.Exec(any, any)
            So(output, ShouldEqual, outputExpect)
            So(err, ShouldBeNil)
        })
    })
}
```

Patch是Monkey提供给用户用于函数打桩的API，Patch的返回值是一个PatchGuard对象指针，主要用于在测试结束时删除当前的补丁，Patch函数的声明如下：

```go
func Patch(target, replacement interface{}) *PatchGuard
```

#### 2 方法打桩
假设数据库Etcd有一个方法Get，当用户输入团队id时，该方法将返回团队成员的名字列表：
```go
type Etcd struct {
}

func (e *Etcd) Get(id string) []string {
    names := make([]string, 0)
    ...
    return names
}
```
我们对Get方法的打桩代码如下：
```go
var e *Etcd
guard := PatchInstanceMethod(reflect.TypeOf(e), "Get", func(_ string) []string {
    return []string{"LiLei", "HanMeiMei", "ZhangMing"}
})
defer guard.Unpatch()
```
## 4 缺陷及应对方案
#### （1）inline函数打桩无效
inline函数不是在调用时发生控制转移，而是在编译时将函数体嵌入到每一个调用处，所以inline函数在调用时没有地址。inline函数没有地址的特性导致了Monkey框架的第一个缺陷：对inline函数打桩无效。
** 应对方案：**  通过命令行参数`-gcflags=-l`禁止inline（内联）
```bash
go test -gcflags=-l -v func_test.go -test.run TestIsEqual
```

#### （2）首字母小写不能打桩
反射机制的这种差异导致了Monkey框架的第二个缺陷：在go1.6版本中可以成功打桩的首字母小写的方法，当go版本升级后Monkey框架会显式触发panic，表示`unknown method`。
正交设计四原则告诉我们，要向稳定的方向依赖。首字母小写的方法或函数不是public的，仅在包内可见，不是一个稳定的依赖方向。如果在UT测试中对首字母小写的方法或函数打桩的话，会导致重构的成本比较大。
**应对方案：** 不管现在团队使用的go版本是哪一个，都不要对首字母小写的方法或函数打桩，不但可以确保测试用例在go版本升级前后的稳定性，而且能有效降低重构的成本。