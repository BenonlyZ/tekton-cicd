1.使用 kubectl apply 命令安装最新版本的 Tekton Triggers 及其依赖项：
kubectl apply --filename https://storage.googleapis.com/tekton-releases/triggers/latest/release.yaml
kubectl apply --filename https://storage.googleapis.com/tekton-releases/triggers/latest/interceptors.yaml
注意两个都需要安装，interceptors是拦截器
如果要按版本安装，则执行：
kubectl apply -f https://storage.googleapis.com/tekton-releases/triggers/previous/v0.1.0/release.yaml

2.如果遇到镜像拉取失败，则需要把上面images替换成以下
tektondev/trigger-eventlistenersink
tektondev/triggers-controller
tektondev/trigger-webhook
tektondev/trigger-interceptors

3.安装完毕后，会在tekton-pipeline namespace下看到以下pod：
tekton-triggers-controller-*
tekton-triggers-core-interceptors-*
tekton-triggers-webhook-76874f5f58-*

4.部署完后，会生成名称为el-gitlab-listener的svc，需要把它对外暴露，才能在github上配置tekton的webhook。
这里用内网穿透的方法，当然你也可以直接用ingress等
先执行：
kubectl port-forward -n default service/el-gitlab-listener 8080:8080
映射本地8080端口
然后，安装ngrok
安装完后执行:  ngrok http 8080
ngrok会生成一个指向你本地8080端口的域名，用这个域名配置github的webhook即可
