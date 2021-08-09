package model

type PayloadKey int32

const (
	PKSourceBranch    PayloadKey = iota // Исходная ветка
	PKTargetBranch                      // Целевая ветка
	PKHeader                            // Имя (прим. названия МРа)
	PKTriggeredByUser                   // Имя пользователя, который инициировал хук
	PKCreatedByUser                     // Пользователь создавший (мр, изменения в ветку и тд)

	PKmrAssignedToUser // Пользователь, на которого назначили МР
	PKmrApprovedBy     // Пользователь, который одобрил МР
	PKmrClosedBy       // Пользователь, который закрыл МР
	PKmrMergedBy       // Пользователь, сливший МР
	PKmrUpdatedBy      // Пользователь, который обновил МР
)

const (
	MRSubInfoKey_AssignedTo = "assigned_to"
	MRSubInfoKey_ApprovedBy = "approved_by"
	MRSubInfoKey_ClosedBy   = "closed_by"
	MRSubInfoKey_MergedBy   = "merged_by"
	MRSubInfoKey_UpdatedBy  = "updated_by"
	MRSubInfoKey_Link       = "mr_link"
	MRSubInfoKey_PipeLink   = "pipeline_link"
)
