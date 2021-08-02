## 1 什么是ginkgo
ginkgo是一个用go写的BDD(Behavior Driven Development)的测试框架，一般用于Go服务的集成测试。

ginkgo的特点
BDD的代码风格

```golang
Describe("delete app api",func(){
    It("should delete app permanently",func(){...})
    It("should delete app failed if services existed", func(){...})
）
```

（1）Ginkgo定义的DSL语法(Describe/Context/It)可以非常方便的帮助大家组织和编排测试用例。
（2）在BDD模式中，测试用例的标题书写，要非常注意表达，要能清晰的指明用例测试的业务场景。只有这样才能极大的增强用例的可读性，降低使用和维护的心智负担。
（3）可读性这一点，在自动化测试用例设计原则上，非常重要。因为测试用例不同于一般意义上的程序，它在绝大部分场景下，看起来都像是一段段独立的方法，每个方法背后隐藏的业务逻辑也是细小的，不具通识性。这个问题在用例量少的情况下，还不明显。但当用例数量上到一定量级，你会发现，如果能快速理解用例到底是能做什么的，真的非常重要。而这正是BDD能补足的地方。
（4）不过还是要强调，Ginkgo只是提供对BDD模式的支持，你的用例最终呈现的效果，还是依赖你自己的书写。

## 2 ginkgo使用
#### （1）ginkgo安装

```go
$ go get github.com/onsi/ginkgo/ginkgo
$ go get github.com/onsi/gomega/...
```

#### （2）ginkgo常用模块

常用的10个：It、Context、Describe、BeforeEach、AfterEach、JustBeforeEach、BeforeSuite、AfterSuite、By、Fail

**It模块:**

It是测试例的基本单位，即It包含的代码就算一个测试用例， It可以理解维测试代码的最小单元。

**Context和Describe 模块**

Context和Describe的功能都是将一个或多个测试例归类个describe可以包含多个context，一个context可以包含多个IT模块

**BeforeEach模块**

BeforeEach是每个测试例执行前执行该段代码。比如创建数据库连接就可以使用BeforeEach ，每个BeforeEach只在当前域内起作用。执行顺序是同一层级的顺序执行，不同层级的从外层到里层以此执行。类似与 全局变量和局部变量的区别

**AfterEach模块**

AfterEach是每个测试例执行后执行该段代码比如销毁数据库连接，一般用于测试例执行完成后进行数据清理，也可以用于结果判断。

**JustBeforeEach**

JustBeforeEach是在BeforeEach执行之后，测试例执行之前执行测试用例执行前的一些前置操作。

**BeforeSuite模块**

`BeforeSuite` 是在该测试集执行前执行，即该文件夹内的测试例执行之前`BeforeSuite`和AfterSuite写在` _suite_test.go` 文件中，会在所有测试例执行之前和之后执行如果`BeforeSuite`执行失败，则这个测试集都不会被执行。

> Tip：使用`ctrl+C`中断执行时，`AfterSuite`仍然会被执行，需要再使用一次`ctrl+C`中断

**AfterSuite模块**

AfterSuite是在该测试集执行后执行，即该文件夹内的测试例执行完后。

**By模块**

By是打印信息，内容只能是字符串，只会在测试例失败后打印，一般用于调试和定位问题

**Fail模块**

Fail是标志该测试例运行结果为失败，并打印里面的信息

**Specify模块**

Specify和It功能完全一样，It属于其简写

#### （3）ginkgo的三个标志
F、X和P，可以用在Describe、Context、It等任何包含测试例的模块

F含义Focus，使用后表示只执行该模块包含的测试
Tip：当里层和外层都存在Focus时，外层的无效，即下面代码只会执行B测试用例

P的含义是Pending，即不执行，用法和F一样，规则的外层的生效

X和P的含义一样
还有一个跳过测试例的方式是在代码中加Skip

```go
It("should do something, if it can", func() {
    if !someCondition {
        Skip("special condition wasn't met")
    }

    // assertions go here
})
```

#### （4）ginkgo 并发设置

ginkgo -p 使用默认并发数
ginkgo -nodes=N 自己设定并发数

默认并发数是用的参数runtime.NumCPU()值，即逻辑CPU个数，大于4时，用runtime.NumCPU()-1

#### （5）ginkgo显示实时日志
如果需要显示实时日志，需要添加 -stream参数

并发执行时打印的日志是汇总后经过合并处理再打印的，所以看起来比较规范，
每个测试例的内容也都打印在一起，但时不实时，如果需要实时打印，加-stream参数，缺点是每个测试例日志交叉打印

#### （6）ginkgo goroutine设置
在平时的代码中，我们经常会看到需要做异步处理的测试用例。但是这块的逻辑如果处理不好，用例可能会因为死锁或者未设置超时时间而异常卡住，非常的恼人。好在Ginkgo专门提供了原生的异步支持，能大大降低此类问题的风险。类似用法：

```go
 It("should post to the channel, eventually", func(done Done) {
    c := make(chan string, 0)

    go DoSomething(c)
    Expect(<-c).To(ContainSubstring("Done!"))
    close(done)
}, 0.2)
```
0.2的单位是秒

#### （7）ginkgo性能测试
使用`Measure`这个模块

```go
Measure("it should do something hard efficiently", func(b Benchmarker) {
    runtime := b.Time("runtime", func() {
        output := SomethingHard()
        Expect(output).To(Equal(17))
    })

    Ω(runtime.Seconds()).Should(BeNumerically("<", 0.2), "SomethingHard() shouldn't take too long.")
    
    b.RecordValue("disk usage (in MB)", HowMuchDiskSpaceDidYouUse())
}, 10)
```

#### （8）ginkgo命令行使用

打开电脑终点是使用 ginkgo help即可
比较常用的命令

```bash
ginkgo bootstrap <FLAGS>  创建测试关联文件
ginkgo generate <filename(s)> 生成测试代码模版
ginkgo nodot  更新测试套件中的nodot声明
ginkgo convert /path/to/package  转变包文件格式为ginkgo可识别的格式
ginkgo unfocus (or ginkgo blur) 递归地取消当前目录下所有集中测试的焦点
ginkgo version   输出版本信息
```

## 3 ginkgo初级实战

1. 本地gopath/src下创建一个项目books
2. 创建 books.go 文件

```golang
package books

import (
	"encoding/json"
	"fmt"
)

type Book struct  {
	Title  string
	Author string
	Pages  int
}

func (b Book) AuthorLastName() (interface{}, interface{}) {
    fmt.Println("AuthorLastName")
	return b.Author,b.Title
}
func (b Book) CategoryByLength() string{
	if b.Pages >300{
		return "NOVEL"
	} else if b.Pages<100 {
		return "SMALL STORY"
	} else {
		return "SHORT STORY"
	}
}
func NewBookFromJSON(json json.Decoder) error {
	book:=json.Decode(json)
	return book
}
```

3. go mod init
当前文件夹下会生成两个文件 `go.mod  go.sum`

4. 创建测试关联文件
进入项目当前路径下执行 ginkgo bootstrap 项目下会创建一个 `books_suite_test.go`文件，为测试用例入口

```go
package books_test

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "testing"
)

func TestBooks(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Books Suite")
}
```

5. 创建测试模版文件 ginkgo generate books 当前的包名字叫 books所以 创建出来的测试模版文件叫`books_test.go`

```go
package books_test

import (
    . "/path/to/books"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Book", func() {

})
```

这是一个通用的测试模版，我们需要根据自己的测试需求编写对应的测试规格
修改后的如下

```go
package books_test

import (
	. "books"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Book", func() {
	var (
		longBook  Book
		shortBook Book
		smallBook Book
	)

	BeforeEach(func() {
		longBook = Book{
			Title:  "Les Miserables",
			Author: "Victor Hugo",
			Pages:  1488,
		}
	
		shortBook = Book{
			Title:  "Fox In Socks",
			Author: "Dr. Seuss",
			Pages:  240,
		}
		smallBook = Book{
			Title:  "go program",
			Author: "caochunhui",
			Pages:  20,
		}
	})


	Describe("Categorizing book length", func() {
		Context("With more than 300 pages", func() {
			It("should be a novel", func() {
				Expect(longBook.CategoryByLength()).To(Equal("NOVEL"))
			})
		})
	
		Context("With fewer than 300 pages", func() {
			It("should be a short story", func() {
				Expect(shortBook.CategoryByLength()).To(Equal("SHORT STORY"))
			})
		})
		Context("With fewer than 100 pages", func() {
			It("should be a small story", func() {
				Expect(smallBook.CategoryByLength()).To(Equal("SMALL STORY))
			})
		})
	})
})
```

当前这个测试代码中我们写了三个用例，分别测试books源代码中CategoryByLength功能所对应的三个逻辑，即根据书页的多少判断书本的类型，测试代码只需要覆盖所有的逻辑即可。

6. 执行测试代码
这块有三种执行办法
```bash
go test -v
ginkgo
还可以 ginko build 编译成二进制，执行二进制
```