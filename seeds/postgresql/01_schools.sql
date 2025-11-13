-- Seeds de Escuelas para Testing
-- Crea 5 escuelas de ejemplo

INSERT INTO schools (id, name, domain, created_at, updated_at) VALUES
('550e8400-e29b-41d4-a716-446655440001', 'Instituto Tecnológico Superior', 'its.edu.mx', NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440002', 'Universidad Nacional', 'un.edu.ar', NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440003', 'Colegio Bilingüe Internacional', 'cbi.edu.co', NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440004', 'Academia de Ciencias Aplicadas', 'aca.edu.pe', NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440005', 'Centro Educativo Digital', 'ced.edu.cl', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

SELECT 'Seeds de escuelas cargados: 5 escuelas' AS resultado;
