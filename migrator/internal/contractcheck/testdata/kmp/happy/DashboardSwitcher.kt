package com.edugo.kmp.screens.ui

// Fixture happy: switcher que decide qué dashboard renderizar.
// Contiene literales "dashboard-<role>" y un when (role.code) { ... }.

@Composable
fun DashboardSwitcher(role: Role) {
    when (role.code) {
        "teacher" -> "dashboard-teacher"
        "student" -> "dashboard-student"
        "guardian" -> "dashboard-guardian"
        else -> "dashboard-home"
    }
}

// Comentario con un literal que NO debe ser extraído:
//   override val screenKey = "phantom-from-comment"
/* También bloqueado:
   override val screenKey = "phantom-from-block"
*/
