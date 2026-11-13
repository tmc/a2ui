package a2ui

import "strings"

// DataBinding references a value in the data model by JSON Pointer path.
// It is encoded as {"@path": ...}.
type DataBinding struct {
	Path string `json:"@path"`
}

// FunctionCall invokes a named catalog function.
// It is encoded as {"@call": ..., "args": ...}.
// CatalogID overrides the surface's default catalog.
//
// The values in Args are dynamic. They are Go values such as a
// [DynamicString], or JSON values as on the wire, in which an object
// with an "@path" key is a data binding, one with an "@call" key is a
// function call, and in any other object a key with a literal leading
// "@" is escaped as described for [EscapeKey].
type FunctionCall struct {
	Call      string         `json:"@call"`
	CatalogID string         `json:"catalogId,omitempty"`
	Args      map[string]any `json:"args,omitempty"`
}

// IsReservedKey reports whether key is reserved for protocol
// directives in a dynamic object: whether it starts with a single "@",
// as "@path", "@call" and "@" do, but "@@path" does not.
//
// In an object that is evaluated as a dynamic value, such as a
// [DynamicValue] object or an object in [FunctionCall.Args], only the
// directives "@path" and "@call" may have a reserved key. Any other
// reserved key is an error.
func IsReservedKey(key string) bool {
	return strings.HasPrefix(key, "@") && !strings.HasPrefix(key, "@@")
}

// EscapeKey returns key escaped for use in an object that is evaluated
// as a dynamic value. A key with a leading "@" is escaped by doubling
// the "@", so that "@type" becomes "@@type"; other keys are unchanged.
func EscapeKey(key string) string {
	if strings.HasPrefix(key, "@") {
		return "@" + key
	}
	return key
}

// UnescapeKey reverses [EscapeKey]: it removes one "@" from a key that
// starts with "@@", so that "@@type" becomes "@type". Other keys are
// unchanged.
func UnescapeKey(key string) string {
	if strings.HasPrefix(key, "@@") {
		return key[1:]
	}
	return key
}

// ChildList is either a static list of component IDs or a dynamic template.
// Exactly one of IDs or Template is set.
type ChildList struct {
	IDs      []string
	Template *ChildTemplate
}

// ChildTemplate generates a dynamic list of children from a data model list.
type ChildTemplate struct {
	ComponentID string `json:"componentId"`
	Path        string `json:"path"`
}

// CheckRule is a single validation check applied to an input component.
// Condition evaluates to a [ValidationResult]; Message is a fallback
// error message.
type CheckRule struct {
	Condition DynamicValidationResult `json:"condition"`
	Message   string                  `json:"message,omitempty"`
}

// AccessibilityAttributes enhance accessibility for assistive technologies.
type AccessibilityAttributes struct {
	Label       *DynamicString  `json:"label,omitempty"`
	Description *DynamicString  `json:"description,omitempty"`
	Live        AccessibleLive  `json:"live,omitempty"`
	Hidden      *DynamicBoolean `json:"hidden,omitempty"`
}

// AccessibleLive controls screen reader announcements for dynamic updates.
type AccessibleLive string

const (
	AccessibleLiveOff       AccessibleLive = "off"
	AccessibleLivePolite    AccessibleLive = "polite"
	AccessibleLiveAssertive AccessibleLive = "assertive"
)

// Metadata carries extension metadata on surfaces, components, and actions.
// Keys starting with "a2ui_" are reserved for official extensions.
type Metadata struct {
	Extensions map[string]any `json:"extensions,omitempty"`
}

// Action is an interaction handler that either triggers an agent-side event
// or executes a function. Exactly one field is non-nil.
type Action struct {
	Event        *EventAction  `json:"event,omitempty"`
	FunctionCall *FunctionCall `json:"functionCall,omitempty"`
}

// EventAction triggers an agent-side event.
type EventAction struct {
	Name        string                  `json:"name"`
	UserMessage *DynamicString          `json:"userMessage,omitempty"`
	Context     map[string]DynamicValue `json:"context,omitempty"`
}

// IconNameOrPath is a well-known icon name, a custom SVG path, or a
// binding to an icon name in the data model. Exactly one field is non-nil.
// The three forms are encoded as a string, {"svgPath": ...} and
// {"@path": ...}.
type IconNameOrPath struct {
	Name    *IconName
	SVGPath *DynamicString
	Binding *DataBinding
}

// Index returns a call to the @index system function, the 0-based index
// of the current item in a template list, plus offset.
func Index(offset float64) DynamicNumber {
	call := FunctionCall{Call: "@index"}
	if offset != 0 {
		call.Args = map[string]any{"offset": offset}
	}
	return NumberFunc(call)
}
