

<p align="center"><img style="width:50px" src="./web/assets/static/img/markless.png" /></p>

<p align="center"><a href="./README.md">English</a>｜<a href="./README-zh.md">中文</a></p>


Markless
===

一个非常简单的书签管理器，仅收集，分类，管理书签。

* 简洁轻量
* 部署友好，仅包含一个可执行文件
* 使用`Sqlite` 



功能：

- [x] 支持黑暗模式
- [x] 通过浏览器插件/IOS快捷指令收集网页
- [x] 支持多用户
- [x] 支持未读已读，贴标签分类
- [x] 支持多语言
- [x] 支持快照，保存链接某一时刻的内容，防止链接失效。
- [ ] 导入解析浏览器导出的书签
- [ ] 导出书签（json格式，带标签）



![首页](./example/index.png)

👀 [demo 地址](https://wsh233.cn/webapp/markless)  用户名：`demo` 密码：`demo1234`

使用
===

查看启动命令参数

```bash
markless -h
```



**通过浏览器插件收集网页**

浏览器的代码也是开源，源代码就在`crx`文件里，[下载](./example/markless-chrome-extension.crx)后，解压安装，

![浏览器插件](./example/broswer-extension.png)

启动实例，然后生成`密钥`

![浏览器插件](./example/token.png)

打开浏览器插件选项，复制后粘贴`实例地址和`密钥`

![浏览器插件](./example/broswer-extension-setting.png)



然后就能通过收集网页了，*链接*是必填项，其他如*标题*，*描述*等信息不填写的话程序会自动解析。

![浏览器插件](./example/broswer-collect.png)

**通过IOS快捷指令收藏网页**

[下载](./example/Markless.shortcut)快捷指令，修改并填写`url`为实例地址，请求头部里的`X-Token`填入前面安装浏览器插件获取的token值。

<p align="center"><img style="width:15em" src="./example/ios-shotcut.PNG" /></p>

双击快捷指令后，进入详细信息，开启`在共享表单中显示`



<p align="center"><img style="width:15em" src="./example/enable-share.PNG" /></p>



<p align="center"><img style="width:15em" src="./example/ios-collect.PNG" /></p>

`Safari`里分享网页，就能看见Markless了，点击即可收藏该网页。






感谢
===

由下面两个开源项目启发而成：

* [linkding](https://github.com/sissbruecker/linkding)

* [miniflux](https://github.com/miniflux/v2)

项目用到许多开源包，感谢作者❤️

