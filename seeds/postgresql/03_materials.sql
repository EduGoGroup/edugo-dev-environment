-- ========================================
-- SEEDS: Materiales de Prueba
-- ========================================
-- Ejecutar después de 02_users.sql
-- Materiales educativos de ejemplo para diferentes asignaturas

INSERT INTO materials (id, school_id, uploaded_by_teacher_id, title, description, subject, grade, file_url, file_type, file_size_bytes, status) VALUES
-- ========================================
-- MATERIALES DE FÍSICA (Docente: María González)
-- ========================================
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
 '44444444-4444-4444-4444-444444444444',
 '22222222-2222-2222-2222-222222222221',
 'Introducción a la Física Cuántica',
 'Material educativo sobre conceptos básicos de física cuántica: dualidad onda-partícula, principio de incertidumbre',
 'Física',
 '12th',
 's3://edugo-materials-dev/fisica-cuantica.pdf',
 'application/pdf',
 2048000,
 'ready'),

('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaab',
 '44444444-4444-4444-4444-444444444444',
 '22222222-2222-2222-2222-222222222221',
 'Mecánica Newtoniana - Leyes del Movimiento',
 'Las tres leyes de Newton explicadas con ejemplos prácticos y ejercicios',
 'Física',
 '10th',
 's3://edugo-materials-dev/mecanica-newton.pdf',
 'application/pdf',
 1536000,
 'ready'),

('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaac',
 '44444444-4444-4444-4444-444444444444',
 '22222222-2222-2222-2222-222222222221',
 'Termodinámica - Principios Fundamentales',
 'Estudio de calor, temperatura y las leyes de la termodinámica',
 'Física',
 '11th',
 's3://edugo-materials-dev/termodinamica.pdf',
 'application/pdf',
 1824000,
 'processing'),

-- ========================================
-- MATERIALES DE MATEMÁTICAS (Docente: Carlos Rodríguez)
-- ========================================
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbba',
 '44444444-4444-4444-4444-444444444444',
 '22222222-2222-2222-2222-222222222222',
 'Álgebra Lineal - Matrices y Determinantes',
 'Ejercicios y teoría sobre matrices, determinantes y sistemas de ecuaciones',
 'Matemáticas',
 '11th',
 's3://edugo-materials-dev/algebra-matrices.pdf',
 'application/pdf',
 1524000,
 'ready'),

('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
 '44444444-4444-4444-4444-444444444444',
 '22222222-2222-2222-2222-222222222222',
 'Cálculo Diferencial - Límites y Derivadas',
 'Conceptos de límites, continuidad y derivadas con aplicaciones',
 'Matemáticas',
 '12th',
 's3://edugo-materials-dev/calculo-derivadas.pdf',
 'application/pdf',
 2256000,
 'ready'),

('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbc',
 '55555555-5555-5555-5555-555555555555',
 '22222222-2222-2222-2222-222222222222',
 'Geometría Euclidiana - Teoremas Fundamentales',
 'Teoremas de Tales, Pitágoras y aplicaciones en problemas geométricos',
 'Matemáticas',
 '9th',
 's3://edugo-materials-dev/geometria.pdf',
 'application/pdf',
 1128000,
 'ready'),

-- ========================================
-- MATERIALES DE HISTORIA (Docente: Ana Martínez)
-- ========================================
('cccccccc-cccc-cccc-cccc-ccccccccccca',
 '55555555-5555-5555-5555-555555555555',
 '22222222-2222-2222-2222-222222222223',
 'Historia de Chile - Siglo XX',
 'Principales eventos históricos del siglo XX en Chile',
 'Historia',
 '9th',
 's3://edugo-materials-dev/historia-chile-xx.pdf',
 'application/pdf',
 3072000,
 'ready'),

('cccccccc-cccc-cccc-cccc-cccccccccccb',
 '55555555-5555-5555-5555-555555555555',
 '22222222-2222-2222-2222-222222222223',
 'Segunda Guerra Mundial - Causas y Consecuencias',
 'Análisis detallado de la Segunda Guerra Mundial',
 'Historia',
 '10th',
 's3://edugo-materials-dev/segunda-guerra.pdf',
 'application/pdf',
 4096000,
 'ready'),

-- ========================================
-- MATERIALES DE INGLÉS (Docente: Roberto Silva)
-- ========================================
('dddddddd-dddd-dddd-dddd-ddddddddddda',
 '44444444-4444-4444-4444-444444444444',
 '22222222-2222-2222-2222-222222222224',
 'English Grammar - Tenses Review',
 'Complete review of all English tenses with exercises',
 'Inglés',
 '8th',
 's3://edugo-materials-dev/english-tenses.pdf',
 'application/pdf',
 1256000,
 'ready'),

('dddddddd-dddd-dddd-dddd-dddddddddddd',
 '44444444-4444-4444-4444-444444444444',
 '22222222-2222-2222-2222-222222222224',
 'Reading Comprehension - Advanced Level',
 'Advanced reading exercises with vocabulary building',
 'Inglés',
 '11th',
 's3://edugo-materials-dev/reading-advanced.pdf',
 'application/pdf',
 1872000,
 'processing'),

-- ========================================
-- MATERIALES DE BIOLOGÍA (Docente: Patricia Vargas)
-- ========================================
('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeea',
 '66666666-6666-6666-6666-666666666666',
 '22222222-2222-2222-2222-222222222225',
 'Biología Celular - Estructura y Función',
 'Estudio de la célula, sus organelos y funciones vitales',
 'Biología',
 '9th',
 's3://edugo-materials-dev/biologia-celular.pdf',
 'application/pdf',
 2560000,
 'ready'),

('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee',
 '66666666-6666-6666-6666-666666666666',
 '22222222-2222-2222-2222-222222222225',
 'Genética y Herencia - Leyes de Mendel',
 'Principios de genética, herencia mendeliana y problemas',
 'Biología',
 '10th',
 's3://edugo-materials-dev/genetica-mendel.pdf',
 'application/pdf',
 1984000,
 'ready')

ON CONFLICT (id) DO NOTHING;

-- Verificación
SELECT 'Seeds de materiales cargados: ' || COUNT(*) || ' materiales' AS resultado FROM materials;
