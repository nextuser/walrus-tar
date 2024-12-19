package tar

import (
	"fmt"
	"log"
	"os"
)

var enableDebug = false

func debug(a ...any) {
	if enableDebug {
		fmt.Println(a...)
	}
}

// 定义一个用来打印的函数，少写点代码，因为要处理很多次的 err
// 后面其他示例还会继续使用这个函数，就不单独再写，望看到此函数了解
func ErrPrintln(err error) {
	if err != nil {
		log.Println("error:", err)
		os.Exit(1)
	}
}
