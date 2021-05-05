# easy-mail

这是一个旨在方便终端用户通过程序发送邮件的程序 (自动化程序的重要一环)

主要特点：

* 自动获取smtp服务器
* 可选保存发送方信息
* 可作为普通发件客户端

场景:

* 自动化发布非结构化消息
* 手动执行，发送服务器文件数据到指定用户

> thanks to gomail

## 安装方式

* 如果有golang环境，可以通过 `go get` 安装

```bash
go get github.com/hellflame/easy-mail
```

* 也可以下载静态编译的可执行文件

从这里下载最新: [release](https://github.com/hellflame/easy-mail/releases)

将对应操作系统的可执行文件下载下来，解压或直接可用

```bash
unzip easy-mail-linux.zip
# or mac os
# unzip easy-mail-darwin.zip

# 直接执行
./easy-mail

# 拷贝至某一个可执行文件查询路径中
mv easy-mail /usr/local/bin
```

## 使用说明

```bash
usage: easy-mail [-h] [-f FROM] [-t TO [TO ...]] [-s SUBJECT] [-c CONTENT] [--content-path PATH] [--content-type TYPE] [--attach PATH [PATH ...]] [--smtp SMTP] [-p PASSWORD] [-g] [-a PATH] [-v]

easily send mail from command line

optional arguments:
  -h, --help                        show this help message
  -f FROM, --from FROM              email send from
  -t TO, --to TO                    recv address list
  -s SUBJECT, --subject SUBJECT     email subject
  -c CONTENT, --content CONTENT     email content
  --content-path PATH               email content path
  --content-type TYPE               email content type
  --attach PATH                     attach file path list
  --smtp SMTP                       manually set smtp address like: smtp.abc.com:465 it can be auto find if not set
  -p PASSWORD, --password PASSWORD  email password
  -g, --generate                    save auth to file for simple use
  -a PATH, --auth PATH              auth file path
  -v, --version                     show version of easy-mail

more info @ https://github.com/hellflame/easy-mail
```

#### 1. 选择保存发件账户信息

> 如果在特定主机发送邮件，方便以后使用 (推荐优先设定)

```bash
easy-mail -f xx@a.b -p you-password --smtp smtp.a.b:587 -g
```

执行后将在用户目录下将用户信息存放于 `.easy-mail.cred` ，其中包含用户名，密码以及指定的smtp服务器

__用户信息以明文形式存储，需要注意__

若不通过`--smtp`指定 smtp 服务器，则程序会通过发件人邮箱查询 `MX` 记录，或者以 `smtp.a.b` 类似的主机名作为smtp服务器，并与 `465, 25, 587` 这几个端口依次组合，尝试发件，可能存在响应较慢甚至发送失败的情况，此时可以手动指定官方给出的 smtp 服务器

> 以下命令均假设已保存用户信息，否则需要在每行命令中指定 发件账户，密码，smtp服务器参数

#### 2. 发送简单文本邮件

```bash
easy-mail -t a@b.c -s 'this is a simple title' -c 'see you tommorow'
```

以上以 `this is a simple title` 作为主题，向邮箱 `a@b.c` 投递消息，正文为 : `see you tommorow`

* 若需要同时发送给多个邮箱，使用空格隔开即可，命令如下：

```bash
easy-mail -t user1@b.c user2@d.e
```

* 邮件正文可指定内容媒体类型，如 `text/html` , 默认为 `text/plain`

```bash
easy-mail -t a@b.c -s 'this is a simple title' -c '<h1>see you <span style="color: red">tommorow</span></h1>' --content-type text/html
```

* 邮件可从文件读取正文内容

```bash
easy-mail -t a@b.c -s 'this is a simple title' --content-path /path/to/file
```

> 此时正文媒体类型依然需要手动指定，默认为 text/plain

#### 3. 发送包含附件的邮件

```bash
easy-mail -t a@b.c -s 'with attaches' -c 'you need to see detail form attaches' --attach /path/attach
```

可发送多个附件，如：

```bash
easy-mail -t a@b.c --attach /path/attach1 /path/attach2
```



