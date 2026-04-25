CREATE TABLE IF NOT EXISTS users(
  id bigint auto_increment PRIMARY KEY,
  first_name varchar(255) NOT NULL,
  last_name varchar(255),
  nickname varchar(80) NOT NULL,
  email varchar(255) NOT NULL,
  password_digest mediumtext,
  gender char(1),
  mobile varchar(20),

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,


  UNIQUE INDEX idx_mobile(mobile),
  UNIQUE INDEX idx_nickname (nickname),
  UNIQUE INDEX idx_email (email)
);