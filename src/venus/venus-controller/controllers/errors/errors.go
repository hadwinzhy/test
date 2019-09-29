package errors

type Error struct {
	ErrorCode
	Detail string `json:"detail"`
}

type ErrorResponse struct {
	Errors []*Error `json:"errors"`
}

type ErrorCode struct {
	HTTPStatus int    `json:"http_status"`
	Code       int    `json:"code"`
	Title      string `json:"title"`
	TitleZH    string `json:"title_zh"`
}

var (
	ErrorRouterNotExist       = ErrorCode{404, 1000, "Router Not Exist", "路由不存在"}
	ErrorPassowrdInvalid      = ErrorCode{400, 1001, "Password invalid", "密码无效"}
	ErrorInvalidParams        = ErrorCode{400, 1003, "Invalid params", "参数无效"}
	ErrorTokenMissing         = ErrorCode{401, 1002, "Token is missing", "缺少token"}
	ErrorTokenInvalid         = ErrorCode{401, 1004, "Invalid Token", "失效token"}
	ErrorCompanyNameDup       = ErrorCode{400, 1005, "Company already exists", "公司已经存在"}
	ErrorNoCompany            = ErrorCode{400, 1006, "Admin does not belong to any company", "管理员不属于任何公司"}
	ErrorDB                   = ErrorCode{400, 1007, "Error from database", "数据库错误"}
	ErrorAdminPhoneDup        = ErrorCode{400, 1008, "Admin phone already exists", "注册手机号已存在"}
	ErrorRecordNotFound       = ErrorCode{404, 2001, "Resource is not found", "没有找到相应记录"}
	ErrorUnexpected           = ErrorCode{400, 4001, "Unexpected error", "遇到预期以外的错误"}
	ErrorFileTooLarge         = ErrorCode{400, 4002, "File is Too Large", "文件过大"}
	ErrorFileSaveFailed       = ErrorCode{400, 4003, "File save failed", "文件保存失败"}
	ErrorNotPermitted         = ErrorCode{403, 2002, "Action is not permitted", "没有权限"}
	ErrorAdminPending         = ErrorCode{403, 2003, "Admin state is now pending", "管理员正在审核中"}
	ErrorAdminRejected        = ErrorCode{403, 2004, "Admin state is now rejected", "管理员申请被拒绝"}
	ErrorReportTimeInterval   = ErrorCode{400, 2005, "The time interval over five day", "时间跨度大于5天"}
	ErrorRangeConflicted      = ErrorCode{400, 2006, "Range Conflicted", "范围冲突"}
	ErrorTooManyItems         = ErrorCode{400, 2007, "Too many items", "项目数过多"}
	ErrorCustomerGroupType    = ErrorCode{400, 2008, "customer group type invalid", "会员组类型不存在"}
	ErrorDateRejected         = ErrorCode{400, 4004, "Activity end time should be after start time", "活动的结束日期不能早于活动的开始日期"}
	ErrorDeviceMacAddress     = ErrorCode{400, 1009, "Device mac address already exists", "设备的mac地址已经存在"}
	ErrorSmFloorName          = ErrorCode{400, 1010, "Floor name should contain number and letter", "楼层名称只能包含数字、字母和汉字"}
	ErrorSmFloorNameMaxLength = ErrorCode{400, 1011, "The length of floor name should be less than 15", "楼层名称最长为15个字符"}
	ErrorSmBusinessTypeLength = ErrorCode{400, 1012, "The length of business type name should be less than 20", "业态类型的长度最长为20"}
	ErrorSmBusinessTypeCount  = ErrorCode{400, 1013, "Business type count should be less than 20", "业态类型的个数最大为20"}
	ErrorSmShopNameExists     = ErrorCode{400, 1014, "Shop name is already exists", "商铺名称已经存在"}
	ErrorRegionNameDup        = ErrorCode{400, 1015, "Region Name already exists", "区域名称已存在"}
	ErrorPictureQuality       = ErrorCode{400, 1016, "Picture Error", "无人脸、图片大小不符合规范或者注册照片质量不符合规范"}
	ErrorTimeout              = ErrorCode{400, 1017, "Message reply timeout", "响应超时"}
)

// ResponseWithErrorCode will abort with error code
func MakeErrorWithErrorCode(errorCode *ErrorCode, message string) Error {
	return Error{ErrorCode: *errorCode, Detail: message}
}

func MakeDBError(message string) Error {
	return Error{ErrorCode: ErrorDB, Detail: message}
}

func MakeNotPermittedError() Error {
	return Error{ErrorCode: ErrorNotPermitted}
}

func MakeInvalidaParamsError(message string) Error {
	return Error{ErrorCode: ErrorInvalidParams, Detail: message}
}

func MakeNotFoundError(message string) Error {
	return Error{ErrorCode: ErrorRecordNotFound, Detail: message}
}

func MakeUnexpectedError(message string) Error {
	return Error{ErrorCode: ErrorUnexpected, Detail: message}
}

func MakePictureError(message string) Error {
	return Error{ErrorCode: ErrorPictureQuality, Detail: message}
}
