# 随机图片 api -- go 实现

- 采用的是 go 来实现的, 采用 gin 框架, 处理图像采用 bild 包,只是使用了调整大小的功能. ~~原本使用 bimg 包, 但是依赖太严重了(因此不能交叉编译)~~

- 图片来源于互联网, 主要来自 c 站

## 使用方法 (推荐 docker)

由于基于: `scratch` 镜像, 比较简洁,无法进入容器, 需要的话可以改为 `alpine`

```yaml
services:
  randomimages:
    container_name: randomimages-server
    image: qqzsc/randomimages:latest
    # image: ghcr.io/zoushucai/randomimages:latest
    restart: always
    ports:
      - "9113:9113"
    volumes:
      # 需要准好图片资源文件夹
      - ./../assets:/app/assets
      - ./logs:/app/logs
```

然后运行 `sudo docker-compose up` 就可以通过 `http://localhost:9113/v1/random` 访问了

## 参数说明

- 假设: `url = http://localhost:9113/v1/random`

- 访问 `url?format=xxx&sub=xxx&width=xxx&index=xxx` 即可获取指定的图片,(图片的类型包括: `jpeg, png, gif`)

  - `sub`: 文件夹名称, 例如 `url?sub=cat` 即可获取 `assets/cat` 文件夹下的图片
  - `width`: 图片宽度, 默认为 0
  - `height`: 图片高度, 默认为 0

    - 如果 `width` 和 `height` 都为 0, 则默认为原始图片的宽度
    - 如果二者只要一个为 0, 则调整为指定的宽高(保持比例)
    - 如果都大于 0, 则调整到指定的宽高(拉伸)

  - `index`: 索引, 例如 `url?index=1` 即可获取索引为 1 的图片(按顺序), 如果未指定, 则默认为随机获取

  - 上述选项是可以组合使用的

- 最好,使用域名+反向代理的方式来访问, 不然可能会有跨域问题

```bash

# example -- 直接返回图片
http://localhost:9113/v1/random?sub=dongman&index=0


# 如果要返回 json,
curl -H "Accept: application/json" http://localhost:9113/v1/random

# 返回
{"file":"localhost:9113/assets/renwu/00059-3052920575-masterpiece.png","message":"success"}
# 这里的 file 是图片的 url, 直接把 assets 目录下的所有图片进行了挂载, 因此可以通过这个地址直接访问图片
# message 是返回的提示信息


# 还可以通过 swagger 来查看接口文档(已删除, 这个会使得编译后的体积大 3 倍,算了,删除)
# http://localhost:9113/swagger/index.html

```

## 编译安装

- 1. 首先,准备好 go 环境, 以及必要的依赖.

```bash
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
sudo tar -xzvf go1.23.0.linux-amd64.tar.gz -C /usr/local
echo 'export PATH="$PATH:/usr/local/go/bin"' >> ~/.zshrc
source ~/.zshrc
echo 'export PATH=$(go env GOPATH)/bin:$PATH' >> ~/.zshrc
source ~/.zshrc
rm go1.23.0.linux-amd64.tar.gz

```

- 2. 准备好图片资源, 把图片放到 `assets` 目录下, 如果要对图片进行分类, 可以在 `assets` 目录下新建文件夹, 把图片放到对应的文件夹下. 例如 `assets/cat` 文件夹下的图片, 就是猫的图片.

- 运行 `go mod init mygo` 来初始化项目

- 运行 `go mod tidy` 来更新依赖

- 运行 `go run main.go` 来启动服务

- 访问 `http://localhost:9113/v1/random` 即可随机获取一张图片

- 如果访问成功, 则可以 `go build -o main-mac-arm64 main.go` 来编译成二进制文件, 然后通过 `./main-mac-arm64` 来启动服务

```bash
# 编译的时候添加 -ldflags="-s -w" 可以减少二进制文件大小
go build -ldflags="-s -w" -o main-mac-arm64 main.go


# 添加参数, 指定端口和图片路径
./main-mac-arm64 --port=9113 --imgdir=imgs

```

### 关于 go 的 交叉编译

不想折腾, 推荐第三方程序: [goreleaser](https://github.com/goreleaser/goreleaser)

- 注意: 如果项目底层依赖 c 库,则可能交叉编译的时候会失败

## 关于 go 的热重载

- [air](https://github.com/air-verse/air)

## 关于 api 文档 -- swagger

- [swag](https://github.com/swaggo/swag)
- [gin-swagger](https://github.com/swaggo/gin-swagger)

## 利用 action 进行自动编译

- 参考仓库的 `.github/workflows` 文件
