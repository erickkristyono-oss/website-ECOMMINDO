-- PT Ecommindo Jaya Persada — schema for the website backend
CREATE DATABASE IF NOT EXISTS ecommindo
  CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE ecommindo;

CREATE TABLE IF NOT EXISTS users (
  id            INT AUTO_INCREMENT PRIMARY KEY,
  full_name     VARCHAR(150)  NOT NULL,
  phone         VARCHAR(30)   NOT NULL,
  email         VARCHAR(150)  NOT NULL UNIQUE,
  password_hash VARCHAR(255)  NOT NULL,
  created_at    TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS services (
  id          VARCHAR(40)     PRIMARY KEY,
  name        VARCHAR(150)    NOT NULL,
  short_desc  VARCHAR(255)    NOT NULL,
  description TEXT            NOT NULL,
  price       DECIMAL(14,2)   NOT NULL,
  icon        VARCHAR(20)     NOT NULL DEFAULT ''
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS cart_items (
  id         INT AUTO_INCREMENT PRIMARY KEY,
  user_id    INT NOT NULL,
  service_id VARCHAR(40) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uq_user_service (user_id, service_id),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS orders (
  id         INT AUTO_INCREMENT PRIMARY KEY,
  order_code VARCHAR(30) NOT NULL UNIQUE,
  user_id    INT NOT NULL,
  total      DECIMAL(14,2) NOT NULL,
  status     VARCHAR(40) NOT NULL DEFAULT 'Menunggu Konfirmasi Pembayaran',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS order_items (
  id            INT AUTO_INCREMENT PRIMARY KEY,
  order_id      INT NOT NULL,
  service_id    VARCHAR(40) NOT NULL,
  service_name  VARCHAR(150) NOT NULL,
  price         DECIMAL(14,2) NOT NULL,
  FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
) ENGINE=InnoDB;

INSERT INTO services (id, name, short_desc, description, price, icon) VALUES
('app-dev', 'Jasa Pembuatan Aplikasi', 'Aplikasi web & mobile custom sesuai kebutuhan bisnis Anda.', 'Kami merancang dan membangun aplikasi web maupun mobile dari nol, mulai dari perencanaan arsitektur, desain UI/UX, hingga pengembangan dan deployment, disesuaikan dengan proses bisnis Anda.', 12000000.00, '📱'),
('server-build', 'Jasa Membangun Server', 'Setup dan konfigurasi server yang stabil, cepat, dan aman.', 'Instalasi, konfigurasi, dan optimasi server (on-premise maupun cloud) agar aplikasi dan layanan bisnis Anda berjalan stabil, cepat, dan siap menangani beban tinggi.', 15000000.00, '🖥️'),
('cyber-security', 'Keamanan Siber', 'Audit, proteksi, dan mitigasi risiko keamanan sistem Anda.', 'Layanan audit keamanan, penetration testing, hardening sistem, dan implementasi kebijakan keamanan untuk melindungi data serta infrastruktur digital perusahaan Anda.', 18000000.00, '🔒'),
('infra-management', 'Manajemen Infrastruktur', 'Kelola infrastruktur IT Anda secara menyeluruh dan efisien.', 'Pengelolaan infrastruktur IT secara end-to-end: monitoring, maintenance, scaling, dan dokumentasi, sehingga tim Anda dapat fokus pada bisnis inti tanpa khawatir soal operasional teknis.', 20000000.00, '🗂️'),
('software-dev', 'Pengembangan Perangkat Lunak', 'Pengembangan fitur, integrasi, dan pemeliharaan perangkat lunak.', 'Pengembangan berkelanjutan untuk perangkat lunak yang sudah berjalan: penambahan fitur baru, integrasi sistem pihak ketiga, refactoring, dan pemeliharaan rutin.', 8000000.00, '💻')
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  short_desc = VALUES(short_desc),
  description = VALUES(description),
  price = VALUES(price),
  icon = VALUES(icon);
