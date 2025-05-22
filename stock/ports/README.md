# ports 文件夹

由于 protobuf 无法像 openapi 一样生成客户端和服务端分离的代码，所以只能将生成的代码放在 common 中供客户端和服务端一起使用。

因此 protobuf 生成的代码的服务端部分无法存在于 ports 用于与外界交互的文件夹中。
