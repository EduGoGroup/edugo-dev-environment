# 🎮 MODO ULTRA FÁCIL ACTIVADO - Para Devs Frontend que les da pereza leer

## TL;DR para los impacientes:

```bash
# 1. ¿Tienes Docker? Ábrelo.
open -a Docker

# 2. Clona esto
git clone https://github.com/EduGoGroup/edugo-dev-environment.git
cd edugo-dev-environment

# 3. ¿Ya hiciste docker login ghcr.io? 
# ¿No? Hazlo. ¿Sí? Continúa.

# 4. Levanta todo (en serio, es UNO solo)
cd docker && docker-compose up -d

# 5. Espera 30 segundos mientras vas por café ☕

# 6. ¿Ya regresaste? Prueba esto:
curl http://localhost:8081/health
# Si ves {"status":"healthy"} → FELICIDADES, YA TERMINASTE
```

---

## 🤦 "Pero es que yo no sé usar Docker..."

Hermano/a, si sabes usar `npm install`, sabes usar Docker.  
Docker Desktop es literalmente un ícono que clickeas y ya.

---

## 🙄 "¿Y los datos de prueba dónde están?"

Ya están adentro, flojo/a. 8 usuarios, escuelas, cursos, todo.

```javascript
// Usuario de prueba (para que no digas que no te dimos)
const user = {
  email: 'admin@edugo.test',
  password: 'admin123' // Sí, es admin123. No, no es seguro. Es DESARROLLO.
}
```

---

## 😤 "Es que a mí no me funciona..."

Claro, porque no leíste. Aquí te lo mastico:

**Problema #1: "Cannot connect to Docker daemon"**  
👉 Abre Docker Desktop, genio.

**Problema #2: "pull access denied"**  
👉 `docker login ghcr.io` - Te lo dijimos en el paso 3.

**Problema #3: "Port already in use"**  
👉 Ya tienes algo corriendo en 8081. Apágalo.  
```bash
lsof -ti:8081 | xargs kill -9
```

**Problema #4: "No me funciona y no sé por qué"**  
👉 `docker-compose logs -f` - Lee los logs como adulto responsable.

---

## 💀 Excusas que NO aceptamos:

❌ "Es muy complicado"  
→ Son 4 comandos. CUATRO.

❌ "No tengo tiempo"  
→ Toma 2 minutos. Menos de lo que tardaste en quejarte en Slack.

❌ "No entiendo Docker"  
→ No necesitas entenderlo. Solo clickear el ícono.

❌ "¿Y si rompo algo?"  
→ Es un ambiente LOCAL. Si lo rompes, `docker-compose down -v` y vuelves a empezar.

❌ "Prefiero que backend me dé un ambiente en la nube"  
→ Este ES tu ambiente. Literalmente hicimos TODO el trabajo por ti.

---

## 🎯 Lo que REALMENTE necesitas saber:

```
API Mobile:  http://localhost:8081
API Admin:   http://localhost:8082
RabbitMQ UI: http://localhost:15672 (user: edugo / pass: edugo123)

Eso es TODO.
```

---

## 🏆 Si llegaste hasta aquí y AÚN no lo has levantado:

Eres oficialmente la persona más procrastinadora del equipo.  
Felicidades. 🎊

Ahora cierra este mensaje y ejecuta los 4 comandos.  
Tu yo del futuro te lo agradecerá cuando estés debuggeando a las 2am  
y necesites probar algo contra el backend.

---

## ✨ Bonus: Cómo impresionar a tu líder técnico

```bash
# Mientras tus compañeros preguntan "¿ya está el backend?"
# Tú ya estás desarrollando con datos reales

git clone https://github.com/EduGoGroup/edugo-dev-environment.git
cd edugo-dev-environment/docker
docker-compose up -d

# 30 segundos después:
"Ya terminé mi feature, ¿alguien más necesita ayuda?" 😎
```

---

**P.D.:** Si TODAVÍA tienes problemas después de leer esto,  
probablemente el problema no es el README. 🤷

**P.P.D.:** Todo el backend, todas las bases de datos, todo funcional,  
en tu laptop, sin internet, sin depender de nadie.  
Y aún así te quejas. Increíble.

---

💙 Con amor (y un poco de frustración),  
El equipo de Backend que hizo esto mientras ustedes discutían tabs vs spaces
