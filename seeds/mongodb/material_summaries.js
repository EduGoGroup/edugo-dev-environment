// Seeds de Material Summaries
// Ejemplo de estructura

db.material_summaries.insertMany([
  {
    material_id: "example-material-1",
    summary: "Este es un resumen de ejemplo generado por IA",
    key_points: ["Punto 1", "Punto 2", "Punto 3"],
    created_at: new Date()
  }
]);

print("Seeds de material_summaries cargados");
