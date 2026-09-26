package constants

const (
	MsgRegisterOK         = "注册成功"
	MsgLoginOK            = "登录成功"
	MsgProfileUpdated     = "资料已更新"
	MsgAppCreated         = "申请项目已创建"
	MsgAppStatusUpdated   = "申请状态已更新"
	MsgDocumentSaved      = "文书已保存"
	MsgVersionCreated     = "版本已保存"
	MsgAnnotationAdded    = "批注已添加"
	MsgMaterialUploaded   = "材料已上传"
	MsgTimelineDone       = "节点已标记完成"
	MsgRecommendationMade = "选校方案已推荐"
	MsgMessageSent        = "消息已发送"
	MsgInvalidCredentials = "用户名或密码错误"
	MsgUsernameTaken      = "用户名已存在"
)

// 材料清单 / 提交闸门 业务提示（面向前端展示，使用中文）。
const (
	MsgSubmitMissingMaterials  = "存在必交材料未通过审核，无法提交申请"
	MsgAppNotEditable          = "申请已提交，材料和文书已锁定，不能修改"
	MsgAppReturnedEditable     = "申请已被退回准备中，可继续修改材料和文书"
	MsgMaterialRequiresUpload  = "请先上传材料文件再提交审核"
	MsgMaterialApprovedLocked  = "材料已审核通过，学生不能再修改；如需更换请联系负责顾问"
	MsgMaterialReviewStaffOnly = "只有负责顾问或管理员可以审核材料"
	MsgMaterialManageStaffOnly = "只有负责顾问或管理员可以维护材料清单"
	MsgAppSubmitOwnerOnly      = "只有学生本人可以提交申请"
	MsgAppReturnStaffOnly      = "只有负责顾问或管理员可以退回申请"
	MsgNotResponsibleCounselor = "只有该项目的负责顾问可以执行此操作"
	MsgDocNotEditable          = "申请已提交，文书已锁定，不能编辑或回滚"
	MsgDocStudentOwnerOnly     = "只有学生本人可以编辑文书"
)
