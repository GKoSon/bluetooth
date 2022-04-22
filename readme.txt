第一步 fork代码  https://github.com/tinygo-org/bluetooth.git
第二步 克隆到本地 $ git clone git@github.com:GKoSon/bluetooth.git
第三步 树莓派单独编译这个单例代码【来源是\bluetooth\examples\nusclient\main.go】
第四步 前面单例运行是OK的 
第五步 修改代码 推上去 git push
1--增加rm脚本
2--增加函数
第五步 重复第三步 在代码里面调用我add的函数 编译失败！预计中的
第六步 修改本地文件的 mod

它原始mod是这样
module x

go 1.17

require tinygo.org/x/bluetooth v0.4.0

require (
        github.com/JuulLabs-OSS/cbgo v0.0.2 // indirect
        github.com/fatih/structs v1.1.0 // indirect
        github.com/go-ole/go-ole v1.2.4 // indirect
        github.com/godbus/dbus/v5 v5.0.3 // indirect
        github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
        github.com/muka/go-bluetooth v0.0.0-20210812063148-b6c83362e27d // indirect
        github.com/sirupsen/logrus v1.6.0 // indirect
        golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
        golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6 // indirect
)
修改一下:

replace tinygo.org/x/bluetooth v0.4.0 => github.com/GKoSon/bluetooth a2c99f7
【没有版本号 就写刚刚推上去的commit号 】


测试编译通过了！

从此 我调用TinyGo的包 被转移了
不是转移在本地
而是转移在我自己的开源仓库里面。

你可以把tinyGo \bluetooth\examples\nusclient\main.go 单独放在PI运行
你需要使用我的方法 就使用mod管理 在把仓库replace

缺点时:必须要mod 不然没有人帮你转弯

附上测试代码: examples/nusclient/main.go 已经修改 自己去修改mod文件


