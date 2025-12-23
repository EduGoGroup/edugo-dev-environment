-- ========================================
-- SEEDS: Escuelas de Prueba
-- ========================================
-- Ejecutar después de que el migrator haya creado la estructura
-- Estos datos son para desarrollo local únicamente

INSERT INTO schools (id, name, code, city, country, subscription_tier, max_teachers, max_students) VALUES
-- Escuela Premium - Santiago
('44444444-4444-4444-4444-444444444444',
 'Liceo Técnico Santiago',
 'LTS-001',
 'Santiago',
 'Chile',
 'premium',
 50,
 500),

-- Escuela Basic - Valparaíso
('55555555-5555-5555-5555-555555555555',
 'Colegio Valparaíso',
 'CV-002',
 'Valparaíso',
 'Chile',
 'basic',
 20,
 200),

-- Escuela Standard - Concepción
('66666666-6666-6666-6666-666666666666',
 'Instituto Biobío',
 'IB-003',
 'Concepción',
 'Chile',
 'standard',
 30,
 300),

-- Escuela Premium - Buenos Aires
('77777777-7777-7777-7777-777777777777',
 'Colegio Nacional Buenos Aires',
 'CNBA-004',
 'Buenos Aires',
 'Argentina',
 'premium',
 60,
 600),

-- Escuela Basic - Lima
('88888888-8888-8888-8888-888888888888',
 'Institución Educativa Lima',
 'IEL-005',
 'Lima',
 'Perú',
 'basic',
 25,
 250)

ON CONFLICT (code) DO NOTHING;

-- Verificación
SELECT 'Seeds de escuelas cargados: ' || COUNT(*) || ' escuelas' AS resultado FROM schools;
