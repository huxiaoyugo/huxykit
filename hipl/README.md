## (HIPL) HTTP interface protocol language
http接口协议语言, 用于定义http接口协议标记的语言，旨在简单清晰的表达接口定义的字段名和字段类型。

## 协议书写规则
### 数据类型
+ 基本类型(必须大写)：
	+ STRING, 
	+ INT: 被当做INT32处理，
	+ INT32, 
	+ INT64,
	+ UINT32,
	+ UINT64,
	+ FLOAT,
	+ DOUBLE,
	+ BOOL(或者BOOLEAN)
+ 复合类型：
	+ [ {{ARRAY_TYPE}} ]：数组，ARRAY_TYPE必须用{{}}包裹起来（兼容老版本）
	+ MAP<K_TYPE, V_TYPE>： MAP，K_TYPE,和V_TYPE不用{{}}包裹，
	+ ENUM<ENUM_TYPE>: 枚举，兼容原有逻辑，若只有ENUM没有在<>中指定具体的类型，按照STRING类型处理；
+ 自定义类型：
	+ 对象的字段一旦提交不得删除，不得修改
	+ 添加新的字段只能在最后添加，不得插入都已有字段中间
+ 订阅枚举类型：
	+ 字段的值必须为数字，且不能为0
	+ 数字不能重复

+ 类型和proto3类型对应关系

|类型           |proto3类型                   |
|:----------:  |:-----------------------------------------: | 
| STRING      | google.protobuf.StringValue | 
| INT32/INT     | google.protobuf.StringValue 
| UINT32    | google.protobuf.UInt32Value |
| INT64     | google.protobuf.Int64Value | 
| UINT64     | google.protobuf.UInt64Value |
| BOOL/BOOLEAN        | google.protobuf.BoolValue   | 
| FLOAT       | google.protobuf.FloatValue  | 
| DOUBLE      | google.protobuf.DoubleValue  |
| [{{TYPE}}]         | repeat type                     |
| MAP<K_TYPE,V_TYPE>  | map<K_TYPE, V_TYPE> | 
| STRUCT         |  message | 
| ENUM<TYPE>         | enum | 

### 对象的声明：
+ 所有的申明都放在\```self-json ```代码块中；
+ 代码块中非申明内容，可添加注释，使用// 或者/**/ ；
+ 对象或者枚举的声明、"request","resource"和字段的名称必须使用引号包裹起来；
+ 字段类型可以在最外层用引号包裹，也可以不用。例如{{STRING}}和"{{STRING}}"均可（兼容老版本）。
+ 除了[]所有的类型都需要使用{{}}包裹起来；
```json
// request声明
"request": {
    "url": "dueros://account/GetUserInfo",
    "params": {
		"name": {{STRING}},
		"extra_info": {{RequestParam}}
	}
}
// resource声明
"resource": {
    "header": {
        "namespace": "ai.dueros.resource.account",
        "name": "PhoneInfo",
        "version": "v1.0"
    },  
    "payload": {
        "phone": {{STRING}}
    }
}

// 自定义对象
struct "RequestParam":{
	"name": {{STRING}},
	"age": {{INT32}},
	"score": {{FLOAT}},
	"is_stu": {{BOOL}},
	"list": [ {{STRING}} ], // 数组字段
	"favorit": {{ENUM<Sports>}}, // 枚举类型的字段
	"ext_info": {{MAP<STRING, STRING>}},
	"school": {{School}}  // 自定义对象类型的字段
}

// 自定义对象
struct "School": {
	"name":{{STRING}},
	"age":{{INT32}}
}

// enum类型定义
enum "Sports": {
	"basketball":1,
	"football":2
}
```

### 如何转化？
+ 将所有的自定义对象 ==> proto message/enum；
+ request.params ==> proto message, 以url为依据命名；
+ resource.payload ==> proto message, 以head.namespace+header.name为依据命名；
+ 对于request.Params和resource.ResourcePayload 的oneof列表，采用生成临时文件和原有的文件进行merge，原因在于oneof列表的tag值一旦设置之后就不能改变, 所以需要读取原有的文件，保证已有的tag不变，新增的params，tag与原有的不同。
        


