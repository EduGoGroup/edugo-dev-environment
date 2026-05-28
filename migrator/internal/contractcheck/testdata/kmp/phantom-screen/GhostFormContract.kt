package com.edugo.kmp.screens.dynamic.contracts

// Fixture phantom-screen: declara un screenKey que el seed no contempla.
class GhostFormContract : BaseCrudContract(
    apiPrefix = "academic:",
    basePath = "/api/v1/ghosts",
    resource = "ghosts"
) {
    override val screenKey = "ghost-form"
}
