-- Esquema inicial del sistema de biblioteca

-- Fuerza charset UTF-8 en la importación inicial (evita caracteres corruptos
-- como "Cien aÃ±os de soledad" cuando el cliente MySQL importa con latin1)
SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(150) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  role ENUM('student', 'admin') NOT NULL DEFAULT 'student',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS books (
  id INT AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(200) NOT NULL,
  author VARCHAR(150) NOT NULL,
  isbn VARCHAR(20),
  category VARCHAR(80),
  cover_url VARCHAR(255),
  total_copies INT NOT NULL DEFAULT 1,
  available_copies INT NOT NULL DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS loans (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  book_id INT NOT NULL,
  loaned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  due_at TIMESTAMP NOT NULL,
  returned_at TIMESTAMP NULL,
  status ENUM('active', 'returned') NOT NULL DEFAULT 'active',
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (book_id) REFERENCES books(id)
);

-- Datos de ejemplo para que la lista de libros no arranque vacía
INSERT INTO books (title, author, isbn, category, cover_url, total_copies, available_copies) VALUES
('Cien años de soledad', 'Gabriel García Márquez', '9780307474728', 'Novela', '', 3, 3),
('Clean Code', 'Robert C. Martin', '9780132350884', 'Tecnología', '', 2, 2),
('El principito', 'Antoine de Saint-Exupéry', '9780156012195', 'Infantil', '', 4, 4),
('Sapiens', 'Yuval Noah Harari', '9780062316097', 'Historia', '', 2, 2),
('The Pragmatic Programmer', 'Andrew Hunt', '9780135957059', 'Tecnología', '', 2, 2);
('The Go Programming Language (Addison-Wesley Professional Computing Series) ', 'Alan Donovan', '', 'Tecnología', '', 2, 2);
