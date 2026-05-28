package kmp

import "strings"

// stripComments elimina los comentarios Kotlin (// y /* */) preservando
// las posiciones de línea para que las regex que actúan sobre el texto
// "limpio" puedan reportar la línea correcta del archivo original.
//
// La función NO toca strings entre comillas dobles ("..."), de modo que
// un literal como "https://..." no se interprete como comentario.
//
// El segundo retorno es un slice paralelo a las líneas del texto
// limpio: lineMap[i] = número de línea (1-indexed) en el archivo
// original al que corresponde la línea i del texto saneado.
//
// Diseño: en vez de borrar las líneas comentadas, las reemplazamos por
// líneas vacías. Eso preserva el conteo de líneas y mantiene la
// correspondencia 1-a-1. Por eso lineMap[i] = i+1 cuando la entrada es
// "limpia". Mantenemos el lineMap explícito para futuros cambios donde
// quizás colapsemos líneas (no en esta versión).
func stripComments(src string) (string, []int) {
	var b strings.Builder
	b.Grow(len(src))
	inLineComment := false
	inBlockComment := false
	inString := false
	escape := false

	for i := 0; i < len(src); i++ {
		c := src[i]

		// Reset de comentario de línea al ver \n.
		if inLineComment {
			if c == '\n' {
				inLineComment = false
				b.WriteByte('\n')
			}
			// Dentro de comentario de línea: nada se escribe (excepto el \n).
			continue
		}

		if inBlockComment {
			// Dentro de /* ... */: preservamos los \n para no perder
			// el conteo de líneas.
			if c == '*' && i+1 < len(src) && src[i+1] == '/' {
				inBlockComment = false
				i++ // saltar el '/'
				continue
			}
			if c == '\n' {
				b.WriteByte('\n')
			}
			continue
		}

		if inString {
			b.WriteByte(c)
			if escape {
				escape = false
				continue
			}
			if c == '\\' {
				escape = true
				continue
			}
			if c == '"' {
				inString = false
			}
			continue
		}

		// Fuera de strings y comentarios.
		if c == '"' {
			inString = true
			b.WriteByte(c)
			continue
		}
		if c == '/' && i+1 < len(src) {
			next := src[i+1]
			if next == '/' {
				inLineComment = true
				i++ // saltar el segundo '/'
				continue
			}
			if next == '*' {
				inBlockComment = true
				i++ // saltar el '*'
				continue
			}
		}
		b.WriteByte(c)
	}

	cleaned := b.String()
	lines := strings.Split(cleaned, "\n")
	lineMap := make([]int, len(lines))
	for i := range lines {
		lineMap[i] = i + 1
	}
	return cleaned, lineMap
}

// lineForOffset devuelve el número de línea (1-indexed) en el texto
// limpio para un offset dado (en bytes).
func lineForOffset(cleaned string, offset int, lineMap []int) int {
	if offset < 0 {
		return 0
	}
	if offset > len(cleaned) {
		offset = len(cleaned)
	}
	line := 1 + strings.Count(cleaned[:offset], "\n")
	if line >= 1 && line <= len(lineMap) {
		return lineMap[line-1]
	}
	return line
}
