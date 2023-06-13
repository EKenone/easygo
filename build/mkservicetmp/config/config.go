package config

var Conf = &Config{}

type Config struct {
	ModName  string   `yaml:"mod_name"` //go.mod的module名字
	Proto    Proto    `yaml:"proto"`    //协议
	Internal Internal `yaml:"internal"` //私有代码存放的参数
	Doc      Doc      `yaml:"doc"`      //文档存放路径
}

type Proto struct {
	Source   string `yaml:"source"`   //协议源路径
	Analysis string `yaml:"analysis"` //协议源解析后存放路径
}

type Internal struct {
	Entrance string `yaml:"entrance"` //请求入口文件夹
	Service  string `yaml:"service"`  //业务代码文件夹
	Route    string `yaml:"route"`    //路由启动存放文件，成功的路由会在这里统一启动，请在服务启动的时候调用InitRoute方法
}

type Doc struct {
	Swagger string `yaml:"swagger"` //swagger存放路径
}

type TestPath struct {
	Unit    string `yaml:"unit"`     //单元测试存放路径
	HttpApi string `yaml:"http_api"` //ide web api测试接口存放路径
}
