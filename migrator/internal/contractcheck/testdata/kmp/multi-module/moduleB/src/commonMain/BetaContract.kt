package com.edugo.kmp.modules.beta

class BetaContract : BaseCrudContract(
    apiPrefix = "learning:",
    basePath = "/api/v1/beta",
    resource = "beta"
) {
    override val screenKey = "beta-list"
}
