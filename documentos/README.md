# Documentación - EduGo Dev Environment

Este directorio contiene la documentación actualizada del proyecto `edugo-dev-environment`.

## Índice de Documentos

| Documento | Descripción |
|-----------|-------------|
| [ARQUITECTURA.md](./ARQUITECTURA.md) | Arquitectura del sistema, componentes y flujo de datos |
| [SERVICIOS.md](./SERVICIOS.md) | Detalle de cada servicio Docker y su configuración |
| [GUIA-RAPIDA.md](./GUIA-RAPIDA.md) | Pasos para levantar el ambiente en minutos |
| [FAQ.md](./FAQ.md) | Preguntas frecuentes y troubleshooting |
| [DEPRECADO-MEJORAS.md](./DEPRECADO-MEJORAS.md) | Código deprecado, mejoras pendientes y deuda técnica |

## Propósito del Proyecto

`edugo-dev-environment` es un repositorio puente que permite a los desarrolladores (especialmente frontend) levantar todo el backend de EduGo localmente con Docker, sin necesidad de conocimientos de Go o backend.

## Estructura del Proyecto

```
edugo-dev-environment/
├── docker/                    # Configuración Docker Compose
│   ├── docker-compose.yml     # Configuración principal
│   ├── .env.example           # Variables de entorno de ejemplo
│   └── *.yml                  # Configuraciones alternativas
├── scripts/                   # Scripts de utilidad
│   ├── setup.sh               # Setup inicial
│   ├── validate.sh            # Validación de configuración
│   └── ...
├── migrator/                  # Herramienta de migraciones
├── seeds/                     # Datos de prueba
│   ├── postgresql/
│   └── mongodb/
├── documentos/                # Documentación (este directorio)
└── README.md                  # Guía principal
```

---

**Última actualización:** Diciembre 2025
