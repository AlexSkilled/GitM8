package model

type PayloadKey int32

const (
	PKSourceBranch    PayloadKey = iota // Исходная ветка
	PKTargetBranch                      // Целевая ветка
	PKHeader                            // Имя (прим. названия МРа)
	PKTriggeredByUser                   // Имя пользователя, который инициировал хук
	PKCreatedByUser                     // Пользователь создавший (мр, изменения в ветку и тд)

	PKLink // Ссылка, пока используем как ссылка на МР

	PKmrAssignedToUser // Пользователь, на которого назначили МР
	PKmrClosedBy       // Пользователь, который закрыл МР
	PKmrUpdatedBy      // Пользователь, который обновил МР

)

const (
	MRPattern_Branches   = "branches"
	MRPattern_OpenedBy   = "opened_by"
	MRPattern_AssignedTo = "assigned_to"
	MRPattern_ApprovedBy = "approved_by"
	MRPattern_ClosedBy   = "closed_by"
	MRPattern_MergedBy   = "merged_by"
	MRPattern_UpdatedBy  = "updated_by"
	MRPattern_PipeLink   = "pipeline_link"
)
