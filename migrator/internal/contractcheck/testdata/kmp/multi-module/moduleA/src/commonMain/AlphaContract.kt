package com.edugo.kmp.modules.alpha

class AlphaContract : BaseCrudContract(
    apiPrefix = "platform:",
    basePath = "/api/v1/alpha",
    resource = "alpha"
) {
    override val screenKey = "alpha-list"
}
