// ========================================
// SEEDS: Assessments (Evaluaciones generadas por IA)
// ========================================
// Evaluaciones de ejemplo asociadas a materiales

print("Insertando assessments...");

db.assessments.insertMany([
  // Assessment de Física Cuántica
  {
    material_id: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    title: "Evaluación: Física Cuántica - Conceptos Básicos",
    description: "Evaluación sobre dualidad onda-partícula y principio de incertidumbre",
    questions: [
      {
        id: "q1",
        type: "multiple_choice",
        question: "¿Qué demuestra el experimento de la doble rendija?",
        options: [
          "Que la luz es solo una onda",
          "Que la luz es solo una partícula",
          "La dualidad onda-partícula de la luz",
          "Que la luz no existe"
        ],
        correct_answer: 2,
        explanation: "El experimento muestra que la luz se comporta como onda al pasar por las rendijas, pero como partícula al ser detectada."
      },
      {
        id: "q2",
        type: "multiple_choice",
        question: "El principio de incertidumbre de Heisenberg establece que:",
        options: [
          "Todo es incierto en física cuántica",
          "No podemos medir posición y momento simultáneamente con precisión arbitraria",
          "Las partículas no existen",
          "La velocidad de la luz es variable"
        ],
        correct_answer: 1,
        explanation: "Heisenberg demostró que existe un límite fundamental a la precisión con que podemos conocer ciertas parejas de propiedades."
      },
      {
        id: "q3",
        type: "true_false",
        question: "La superposición cuántica permite que una partícula esté en múltiples estados simultáneamente.",
        correct_answer: true,
        explanation: "Hasta que se realiza una medición, las partículas cuánticas existen en una superposición de estados."
      }
    ],
    total_points: 30,
    passing_score: 18,
    time_limit_minutes: 15,
    created_at: new Date(),
    updated_at: new Date()
  },

  // Assessment de Mecánica Newtoniana
  {
    material_id: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaab",
    title: "Evaluación: Leyes de Newton",
    description: "Evaluación sobre las tres leyes del movimiento",
    questions: [
      {
        id: "q1",
        type: "multiple_choice",
        question: "La Primera Ley de Newton también se conoce como:",
        options: [
          "Ley de acción y reacción",
          "Ley de la inercia",
          "Ley de la gravitación",
          "Ley de conservación"
        ],
        correct_answer: 1,
        explanation: "La primera ley describe la inercia: la tendencia de los objetos a mantener su estado de movimiento."
      },
      {
        id: "q2",
        type: "multiple_choice",
        question: "Si F = ma, ¿qué sucede con la aceleración si duplicamos la fuerza?",
        options: [
          "Se reduce a la mitad",
          "Se mantiene igual",
          "Se duplica",
          "Se cuadruplica"
        ],
        correct_answer: 2,
        explanation: "La aceleración es directamente proporcional a la fuerza aplicada."
      },
      {
        id: "q3",
        type: "short_answer",
        question: "Explica con un ejemplo cotidiano la tercera ley de Newton.",
        sample_answer: "Cuando caminas, empujas el suelo hacia atrás y el suelo te empuja hacia adelante, permitiéndote avanzar.",
        keywords: ["acción", "reacción", "igual", "opuesta"]
      }
    ],
    total_points: 30,
    passing_score: 18,
    time_limit_minutes: 20,
    created_at: new Date(),
    updated_at: new Date()
  },

  // Assessment de Álgebra Lineal
  {
    material_id: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbba",
    title: "Evaluación: Matrices y Determinantes",
    description: "Evaluación sobre operaciones con matrices y cálculo de determinantes",
    questions: [
      {
        id: "q1",
        type: "multiple_choice",
        question: "¿Cuándo una matriz tiene inversa?",
        options: [
          "Siempre",
          "Cuando es cuadrada",
          "Cuando su determinante es distinto de cero",
          "Cuando tiene más filas que columnas"
        ],
        correct_answer: 2,
        explanation: "Una matriz es invertible si y solo si su determinante es diferente de cero."
      },
      {
        id: "q2",
        type: "true_false",
        question: "La multiplicación de matrices es conmutativa (A×B = B×A).",
        correct_answer: false,
        explanation: "La multiplicación de matrices NO es conmutativa. El orden importa."
      },
      {
        id: "q3",
        type: "multiple_choice",
        question: "El determinante de una matriz identidad 3x3 es:",
        options: [
          "0",
          "1",
          "3",
          "Depende de los valores"
        ],
        correct_answer: 1,
        explanation: "El determinante de cualquier matriz identidad es siempre 1."
      }
    ],
    total_points: 30,
    passing_score: 18,
    time_limit_minutes: 25,
    created_at: new Date(),
    updated_at: new Date()
  },

  // Assessment de Biología Celular
  {
    material_id: "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeea",
    title: "Evaluación: La Célula",
    description: "Evaluación sobre estructura y función celular",
    questions: [
      {
        id: "q1",
        type: "multiple_choice",
        question: "¿Cuál organelo es conocido como la 'central energética' de la célula?",
        options: [
          "Núcleo",
          "Ribosoma",
          "Mitocondria",
          "Aparato de Golgi"
        ],
        correct_answer: 2,
        explanation: "La mitocondria produce ATP, la principal fuente de energía celular."
      },
      {
        id: "q2",
        type: "true_false",
        question: "Las células procariotas tienen núcleo definido.",
        correct_answer: false,
        explanation: "Las procariotas NO tienen núcleo definido; su material genético está en el citoplasma."
      },
      {
        id: "q3",
        type: "multiple_choice",
        question: "¿Dónde se encuentra el ADN en una célula eucariota?",
        options: [
          "En el citoplasma",
          "En el núcleo",
          "En la membrana celular",
          "En los ribosomas"
        ],
        correct_answer: 1,
        explanation: "En células eucariotas, el ADN se encuentra dentro del núcleo."
      }
    ],
    total_points: 30,
    passing_score: 18,
    time_limit_minutes: 15,
    created_at: new Date(),
    updated_at: new Date()
  }
]);

print("✅ assessments insertados exitosamente");
