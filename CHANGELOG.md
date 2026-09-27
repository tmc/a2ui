# Changelog

Module github.com/tmc/a2ui stays at v0.x until the upstream A2UI v1.0
specification is marked stable and this API has settled. A future
protocol 2.0 will be module github.com/tmc/a2ui/v2.

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
- a2a: A2UIMIMEType is "application/a2ui+json", and extension versions
  default to v1.0.
- a2uiadk: the ToolContext parameter of SendA2UIJSONToClientTool.Run is
  named tc.

### Added

- a2uischema:
  - The sentinel errors ErrInvalidMessage, ErrVersionMismatch,
    ErrUnknownComponent, ErrUnknownFunction, ErrInvalidTree and
    ErrNotAllowed, for use with errors.Is.
  - ValidationError{Path, Err}, for use with errors.As. Path is a JSON
    pointer to the offending value, such as
    /1/updateComponents/components/0/child.
  - Validation enforces the allowedParents and allowedChildren lists of
    v1.0 catalogs. The basic catalog declares neither.
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
