Box:
----
 分布式的服务端编程框架，由grpc提供rpc支持。
 设计有gateway层，nameserver层，worker层。
 3者都有监听端口，方便相互间通信，实现简单的Ping/Pong协议测定RTT。
 支持带灰度可用。
 暂时只支持IPv4。


关于grpc
===

- 基于http2
- 基于protobuf3
- 默认使用TLS

worker
===
- 定时地向nameServer报活；启动时自动加集群。
- worker包括但不限：告诉ns上、下线，版本(灰度时用)、负载压力、连接数等。
- 业务、逻辑处理、服务进程。
- 接收gateway层转发过来的请求，处理完后把结果交由gateway转发到目标客户端。

gateway
===
- 从nameserver上获取可用worker列表。
- 所监听端口为对外（客户端）入口，对外实现tcp传输协议（暂时只有htttp2）。
- 负债数据包传发到worker上处理，网络流量必需流经此应用。
- 对新请求按woker压力分配；对旧请求（必需带workerId）分配到指定worker上。
- 属于7层负载均衡器，不带业务逻辑，无状态。更应该部署多实例。
- 可以由lvs或DNS(单域名指向多IP)进行多实例部署。


nameserver
===
- master/slave 方式运行解决单点问题；轻量级服务进程。
- gateway与worker必需知道nameserver存在。
- worker列表发生变化时告诉所有gateway。
- workerId的设定：ip:port-->int64；高32位为ip转成的int，低32位为port。
- workerId这样设定的原因是使客户端有机会可以指定worker进行调试。
- 当workerid冲突时，使用随机生成的int32；32位整数，便知是rand生成。


client
===
- 可以是任意客户端，不限平台、不限开发语言。
- 使用tcp方式与gateway通信。
- 实现gateway所支持网络协议即可（暂时只有htttp2）