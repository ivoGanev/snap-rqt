package input

type Action string

const (
	ActionAddRequest                   Action = "add_request"
	ActionRemoveRequest                Action = "remove_request"
	ActionDuplicateRequest             Action = "duplicate_request"
	ActionEditRequestName              Action = "edit_request_name"
	ActionSelectRequest                Action = "select_request"
	ActionAddCollection                Action = "add_collection"
	ActionRemoveCollection             Action = "remove_collection"
	ActionDuplicateCollection          Action = "duplicate_collection"
	ActionEditCollectionName           Action = "edit_collection_name"
	ActionFocusCollections             Action = "focus_collections"
	ActionFocusRequests                Action = "focus_requests"
	ActionToggleCatalogMode            Action = "toggle_catalog_mode"
	ActionSwapPuppetModes              Action = "swap_puppet_modes"
	ActionQuit                         Action = "quit"
	ActionSwapFocus                    Action = "swap_focus"
	ActionModalSave                    Action = "modal_save"
	ActionModalCancel                  Action = "modal_cancel"
	ActionRequestEditorSwitchToBody    Action = "request_editor_switch_to_body"
	ActionRequestEditorSwitchToHeaders Action = "request_editor_switch_to_headers"
	ActionRequestEditorEdit            Action = "request_editor_edit"
	ActionRequestEditorDone            Action = "request_editor_input_done"
	ActionHeaderBarEditUrl             Action = "header_bar_edit_url"
	ActionHeaderBarSelectMethod        Action = "header_bar_select_method"
	ActionHeaderBarUrlApply            Action = "header_bar_url_apply"
	ActionLoseFocus                    Action = "lose_focus"
	ActionExitInputMode                Action = "exit_input_mode"
)
