package constant

//项目状态
const (
	ProjectStatusEnable  int = iota + 1 //启用
	ProjectStatusDisable                //禁用
	ProjectStatusCancel                 //注销
)

// 是否删除
const (
	// 是
	IsDelete = 1
	// 否
	IsNotDelete = 2
)
