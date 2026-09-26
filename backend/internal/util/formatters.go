package util

import (
	"time"
)

// Shared formatters: date, tuition range, status/role/doc-type text.

// FormatDate renders YYYY-MM-DD.
func FormatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

// FormatDateTime renders YYYY-MM-DD HH:mm.
func FormatDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04")
}

// AppStatusText maps application status to Chinese text.
func AppStatusText(s string) string {
	switch s {
	case "planning":
		return "规划中"
	case "preparing":
		return "材料准备中"
	case "submitted":
		return "已提交"
	case "waiting":
		return "等待结果"
	case "admitted":
		return "已录取"
	case "rejected":
		return "已拒绝"
	case "waitlisted":
		return "候补名单"
	default:
		return "未知"
	}
}

// RoleText maps a role to Chinese text.
func RoleText(r string) string {
	switch r {
	case "student":
		return "学生"
	case "counselor":
		return "顾问"
	case "admin":
		return "管理员"
	default:
		return "未知"
	}
}

// DocTypeText maps a document type to Chinese text.
func DocTypeText(t string) string {
	switch t {
	case "ps":
		return "个人陈述"
	case "rl":
		return "推荐信"
	case "cv":
		return "简历"
	case "essay":
		return "Essay"
	default:
		return "其他"
	}
}

// MaterialStatusText maps a material status to Chinese text.
func MaterialStatusText(s string) string {
	switch s {
	case "pending":
		return "待上传"
	case "uploaded":
		return "已上传"
	case "approved":
		return "已审核"
	default:
		return "未知"
	}
}
