# mailall

这是一个旨在方便终端用户通过程序发送邮件的程序 (自动化程序的重要一环)

主要特点：

* 自动获取smtp服务器
* 可选保存发送方信息
* 可作为普通发件客户端

> thank to gomail

## 使用说明

```bash
usage: maillall [-h|--help] [--from "<value>"] [--to "<value>" [--to "<value>"
                ...]] [-s|--subject "<value>"] [-c|--content "<value>"]
                [--content-path "<value>"] [--content-type "<value>"] [--attach
                "<value>" [--attach "<value>" ...]] [--smtp "<value>"]
                [--password "<value>"] [-g|--generate] [-a|--auth "<value>"]

                send mail from command line

Arguments:

  -h  --help          Print help information
      --from          email send from
      --to            recv address list
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
```

#### 1. 选择保存发件账户信息

> 如果在特定主机发送邮件，方便以后使用

```bash
mailall --from xx@a.b --password you-password --smtp smtp.a.b:587 -g
```

执行后将在用户目录下将用户信息存放于 `.mailall.auth` ，其中包含用户名，密码以及指定的smtp服务器

__用户信息以明文形式存储，需要注意__

若不通过`--smtp`指定 smtp 服务器，则程序会通过发件人邮箱查询 `MX` 记录，并以 `25` 端口作为 smtp 服务器 (发现大多数通过DNS可查询到的`MX` 服务器只支持 25 端口，即未加密端口，安全起见，还是__手动指定支持ssl的smtp服务器__)

