-- gbstudyapply 留学申请管理系统 初始化脚本（PostgreSQL）
-- 由 postgres 官方镜像 /docker-entrypoint-initdb.d 首次启动时自动执行。

CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(64) NOT NULL,
  email VARCHAR(128) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  real_name VARCHAR(64),
  phone VARCHAR(32),
  role VARCHAR(16) NOT NULL DEFAULT 'student',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT uni_users_username UNIQUE (username),
  CONSTRAINT uni_users_email UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS universities (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  country VARCHAR(64),
  city VARCHAR(64),
  ranking INT DEFAULT 0,
  top_majors JSONB DEFAULT '[]',
  application_deadline VARCHAR(32),
  tuition_range VARCHAR(64),
  requirements JSONB DEFAULT '{}',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT uni_universities_name UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS application_projects (
  id BIGSERIAL PRIMARY KEY,
  student_id BIGINT NOT NULL,
  counselor_id BIGINT DEFAULT 0,
  university_id BIGINT NOT NULL,
  major VARCHAR(128) NOT NULL,
  round VARCHAR(32),
  status VARCHAR(16) NOT NULL DEFAULT 'planning',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_app_student FOREIGN KEY (student_id) REFERENCES users(id),
  CONSTRAINT fk_app_univ FOREIGN KEY (university_id) REFERENCES universities(id)
);

CREATE TABLE IF NOT EXISTS documents (
  id BIGSERIAL PRIMARY KEY,
  application_id BIGINT NOT NULL,
  doc_type VARCHAR(16) NOT NULL,
  title VARCHAR(255) NOT NULL,
  content TEXT,
  current_version INT DEFAULT 1,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_doc_app FOREIGN KEY (application_id) REFERENCES application_projects(id)
);

CREATE TABLE IF NOT EXISTS document_versions (
  id BIGSERIAL PRIMARY KEY,
  document_id BIGINT NOT NULL,
  content TEXT,
  version_no INT DEFAULT 1,
  change_summary VARCHAR(255),
  created_by BIGINT DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_version_doc FOREIGN KEY (document_id) REFERENCES documents(id)
);

CREATE TABLE IF NOT EXISTS annotations (
  id BIGSERIAL PRIMARY KEY,
  document_id BIGINT NOT NULL,
  counselor_id BIGINT NOT NULL,
  content VARCHAR(1000) NOT NULL,
  start_offset INT DEFAULT 0,
  end_offset INT DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_ann_doc FOREIGN KEY (document_id) REFERENCES documents(id)
);

CREATE TABLE IF NOT EXISTS material_items (
  id BIGSERIAL PRIMARY KEY,
  application_id BIGINT NOT NULL,
  name VARCHAR(128) NOT NULL,
  category VARCHAR(64),
  is_required BOOLEAN DEFAULT TRUE,
  status VARCHAR(16) NOT NULL DEFAULT 'pending',
  file_url VARCHAR(255),
  uploaded_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_mat_app FOREIGN KEY (application_id) REFERENCES application_projects(id)
);

CREATE TABLE IF NOT EXISTS timeline_nodes (
  id BIGSERIAL PRIMARY KEY,
  application_id BIGINT NOT NULL,
  title VARCHAR(255) NOT NULL,
  node_type VARCHAR(32),
  due_date DATE,
  is_done BOOLEAN DEFAULT FALSE,
  reminder_sent BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_node_app FOREIGN KEY (application_id) REFERENCES application_projects(id)
);

CREATE TABLE IF NOT EXISTS recommendations (
  id BIGSERIAL PRIMARY KEY,
  student_id BIGINT NOT NULL,
  counselor_id BIGINT NOT NULL,
  university_ids JSONB DEFAULT '[]',
  reason TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_rec_student FOREIGN KEY (student_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS messages (
  id BIGSERIAL PRIMARY KEY,
  sender_id BIGINT DEFAULT 0,
  receiver_id BIGINT NOT NULL,
  content VARCHAR(2000) NOT NULL,
  is_read BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 种子数据
INSERT INTO users (username, email, password_hash, real_name, phone, role) VALUES
  ('admin', 'admin@gbstudyapply.local', '$2a$10$VGETME6mK/u27yF1UwKHkuh0b36LjEpJjw2c4J2L7wPph1pcG0cVO', '系统管理员', NULL, 'admin'),
  ('student1', 'student1@gbstudyapply.local', '$2a$10$aumfvegF9M4D1Im6VWY.POc0EMb3hAy3GI/L5vi2VMUtlU6My6oUO', '李同学', '13800000001', 'student'),
  ('counselor1', 'counselor1@gbstudyapply.local', '$2a$10$DEQ/fzKZ3l.vu3LQ8AG9hODA.Q6BTelg5xgiHqPoVYzKPtj6cMLme', '王顾问', '13800000002', 'counselor');

INSERT INTO universities (name, country, city, ranking, top_majors, application_deadline, tuition_range, requirements) VALUES
  ('哥伦比亚大学', '美国', '纽约', 3, '["计算机科学","金融工程"]', '2027-01-15', '$65,000-$70,000', '{"gpa":3.5,"language":"TOEFL 100","test":"GRE 320"}'),
  ('伦敦大学学院', '英国', '伦敦', 8, '["教育学","建筑"]', '2027-01-31', '£29,000-£34,000', '{"gpa":3.3,"language":"IELTS 7.0","test":"无需GRE"}'),
  ('新加坡国立大学', '新加坡', '新加坡', 11, '["计算机","电子工程"]', '2027-02-15', 'SGD 40,000', '{"gpa":3.4,"language":"TOEFL 95","test":"GRE 315"}'),
  ('墨尔本大学', '澳大利亚', '墨尔本', 24, '["商科","医学"]', '2027-03-31', 'AUD 45,000', '{"gpa":3.2,"language":"IELTS 6.5","test":"无需GRE"}'),
  ('多伦多大学', '加拿大', '多伦多', 21, '["工程","生命科学"]', '2027-01-10', 'CAD 55,000', '{"gpa":3.4,"language":"TOEFL 100","test":"GRE 320"}');

INSERT INTO application_projects (student_id, counselor_id, university_id, major, round, status) VALUES
  (2, 3, 1, '计算机科学', '秋季', 'planning'),
  (2, 3, 2, '教育学', '秋季', 'submitted'),
  (2, 3, 3, '计算机', '春季', 'waiting'),
  (2, 3, 4, '商科', '秋季', 'admitted');

INSERT INTO documents (application_id, doc_type, title, content, current_version) VALUES
  (1, 'ps', '个人陈述 v1', '我的学术背景与申请动机……', 1),
  (2, 'cv', '简历', '教育经历：……', 1);

INSERT INTO document_versions (document_id, content, version_no, change_summary, created_by) VALUES
  (1, '我的学术背景与申请动机……', 1, '初始版本', 2),
  (2, '教育经历：……', 1, '初始版本', 2);

INSERT INTO annotations (document_id, counselor_id, content, start_offset, end_offset) VALUES
  (1, 3, '建议补充研究经历段落', 0, 20);

INSERT INTO material_items (application_id, name, category, is_required, status) VALUES
  (1, '成绩单', '学术材料', TRUE, 'pending'),
  (1, '推荐信 x2', '文书', TRUE, 'uploaded'),
  (1, '语言成绩', '标化', TRUE, 'approved');

INSERT INTO timeline_nodes (application_id, title, node_type, due_date, is_done) VALUES
  (1, 'GRE 考试', 'language_exam', CURRENT_DATE + 20, FALSE),
  (1, '文书终稿', 'essay_deadline', CURRENT_DATE + 45, FALSE),
  (2, '提交申请', 'application_deadline', CURRENT_DATE + 10, FALSE);

INSERT INTO recommendations (student_id, counselor_id, university_ids, reason) VALUES
  (2, 3, '[1,2,4]', '冲稳保三档组合：冲刺哥大、稳妥UCL、保底墨尔本');

INSERT INTO messages (sender_id, receiver_id, content, is_read) VALUES
  (3, 2, '已收到你的申请材料，下周沟通文书修改意见。', FALSE),
  (0, 2, '系统提醒：GRE 考试节点即将截止，请及时准备。', FALSE);
