package proto

var (
	// basic
	StdSuccess  = &Base{200, "succes"}
	BadRquest   = &Base{400, "参数错误"}
	InternalErr = &Base{500, "服务器内部错误"}

	// user
	BadUserName      = &Base{1000, "非法的用户名"}
	WrongPassword    = &Base{1001, "密码错误"}
	UserNotFound     = &Base{1002, "该用户不存在"}
	UserAlreadyExist = &Base{1003, "用户已存在"}
	BadPassword      = &Base{1004, "密码不合法"}
	EmptyUserName    = &Base{1005, "用户名不能为空"}
	EmptyPassword    = &Base{1006, "密码不能为空"}
	InvalidToken     = &Base{1007, "token不合法"}
	LiginFailed      = &Base{1008, "登录失败"}

	// file
	FileNotFound     = &Base{2000, "文件不存在"}
	DirNotExist      = &Base{2001, "目录不存在"}
	FileAlreadyExist = &Base{2002, "文件名已存在"}
	InvalidDir       = &Base{2003, "非法的目录名"}
	ReadFileFail     = &Base{2004, "读取文件失败"}
	GetTokenFail     = &Base{2005, "获取token失败"}
)
