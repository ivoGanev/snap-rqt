package input

type Action string

const (
	ActionAddRequest          Action = "add_request"
	ActionRemoveRequest       Action = "remove_request"
	ActionDuplicateRequest    Action = "duplicate_request"
	ActionEditRequestName     Action = "edit_request_name"
	ActionAddCollection       Action = "add_collection"
	ActionRemoveCollection    Action = "remove_collection"
	ActionDuplicateCollection Action = "duplicate_collection"
	ActionEditCollectionName  Action = "edit_collection_name"
	ActionFocusCollections    Action = "focus_collections"
	ActionFocusRequests       Action = "focus_requests"
	ActionToggleViewMode      Action = "toggle_view_mode"
	ActionQuit                Action = "quit"
	ActionSwapFocus           Action = "swap_focus"
	ActionModalSave           Action = "modal_save"
	ActionModalCancel         Action = "modal_cancel"
	ActionSwitchToBody        Action = "switch_to_body"
	ActionSwitchToHeaders     Action = "switch_to_headers"
)
