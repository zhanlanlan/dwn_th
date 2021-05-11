package proto

var (
	// basic
	StdSuccess  = &Base{Code: 200, Msg: "succes"}
	BadRquest   = &Base{Code: 400, Msg: "参数错误"}
	InternalErr = &Base{Code: 500, Msg: "服务器内部错误"}

	// user
	BadUserName      = &Base{Code: 1000, Msg: "非法的用户名"}
	WrongPassword    = &Base{Code: 1001, Msg: "密码错误"}
	UserNotFound     = &Base{Code: 1002, Msg: "该用户不存在"}
	UserAlreadyExist = &Base{Code: 1003, Msg: "用户已存在"}
	BadPassword      = &Base{Code: 1004, Msg: "密码不合法"}
	EmptyUserName    = &Base{Code: 1005, Msg: "用户名不能为空"}
	EmptyPassword    = &Base{Code: 1006, Msg: "密码不能为空"}
	InvalidToken     = &Base{Code: 1007, Msg: "token不合法"}
	LiginFailed      = &Base{Code: 1008, Msg: "登录失败"}

	// file
	FileNotFound     = &Base{Code: 2000, Msg: "文件不存在"}
	DirNotExist      = &Base{Code: 2001, Msg: "目录不存在"}
	FileAlreadyExist = &Base{Code: 2002, Msg: "文件名已存在"}
	InvalidDir       = &Base{Code: 2003, Msg: "非法的目录名"}
	ReadFileFail     = &Base{Code: 2004, Msg: "读取文件失败"}
)
