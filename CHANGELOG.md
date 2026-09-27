# Changelog

Module github.com/tmc/a2ui stays at v0.x until the upstream A2UI v1.0
specification is marked stable and this API has settled. A future
protocol 2.0 will be module github.com/tmc/a2ui/v2.

## v0.5.0

Function evaluation and component checks in a2uistate.

### Added

- a2uistate.Evaluator resolves Dynamic* values against a DataModel and
  evaluates function calls. Its Funcs map names to Func
  implementations; nil means BasicFunctions. A non-nil map replaces the
  basic functions entirely, so extend the map that BasicFunctions
  returns. Arguments are resolved recursively, including arrays and
  nested calls, before the function is called. The catalog ID of a
  call is ignored.
  - ResolveValue, ResolveString, ResolveNumber, ResolveBoolean and
    ResolveStringList return (T, error).
  - ResolveArgs resolves the arguments of a call without calling it,
    for renderers that perform actions.
  - Check evaluates a component's []a2ui.CheckRule and returns the
    failing a2ui.ValidationResults in order. The message comes from the
    validation result, then the rule, then "Validation failed". The
    severity defaults to SeverityError. A condition bound to a missing
    value fails. Evaluation errors fail the check and are returned
    joined.
- a2uistate.BasicFunctions returns a new map with the basic catalog
  functions and @index: formatString, formatNumber, formatCurrency,
  formatDate, pluralize, required, regex, length, numeric, email, and,
  or, not and openUrl. formatString parses ${...} templates as web_core
  does, including nested calls, relative paths and ${@index(offset: 1)}.
  Formatting is deterministic, en-US only and stdlib only. The
  BasicFunctions documentation lists where it differs from web_core:
  - A fixed table of 22 CLDR currency symbols.
  - formatDate accepts ISO 8601 only. It reads its pattern in runs of
    the same letter, as Unicode TR35 does, so EEE is "Mon" (web_core
    gives "MonMonMon"). Letters that are not fields, such as YYYY, are
    copied as they are.
  - regex uses RE2.
  - Validation results carry no default message.
  - and, or and not accept validation results.
  - numeric accepts numeric strings.
  - Unknown arguments are ignored.
- formatCurrency reports an empty currency as ErrInvalidArgs.
- Sentinel errors for errors.Is: ErrNoValue, ErrWrongType,
  ErrUnknownFunction, ErrInvalidArgs and ErrAction. Messages start
  with "a2uistate:".
- Actions are never performed during resolution. Evaluating openUrl
  fails with ErrAction.

### Changed

- The generated function-call helpers follow the catalog schema.
  Required parameters come first, then optional ones, each in the order
  of the schema's args.properties.
  Optional plain scalars are pointers, and nil leaves them out, as in
  a2ui.Length(v, new(8), nil). Optional Dynamic* arguments are left out
  when they are zero. Required arguments are always sent. Previously,
  the parameters were in alphabetical order, Length and Numeric always
  sent "min": 0 and "max": 0, and unset Dynamic* arguments failed to
  marshal. These signatures changed:
  - Length(max int, min int, value DynamicString) is now
    Length(value DynamicString, min *int, max *int).
  - Numeric(max float64, min float64, value DynamicNumber) is now
    Numeric(value DynamicNumber, min *float64, max *float64).
  - Pluralize(few, many, one, other, two DynamicString,
    value DynamicNumber, zero DynamicString) is now
    Pluralize(value DynamicNumber, other, zero, one, two, few,
    many DynamicString).
  - FormatCurrency(currency DynamicString, decimals DynamicNumber,
    grouping DynamicBoolean, value DynamicNumber) is now
    FormatCurrency(value DynamicNumber, currency DynamicString,
    decimals DynamicNumber, grouping DynamicBoolean).
  - FormatNumber(decimals DynamicNumber, grouping DynamicBoolean,
    value DynamicNumber) is now FormatNumber(value DynamicNumber,
    decimals DynamicNumber, grouping DynamicBoolean).
  - FormatDate(format DynamicString, value DynamicValue) is now
    FormatDate(value DynamicValue, format DynamicString).
  - Regex(pattern string, value DynamicString) is now
    Regex(value DynamicString, pattern string).
  And, Email, FormatString, Not, OpenURL, Or and Required are unchanged.
- The DataModel Resolve methods keep their (T, bool) signatures but now
  evaluate calls to the basic functions, and only those. They
  previously reported false for any function call. A call that fails,
  including an unknown function or an action, still reports false and
  the error is dropped. Use the Evaluator methods for other functions
  and for the error.
- ResolveValue resolves the dynamic values inside an array.
- ResolveString no longer escapes <, > and & when it converts objects
  and arrays to JSON, which matches JSON.stringify.

## v0.4.0

Renderer-side helpers and stricter validation.

### Added

- Package a2uistate holds the renderer state of 1.x surfaces:
  - DataModel: a JSON value addressed by JSON Pointers. Get and Set
    implement updateDataModel: intermediate objects and arrays are
    created, nil deletes, "" and "/" address the whole model, and array
    indexes are capped at 10000.
  - ResolvePath, plus ResolveValue, ResolveString, ResolveNumber,
    ResolveBoolean and ResolveStringList on DataModel. They resolve
    Dynamic* values in a list-template scope. ResolveString formats
    numbers as JavaScript's String does (42, 1.5, 1e+21, 1e-7). Function
    calls are not evaluated.
  - Surface (NewSurface, Apply, Component, Root, Data, CatalogID,
    Metadata, SendDataModel, Deleted) and Surfaces (Apply, Surface, Add).
    They route messages by surface ID. It is an error for a message to
    address the wrong surface, an unknown surface or a deleted surface,
    or to repeat createSurface. A failed createSurface leaves the surface
    unchanged. A deleteSurface through Surfaces also marks a retained
    *Surface deleted. The zero Surface is usable, with the ID "".
  - The a2uistate types are not safe for concurrent use.
- a2ui.BasicCatalogID is the basic catalog ID, generated from the
  catalog.
- An example of new(a2ui.StringLiteral("x")) for optional pointer
  fields. No pointer-helper functions were added.

### Changed

- a2uischema: ParseMessages, ValidateJSON and ValidateExample reject
  fields the message types do not define, whatever their value
  (ErrInvalidMessage, with a JSON Pointer path). This includes
  misspelled fields and the returnType of function calls, which was
  removed in 1.x. The check looks up each field in the Go type decoded
  at its position. The data model, function arguments and custom
  component properties still accept any fields. The a2ui package's
  json decoding stays lenient.
- a2uistream: ParseAndValidate calls the validator's
  ValidateJSON([]byte) error if it has one, so that
  a2uischema.Validator also checks streamed responses for unknown
  fields. A nil validator still checks only versions and decoding.

## v0.3.0

The module now supports A2UI 1.x only. The root package holds the 1.x
types, generated from specification/v1_0 and catalogs/basic/v1 in
[a2ui-project/a2ui](https://github.com/a2ui-project/a2ui). Programs that
need A2UI v0.9 should stay on v0.2.0.

### Removed

- The v09 and v091 packages, and all v0.8, v0.9 and v0.9.1 support.
- a2uistream: V08Parser, V08Reader, V08Part and NewV08Parser/NewV08Reader.
- a2uischema:
  - Version, Version09, Version091 and Catalog.Version.
  - ValidateVersionMessages.
  - ValidationCode and its constants.
  - The Code, Component, Ref, Function and Message fields of ValidationError.
- a2a:
  - Part, Extension, CreatePart, NewExtension, Data and IsPart. Use DataPart,
    AgentExtension, CreateDataPart, NewAgentExtension, A2UIData and
    IsA2UIPart.
  - MIMETypeForVersion, CreateDataPartForVersion, Versioned and
    IsA2UIMIMEType.
  - The versioned MIME type constants. The v0.9 type
    "application/json+a2ui" is no longer recognized.

### Changed

- Root package: the v0.9 names (ServerMessage, ClientCapabilities and so
  on) are replaced by the 1.x types: AgentMessage, RendererMessage,
  RendererCapabilities, FunctionResponse and the rest.
- IconNameOrPath is {Name, SVGPath, Binding}. It covers all three v1.0
  forms: an icon name, {"svgPath": ...} and a {"path": ...} data binding.
- FunctionResponse has no HasValue field. A response with a nil Error is a
  value response, and a nil Value marshals as "value": null.
- AgentMessage and RendererMessage marshal an empty Version as
  a2ui.Version. The VersionString methods are removed.
- Union types (DynamicString, DynamicValue, Action, IconNameOrPath and
  the rest) treat JSON null as a no-op when decoding. A zero union fails
  to marshal, instead of encoding as "". A zero ChildList still encodes
  as [].
- a2uistream: ResponsePart.Messages is []a2ui.AgentMessage.
  ParseAndValidate, Parser and Reader decode 1.x only, and
  ParseAndValidate returns an error for content of any other version.
- a2uischema:
  - NewSchemaManager(catalogs, acceptsInlineCatalogs, modifiers...) takes
    no version.
  - BasicCatalogConfig() and BasicCatalogProvider() take no version and
    return no error.
  - ParseMessages and ValidateMessages use []a2ui.AgentMessage.
  - SelectedCatalog takes *a2ui.RendererCapabilities instead of any.
  - GenerateSystemPrompt takes a PromptOptions struct instead of nine
    positional arguments.
  - Catalog.ServerToClientSchema is renamed MessageSchema.
  - Error messages start with "a2uischema:".
  - ValidateMessages requires a root component only for surfaces created
    in the same messages. Updates to other surfaces may omit root and
    refer to components sent earlier, as the spec allows.
- a2a: A2UIMIMEType is "application/a2ui+json", and extension versions
  default to v1.0. MarshalA2UIData deep-copies map payloads.
- a2uiadk: the ToolContext parameter of SendA2UIJSONToClientTool.Run is
  named tc.

### Added

- a2ui: Component.Custom and CustomComponent{Type, Properties}, for
  components from custom and inline catalogs.
- a2uischema:
  - The sentinel errors ErrInvalidMessage, ErrVersionMismatch,
    ErrUnknownComponent, ErrUnknownFunction, ErrInvalidTree,
    ErrUnallowedParent and ErrUnallowedChild, for use with errors.Is.
    The last two match the spec's UNALLOWED_PARENT and UNALLOWED_CHILD
    codes.
  - ValidationError{Path, Err}, for use with errors.As. Path is a JSON
    pointer to the offending value, such as
    /1/updateComponents/components/0/child.
  - Validation enforces the allowedParents and allowedChildren lists of
    v1.0 catalogs. The basic catalog declares neither.
  - Custom components are checked against the catalog, and a type the
    catalog does not define is an ErrUnknownComponent.
  - Inline catalog components and functions are listed in the merged
    schema's $defs.anyComponent and $defs.anyFunction.
- a2uistream: ErrInvalidPayload.

## v0.2.0

The last release with A2UI v0.9.

- Every name in the root package is deprecated. Each one forwards to the
  same name in github.com/tmc/a2ui/v09 and is marked //go:fix inline.
  Running `go fix ./...` rewrites callers to import v09 directly.

  With Go 1.27, go fix can silently write nothing for a package where two
  fixes overlap. If a file does not change, run the inline analyzer
  twice instead:

      go run golang.org/x/tools/go/analysis/passes/inline/cmd/inline@latest -fix ./...

- The v010 package is removed. The A2UI v0.10 draft was folded into
  v1.0 upstream.
- a2uibuild, a2uischema and a2uistream import v09 directly.

## v0.1.0

First release. The root package re-exports A2UI v0.9 from v09, with the
v09, v091 and v010 packages and the a2a, a2uiadk, a2uibuild, a2uischema
and a2uistream helpers.
