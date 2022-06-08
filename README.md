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

- 主要优点
1. 大量使用了Go interface特性，将接口的定义和实现解耦，可以支持多种实现方式
2. 接口调用层级跟资源的层级相匹配，调用方式更加友好
3. 多版本共存
```
// 项目客户端 -| 应用客户端 -| 服务客户端 -| 资源名 -| 接口
clientset.Iam().AuthzV1().Authz().Authorize() // 调用了/v1/authz版本的API接口
clientset.Iam().AuthzV2().Authz().Authorize() // 调用了/v2/authz版本的API接口
```

##### 客户端类型
- RESTClient
1. Raw类型的客户端 / 整个SDK的核心
2. 向下通过调用Request模块，来完成HTTP请求方法、请求路径、请求体、认证信息的构建
3. 向上提供Post()、Put()、Get()、Delete()等方法来供客户端完成HTTP请求
4. 可以通过指定HTTP的请求方法、请求路径、请求参数等信息，直接发送HTTP请求: client.Get().AbsPath("/version").Do().Into()

- 基于RESTClient封装的客户端
1. 例如AuthzV1Client、APIV1Client等
2. 执行特定REST资源、特定API接口的请求，方便开发者调用

##### 接口层级
```
// 项目级别的接口
type Interface interface {
    Iam() iam.IamInterface
    Tms() tms.TmsInterface
}

// 应用级别的接口
type IamInterface interface {
    APIV1() apiv1.APIV1Interface
    AuthzV1() authzv1.AuthzV1Interface
}

// 服务级别的接口
type APIV1Interface interface {
    RESTClient() rest.Interface
    SecretsGetter
    UsersGetter
    PoliciesGetter
}

// 资源级别的客户端
type SecretsGetter interface {
    Secrets() SecretInterface
}

// 资源的接口定义
type SecretInterface interface {
    Create(ctx context.Context, secret *v1.Secret, opts metav1.CreateOptions) (*v1.Secret, error)
    Update(ctx context.Context, secret *v1.Secret, opts metav1.UpdateOptions) (*v1.Secret, error)
    Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
    DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
    Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Secret, error)
    List(ctx context.Context, opts metav1.ListOptions) (*v1.SecretList, error)
    SecretExpansion
}
```