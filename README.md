##### SDK
- 通常指封装了后端服务API接口的软件包
- 通常包含相关的文档、示例、封装好的API接口和工具

##### 目录结构
- README.md: 帮助文档
- examples/: 使用示例
- sdk/
1. SDK共享包
2. 封装最基础的通信功能
3. 如果是HTTP服务，基本都是基于net/http包进行封装
- services/{elmt,iam}
1. 把某类服务的API接口封装代码存放在services/<服务名>下

##### 设计方法
- 通过Config配置创建客户端Client: func NewClient(config sdk.Config) (Client, error)
- 创建的Client建议是interface类型
1. 将定义和具体实现解耦
2. Client的每一个方法对应一个API接口，如CreateUser、DeleteUser等