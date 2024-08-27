## 说明

仓库来源: [docker-flarum](https://github.com/crazy-max/docker-flarum)

由于原仓库没有安装第三方插件的功能，所以修改了源码，添加了插件安装功能。

## 基本的思想

- 就是往 Dockerfile 中添加了插件安装的步骤，然后重新构建镜像。

- 基于 github action 构建,地址: [zoushucai/docker-flarum](https://github.com/zoushucai/docker-flarum)
