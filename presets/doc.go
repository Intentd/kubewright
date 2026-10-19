// Package presets names the field settings teams repeat: resource requests
// and limits, probe timings, container and pod security, rolling update,
// HPA and PDB behavior, network policy rules and namespace quotas. A spec
// field names a preset ("medium", "restricted", "safe") instead of spelling
// the settings out.
//
// Builtin returns the built-in presets. A Registry merges your own with
// With, where yours replace built-ins of the same name, and looks them up by
// name, ignoring case.
//
// An Expand function per kind of preset turns a named preset into explicit
// fields, and returns an error for a name the registry does not have:
//
//	reg := presets.Builtin().With(own)
//	res, err := presets.ExpandResources(&types.ResourceRequirements{Preset: "medium"}, reg)
//
// Expand presets before calling a verb: Build never looks one up, and Create
// and Apply refuse a container that still names one.
package presets
