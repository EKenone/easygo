package comply

import (
	"encoding/json"
	"io"
	"os"
)

type Parse struct {
	Group      Group       `json:"group"`      //接口组名称
	Interfaces []Interface `json:"interfaces"` //接口
}

type Group struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type Interface struct {
	Name           string      `json:"name"`             //接口名称
	Desc           string      `json:"desc"`             //接口备注
	Method         string      `json:"method"`           //请求方法
	NotAuth        bool        `json:"not_auth"`         //是否不检验身份
	LastPath       string      `json:"last_path"`        //最后一节路由
	MidType        interface{} `json:"mid_type"`         //中间件类型
	ReqContentType string      `json:"req_content_type"` //接参数形式
	ResContentType string      `json:"res_content_type"` //接参数形式
	Uri            []Field     `json:"uri"`              //uri参数
	Req            Message     `json:"req"`              //请求参数
	Res            Message     `json:"res"`              //返回数据
}

type Field struct {
	Name    string   `json:"name"`     //字段名
	Class   string   `json:"class"`    //字段类型
	Desc    string   `json:"desc"`     //字段备注
	Binding string   `json:"binding"`  //字段校验
	From    string   `json:"from"`     //字段来源
	JsonTag []string `json:"json_tag"` //json tag的补充
}

type Message struct {
	Name    string    `json:"name"`    //结构体名称
	Fields  []Field   `json:"fields"`  //结构体参数
	Message []Message `json:"message"` //请求参数字段结构体
}

func (p *Parse) Parse(path string) (err error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return
	}

	body, err := io.ReadAll(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, p)
	if err != nil {
		return
	}
	return
}

type Doc struct {
	Title   string   `yaml:"title"`
	Route   string   `yaml:"route"`
	Desc    string   `yaml:"desc"`
	Schemes []string `yaml:"schemes"`
	Host    string   `yaml:"host"`
	Ver     string   `yaml:"ver"`
	Auth    struct {
		Security string `yaml:"security"`
		Title    string `yaml:"title"`
		In       string `yaml:"in"`
		Name     string `yaml:"name"`
		Token    string `yaml:"token"`
	} `yaml:"auth"`
	Contact struct {
		Name  string `yaml:"name"`
		Url   string `yaml:"url"`
		Email string `yaml:"email"`
	} `yaml:"contact"`
}
