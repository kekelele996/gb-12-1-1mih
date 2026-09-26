package main

import (
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
)

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.University{},
		&model.ApplicationProject{},
		&model.Document{},
		&model.DocumentVersion{},
		&model.Annotation{},
		&model.MaterialItem{},
		&model.TimelineNode{},
		&model.Recommendation{},
		&model.Message{},
	)
}

func seed(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	logger := slog.Default()

	adminHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	studentHash, _ := bcrypt.GenerateFromPassword([]byte("student123"), bcrypt.DefaultCost)
	counselorHash, _ := bcrypt.GenerateFromPassword([]byte("counselor123"), bcrypt.DefaultCost)
	admin := &model.User{Username: "admin", Email: "admin@gbstudyapply.local", PasswordHash: string(adminHash), RealName: "系统管理员", Role: "admin"}
	student := &model.User{Username: "student1", Email: "student1@gbstudyapply.local", PasswordHash: string(studentHash), RealName: "李同学", Phone: "13800000001", Role: "student"}
	counselor := &model.User{Username: "counselor1", Email: "counselor1@gbstudyapply.local", PasswordHash: string(counselorHash), RealName: "王顾问", Phone: "13800000002", Role: "counselor"}
	if err := db.Create(admin).Error; err != nil {
		return err
	}
	if err := db.Create(student).Error; err != nil {
		return err
	}
	if err := db.Create(counselor).Error; err != nil {
		return err
	}

	unis := []model.University{
		{Name: "哥伦比亚大学", Country: "美国", City: "纽约", Ranking: 3, TopMajors: `["计算机科学","金融工程"]`, ApplicationDeadline: "2027-01-15", TuitionRange: "$65,000-$70,000", Requirements: `{"gpa":3.5,"language":"TOEFL 100","test":"GRE 320"}`},
		{Name: "伦敦大学学院", Country: "英国", City: "伦敦", Ranking: 8, TopMajors: `["教育学","建筑"]`, ApplicationDeadline: "2027-01-31", TuitionRange: "£29,000-£34,000", Requirements: `{"gpa":3.3,"language":"IELTS 7.0","test":"无需GRE"}`},
		{Name: "新加坡国立大学", Country: "新加坡", City: "新加坡", Ranking: 11, TopMajors: `["计算机","电子工程"]`, ApplicationDeadline: "2027-02-15", TuitionRange: "SGD 40,000", Requirements: `{"gpa":3.4,"language":"TOEFL 95","test":"GRE 315"}`},
		{Name: "墨尔本大学", Country: "澳大利亚", City: "墨尔本", Ranking: 24, TopMajors: `["商科","医学"]`, ApplicationDeadline: "2027-03-31", TuitionRange: "AUD 45,000", Requirements: `{"gpa":3.2,"language":"IELTS 6.5","test":"无需GRE"}`},
		{Name: "多伦多大学", Country: "加拿大", City: "多伦多", Ranking: 21, TopMajors: `["工程","生命科学"]`, ApplicationDeadline: "2027-01-10", TuitionRange: "CAD 55,000", Requirements: `{"gpa":3.4,"language":"TOEFL 100","test":"GRE 320"}`},
	}
	if err := db.Create(&unis).Error; err != nil {
		return err
	}

	projects := []model.ApplicationProject{
		{StudentID: student.ID, CounselorID: counselor.ID, UniversityID: unis[0].ID, Major: "计算机科学", Round: "秋季", Status: "planning"},
		{StudentID: student.ID, CounselorID: counselor.ID, UniversityID: unis[1].ID, Major: "教育学", Round: "秋季", Status: "submitted"},
		{StudentID: student.ID, CounselorID: counselor.ID, UniversityID: unis[2].ID, Major: "计算机", Round: "春季", Status: "waiting"},
		{StudentID: student.ID, CounselorID: counselor.ID, UniversityID: unis[3].ID, Major: "商科", Round: "秋季", Status: "admitted"},
	}
	if err := db.Create(&projects).Error; err != nil {
		return err
	}

	docs := []model.Document{
		{ApplicationID: projects[0].ID, DocType: "ps", Title: "个人陈述 v1", Content: "我的学术背景与申请动机……", CurrentVersion: 1},
		{ApplicationID: projects[1].ID, DocType: "cv", Title: "简历", Content: "教育经历：……", CurrentVersion: 1},
	}
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	versions := []model.DocumentVersion{
		{DocumentID: docs[0].ID, Content: "我的学术背景与申请动机……", VersionNo: 1, ChangeSummary: "初始版本", CreatedBy: student.ID},
		{DocumentID: docs[1].ID, Content: "教育经历：……", VersionNo: 1, ChangeSummary: "初始版本", CreatedBy: student.ID},
	}
	if err := db.Create(&versions).Error; err != nil {
		return err
	}
	annotations := []model.Annotation{
		{DocumentID: docs[0].ID, CounselorID: counselor.ID, Content: "建议补充研究经历段落", StartOffset: 0, EndOffset: 20},
	}
	if err := db.Create(&annotations).Error; err != nil {
		return err
	}

	materials := []model.MaterialItem{
		{ApplicationID: projects[0].ID, Name: "成绩单", Category: "学术材料", IsRequired: true, Status: "pending"},
		{ApplicationID: projects[0].ID, Name: "推荐信 x2", Category: "文书", IsRequired: true, Status: "uploaded"},
		{ApplicationID: projects[0].ID, Name: "语言成绩", Category: "标化", IsRequired: true, Status: "approved"},
	}
	if err := db.Create(&materials).Error; err != nil {
		return err
	}

	nodes := []model.TimelineNode{
		{ApplicationID: projects[0].ID, Title: "GRE 考试", NodeType: "language_exam", DueDate: time.Now().AddDate(0, 0, 20), IsDone: false},
		{ApplicationID: projects[0].ID, Title: "文书终稿", NodeType: "essay_deadline", DueDate: time.Now().AddDate(0, 0, 45), IsDone: false},
		{ApplicationID: projects[1].ID, Title: "提交申请", NodeType: "application_deadline", DueDate: time.Now().AddDate(0, 0, 10), IsDone: false},
	}
	if err := db.Create(&nodes).Error; err != nil {
		return err
	}

	recs := []model.Recommendation{
		{StudentID: student.ID, CounselorID: counselor.ID, UniversityIDs: "[1,3,5]", Reason: "冲稳保三档组合：冲刺哥大、稳妥UCL、保底墨尔本"},
	}
	if err := db.Create(&recs).Error; err != nil {
		return err
	}

	msgs := []model.Message{
		{SenderID: counselor.ID, ReceiverID: student.ID, Content: "已收到你的申请材料，下周沟通文书修改意见。", IsRead: false},
		{SenderID: 0, ReceiverID: student.ID, Content: "系统提醒：GRE 考试节点即将截止，请及时准备。", IsRead: false},
	}
	if err := db.Create(&msgs).Error; err != nil {
		return err
	}

	logger.Info("gbstudyapply seed data created",
		"users", 3, "universities", len(unis), "projects", len(projects), "documents", len(docs), "messages", len(msgs))
	return nil
}
