
// request声明
"request": {
    "url": "dueros://account/GetUserInfo",
    "params": {
		"name": "{{STRING}}",
		"extra_info": "{{RequestParam}}"
	}
}

// 自定义对象
struct "RequestParam": {
	"name": {{STRING}},
	"age": {{INT32}},   
	"score": {{FLOAT}},
	"is_stu": "BOOL"
	"list": [ STRING ], // 数组字段
	"favorit": "{{ENUM<Sports>}}", // 枚举类型的字段
	"ext_info": "{{MAP<STRING, STRING>}}",
	"school": "{{School}}"  // 自定义对象类型的字段
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