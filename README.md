<img style="width:50px" src="./assets/static/img/markless.png" />

Markless
===

一个非常简单的书签管理器，仅收集，分类，管理书签。

* 非常轻量简单
* 部署友好，仅包含一个可执行文件
* 使用`Sqlite`
* 支持黑暗模式
* 可通过浏览器插件收集网页
* 可通过IOS快捷指令收集网页
* 支持多用户
* 支持打个多标签
* 

![](./example/index.png)

[demo地址](https://wsh233.cn/webapp/markless) 用户名：`demo` 密码：`demo1234`

使用
===

查看启动命令参数

```bash
markless -h
```

参数说明

* -baseurl: 部署根路由， 默认 `/`
* -databaseurl：sqlite路径 ，默认自动创建在程序所在目录
* -port：端口 ，默认`5000`
* -title：网站名称，默认`markless`
* -adminname：默认管理员名称，默认`admin`
* -password：默认管理员密码，默认`admin1234`



待完成
===

* 导入解析浏览器导出的书签

* 导出书签（json格式，带标签）
* 快照（保存某一时刻的网页内容，防止链接失效）



由下面两个开源项目启发而成编写：

* [linkding](https://github.com/sissbruecker/linkding)

* [miniflux](https://github.com/miniflux/v2)
