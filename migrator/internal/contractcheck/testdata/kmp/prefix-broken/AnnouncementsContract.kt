package com.edugo.kmp.screens.dynamic.contracts

// Reproduce el bug histórico F2·H3.a: el contrato declara apiPrefix
// "academic:" pero el resource "announcements" pertenece al servicio
// "platform:" según la serviceRoutingTable canónica. El cross-validator
// debe emitir un Drift category=service_prefix_mismatch severity=error.
class AnnouncementsContract : BaseCrudContract(
    apiPrefix = "academic:",
    basePath = "/api/v1/announcements",
    resource = "announcements"
) {
    override val screenKey = "announcements-list"
}
