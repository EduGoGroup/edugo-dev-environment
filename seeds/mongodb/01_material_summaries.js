// ========================================
// SEEDS: Material Summaries
// ========================================
// Resúmenes generados por IA para materiales de prueba

print("Insertando material_summaries...");

db.material_summaries.insertMany([
  // Física Cuántica
  {
    material_id: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    summary: "Este material introduce los conceptos fundamentales de la física cuántica, incluyendo la dualidad onda-partícula y el principio de incertidumbre de Heisenberg. Se explora cómo las partículas subatómicas exhiben comportamientos que desafían nuestra intuición clásica.",
    key_points: [
      "La luz y la materia exhiben propiedades tanto de onda como de partícula",
      "El principio de incertidumbre establece límites a lo que podemos medir simultáneamente",
      "Los estados cuánticos se describen mediante funciones de onda",
      "El experimento de la doble rendija demuestra la naturaleza dual de la luz",
      "La superposición cuántica permite que las partículas existan en múltiples estados"
    ],
    topics: ["Física Cuántica", "Dualidad onda-partícula", "Principio de Incertidumbre", "Mecánica Cuántica"],
    difficulty_level: "advanced",
    estimated_reading_time_minutes: 45,
    created_at: new Date(),
    updated_at: new Date()
  },

  // Mecánica Newtoniana
  {
    material_id: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaab",
    summary: "Material completo sobre las tres leyes del movimiento de Newton. Incluye explicaciones detalladas con ejemplos cotidianos y ejercicios prácticos para aplicar cada ley.",
    key_points: [
      "Primera Ley: Un objeto permanece en reposo o movimiento uniforme a menos que actúe una fuerza externa",
      "Segunda Ley: F = ma, la fuerza es igual a masa por aceleración",
      "Tercera Ley: Para cada acción hay una reacción igual y opuesta",
      "Las leyes de Newton son la base de la mecánica clásica",
      "Aplicaciones en ingeniería, deportes y vida cotidiana"
    ],
    topics: ["Mecánica Clásica", "Leyes de Newton", "Fuerza", "Movimiento"],
    difficulty_level: "intermediate",
    estimated_reading_time_minutes: 30,
    created_at: new Date(),
    updated_at: new Date()
  },

  // Álgebra Lineal
  {
    material_id: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbba",
    summary: "Guía completa de álgebra lineal enfocada en matrices y determinantes. Cubre operaciones básicas, propiedades y aplicaciones en sistemas de ecuaciones lineales.",
    key_points: [
      "Una matriz es un arreglo rectangular de números",
      "Las operaciones básicas incluyen suma, resta y multiplicación",
      "El determinante indica si una matriz es invertible",
      "La matriz inversa existe solo cuando el determinante es distinto de cero",
      "Los sistemas de ecuaciones se resuelven usando matrices"
    ],
    topics: ["Álgebra Lineal", "Matrices", "Determinantes", "Sistemas de Ecuaciones"],
    difficulty_level: "intermediate",
    estimated_reading_time_minutes: 40,
    created_at: new Date(),
    updated_at: new Date()
  },

  // Cálculo Diferencial
  {
    material_id: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb",
    summary: "Introducción al cálculo diferencial cubriendo límites, continuidad y derivadas. Material incluye reglas de derivación y aplicaciones en problemas de optimización.",
    key_points: [
      "El límite describe el comportamiento de una función cuando x se acerca a un valor",
      "Una función es continua si no tiene saltos ni discontinuidades",
      "La derivada representa la tasa de cambio instantánea",
      "Reglas de derivación: potencia, producto, cociente y cadena",
      "Las derivadas se usan para encontrar máximos y mínimos"
    ],
    topics: ["Cálculo", "Límites", "Derivadas", "Optimización"],
    difficulty_level: "advanced",
    estimated_reading_time_minutes: 50,
    created_at: new Date(),
    updated_at: new Date()
  },

  // Historia de Chile
  {
    material_id: "cccccccc-cccc-cccc-cccc-ccccccccccca",
    summary: "Recorrido por los eventos más significativos del siglo XX en Chile, desde la cuestión social hasta el retorno a la democracia, pasando por los gobiernos radicales y el período de dictadura.",
    key_points: [
      "La cuestión social marcó las primeras décadas del siglo",
      "Los gobiernos radicales implementaron importantes reformas",
      "El gobierno de la Unidad Popular representó un cambio significativo",
      "El golpe de 1973 inició un período de dictadura militar",
      "En 1990 Chile retornó a la democracia"
    ],
    topics: ["Historia de Chile", "Siglo XX", "Política Chilena", "Democracia"],
    difficulty_level: "intermediate",
    estimated_reading_time_minutes: 35,
    created_at: new Date(),
    updated_at: new Date()
  },

  // Biología Celular
  {
    material_id: "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeea",
    summary: "Estudio detallado de la célula como unidad básica de la vida. Cubre la estructura celular, organelos y sus funciones, diferenciando entre células procariotas y eucariotas.",
    key_points: [
      "La célula es la unidad básica de todos los seres vivos",
      "Las células eucariotas tienen núcleo definido",
      "Los organelos cumplen funciones específicas",
      "La mitocondria es la central energética de la célula",
      "El ADN contiene la información genética"
    ],
    topics: ["Biología", "Célula", "Organelos", "Eucariotas", "Procariotas"],
    difficulty_level: "basic",
    estimated_reading_time_minutes: 25,
    created_at: new Date(),
    updated_at: new Date()
  }
]);

print("✅ material_summaries insertados exitosamente");
