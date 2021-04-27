# easy-mail

这是一个旨在方便终端用户通过程序发送邮件的程序 (自动化程序的重要一环)

主要特点：

* 自动获取smtp服务器
* 可选保存发送方信息
* 可作为普通发件客户端

> thanks to gomail

## 使用说明

```bash
usage: easy-mail [-h|--help] [-f|--from "<value>"] [-t|--to "<value>"]
                 [-s|--subject "<value>"] [-c|--content "<value>"]
                 [--content-path "<value>"] [--content-type "<value>"]
                 [--attach "<value>" [--attach "<value>" ...]] [--smtp
                 "<value>"] [--password "<value>"] [-g|--generate] [-a|--auth
                 "<value>"] [-v|--version]

                 easily send mail from command line

Arguments:

  -h  --help          Print help information
  -f  --from          email send from
  -t  --to            recv address list, separated by ','
  -s  --subject       email title
  -c  --content       simple email content
      --content-path  email content path
      --content-type  email content type
      --attach        attach file path list
      --smtp          manually set smtp address like: smtp.abc.com:587 it can
                      be auto find if not set
      --password      email password
  -g  --generate      generate auth file to simple use
  -a  --auth          auth file path
  -v  --version       show version of easy-mail
```

#### 1. 选择保存发件账户信息

> 如果在特定主机发送邮件，方便以后使用 (推荐优先设定)

```bash
easy-mail --from xx@a.b --password you-password --smtp smtp.a.b:587 -g
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

* 若需要同时发送给多个邮箱，则可用逗号分隔不同邮箱账户，如：

```bash
easy-mail -t first@b.c,second@d.e
```

* 邮件正文可指定内容媒体类型，如 `text/html`

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

若需要发送多个附件，则再接 `--attach` 参数即可，如：

```bash
easy-mail -t a@b.c --attach /path/attach1 --attach /path/attach2
```

