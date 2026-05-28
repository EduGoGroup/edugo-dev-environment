package com.edugo.kmp.screens.dynamic.contracts

import com.edugo.kmp.sdui.engine.contract.CustomEventHandler
import com.edugo.kmp.sdui.engine.contract.EventContext
import com.edugo.kmp.sdui.engine.contract.EventResult

// Fixture happy: contrato canónico de un form CRUD.
class SchoolsFormContract : BaseCrudContract(
    apiPrefix = "academic:",
    basePath = "/api/v1/schools",
    resource = "schools"
) {
    override val screenKey = "schools-form"

    private class CreateHandler : CustomEventHandler {
        override val eventId = "save-new"
        override val requiredPermission: String? = "schools:create"
    }

    private class UpdateHandler : CustomEventHandler {
        override val eventId = "save-existing"
        override val requiredPermission: String? = "schools:update"
    }
}
