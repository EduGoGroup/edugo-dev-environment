-- ========================================
-- SEEDS: Usuarios de Prueba
-- ========================================
-- Ejecutar después de 01_schools.sql
-- Password para todos: "password123" (hash bcrypt)
-- NOTA: En desarrollo usar contraseñas simples, en producción nunca usar seeds

-- Hash bcrypt de "password123" generado con cost 10
-- $2a$10$N9qo8uLOickgx2ZMRZoMy.MqrANdHfQ.8KepbEPQpJ1KAITdP4HHe

INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_active, email_verified) VALUES
-- ========================================
-- ADMINISTRADORES
-- ========================================
('11111111-1111-1111-1111-111111111111',
 'admin@edugo.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrANdHfQ.8KepbEPQpJ1KAITdP4HHe',
 'Admin',
 'Sistema',
 'admin',
 true,
 true),

('11111111-1111-1111-1111-111111111112',
 'superadmin@edugo.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrANdHfQ.8KepbEPQpJ1KAITdP4HHe',
 'Super',
 'Admin',
 'admin',
 true,
 true),

-- ========================================
-- DOCENTES
-- ========================================
('22222222-2222-2222-2222-222222222221',
 'teacher.fisica@edugo.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrANdHfQ.8KepbEPQpJ1KAITdP4HHe',
 'María',
 'González',
 'teacher',
 true,
 true),

('22222222-2222-2222-2222-222222222222',
 'teacher.matematicas@edugo.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrANdHfQ.8KepbEPQpJ1KAITdP4HHe',
 'Carlos',
 'Rodríguez',
 'teacher',
 true,
 true),

('22222222-2222-2222-2222-222222222223',
 'teacher.historia@edugo.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrANdHfQ.8KepbEPQpJ1KAITdP4HHe',
 'Ana',
 'Martínez',
 'teacher',
 true,
 true),

('22222222-2222-2222-2222-222222222224',
 'teacher.ingles@edugo.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrANdHfQ.8KepbEPQpJ1KAITdP4HHe',
 'Roberto',
 'Silva',
 'teacher',
 true,
 true),

('22222222-2222-2222-2222-222222222225',
 'teacher.biologia@edugo.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrANdHfQ.8KepbEPQpJ1KAITdP4HHe',
 'Patricia',
 'Vargas',
 'teacher',
 true,
 true),

-- ========================================
-- ESTUDIANTES
-- ========================================
('33333333-3333-3333-3333-333333333331',
 'student1@edugo.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrANdHfQ.8KepbEPQpJ1KAITdP4HHe',
 'Juan',
 'Pérez',
 'student',
 true,
 true),

('33333333-3333-3333-3333-333333333332',
 'student2@edugo.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrANdHfQ.8KepbEPQpJ1KAITdP4HHe',
 'Camila',
 'López',
 'student',
 true,
 true),

('33333333-3333-3333-3333-333333333333',
 'student3@edugo.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrANdHfQ.8KepbEPQpJ1KAITdP4HHe',
 'Diego',
 'Hernández',
 'student',
 true,
 true),

('33333333-3333-3333-3333-333333333334',
 'student4@edugo.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrANdHfQ.8KepbEPQpJ1KAITdP4HHe',
 'Valentina',
 'García',
 'student',
 true,
 true),

('33333333-3333-3333-3333-333333333335',
 'student5@edugo.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrANdHfQ.8KepbEPQpJ1KAITdP4HHe',
 'Sebastián',
 'Torres',
 'student',
 true,
 true),

('33333333-3333-3333-3333-333333333336',
 'student6@edugo.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrANdHfQ.8KepbEPQpJ1KAITdP4HHe',
 'Isabella',
 'Ramírez',
 'student',
 true,
 true),

('33333333-3333-3333-3333-333333333337',
 'student7@edugo.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrANdHfQ.8KepbEPQpJ1KAITdP4HHe',
 'Matías',
 'Flores',
 'student',
 true,
 true),

('33333333-3333-3333-3333-333333333338',
 'student8@edugo.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrANdHfQ.8KepbEPQpJ1KAITdP4HHe',
 'Sofía',
 'Muñoz',
 'student',
 true,
 true),

('33333333-3333-3333-3333-333333333339',
 'student9@edugo.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrANdHfQ.8KepbEPQpJ1KAITdP4HHe',
 'Benjamín',
 'Castro',
 'student',
 true,
 true),

('33333333-3333-3333-3333-333333333340',
 'student10@edugo.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrANdHfQ.8KepbEPQpJ1KAITdP4HHe',
 'Antonella',
 'Rojas',
 'student',
 true,
 true)

ON CONFLICT (email) DO NOTHING;

-- Verificación
SELECT 'Seeds de usuarios cargados: ' || COUNT(*) || ' usuarios' AS resultado FROM users;
