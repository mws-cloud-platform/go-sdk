// Package interfaces provides interfaces for resource references and IDs.
package interfaces //nolint:revive // leave a name for backward compatibility

// ResourceName is the type used in IDs and refs.
type ResourceName string

// ResourceRef is an interface for resource references.
type ResourceRef interface {
	// ServiceSlug returns a unique, URL-friendly identifier for the service.
	// This is typically a short string used in API paths and URLs.
	// Example: "compute", "vpc".
	ServiceSlug() string
	// IDPath returns an absolute path to the resource including the service slug.
	// Example: "compute/projects/my-project/vms/my-vm".
	IDPath() string
	// Path returns the original resource path as provided during creation.
	// Unlike IDPath, this preserves the original format and may be relative or absolute.
	// Example: "vms/my-vm" or "projects/my-project/vms/my-vm" or "compute/projects/my-project/vms/my-vm".
	Path() string
	// ResourceName returns the name of the resource this reference points to.
	// Example: "my-vm".
	ResourceName() ResourceName
}

// ResourceID is an interface for resource IDs.
type ResourceID interface {
	// ServiceSlug returns a unique, URL-friendly identifier for the service.
	// This is typically a short string used in API paths and URLs.
	// Example: "compute", "vpc".
	ServiceSlug() string
	// ID returns an absolute path to the resource including the service slug.
	// Example: "compute/projects/my-project/vms/my-vm".
	ID() string
	// ResourceName returns the name of the resource this reference points to.
	// Example: "my-vm".
	ResourceName() ResourceName
}
