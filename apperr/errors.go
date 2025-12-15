package apperr

import (
    "errors"
)

// Kind classifies an application error.
type Kind string

const (
    // Invalid indicates the client supplied invalid input.
    Invalid Kind = "invalid"
    // Required indicates a missing required field.
    Required Kind = "required"
    // NotFound indicates the requested resource does not exist.
    NotFound Kind = "not_found"
    // AlreadyExists indicates a unique constraint violation on create/update.
    AlreadyExists Kind = "already_exists"
    // Conflict indicates a state conflict.
    Conflict Kind = "conflict"
    // Unauthorized indicates missing/invalid authentication.
    Unauthorized Kind = "unauthorized"
    // Forbidden indicates authenticated but not permitted.
    Forbidden Kind = "forbidden"
    // Internal indicates an unexpected server error.
    Internal Kind = "internal"
)

// Error is the canonical application error type.
type Error struct {
    Kind   Kind  // machine-readable classification
    Entity string // optional: affected domain entity (e.g., "identity")
    Field  string // optional: which field is invalid/required/unique
    Msg    string // human-friendly message
    Err    error  // wrapped underlying error
}

func (e *Error) Error() string {
    if e == nil {
        return "<nil>"
    }
    if e.Msg != "" {
        return e.Msg
    }
    if e.Err != nil {
        return e.Err.Error()
    }
    return string(e.Kind)
}

// Unwrap supports errors.Is / errors.As.
func (e *Error) Unwrap() error { return e.Err }

// Helper constructors

func New(kind Kind, msg string) *Error { return &Error{Kind: kind, Msg: msg} }

func Wrap(err error, kind Kind, msg string) *Error { return &Error{Kind: kind, Msg: msg, Err: err} }

func RequiredField(field string) *Error { return &Error{Kind: Required, Field: field, Msg: field + " is required"} }

func InvalidField(field, msg string) *Error { return &Error{Kind: Invalid, Field: field, Msg: msg} }

func EntityNotFound(entity string) *Error { return &Error{Kind: NotFound, Entity: entity, Msg: entity + " not found"} }

func UniqueConstraint(entity, field string) *Error {
    return &Error{Kind: AlreadyExists, Entity: entity, Field: field, Msg: entity + " with this " + field + " already exists"}
}

func InternalError(err error, msg string) *Error { return &Error{Kind: Internal, Msg: msg, Err: err} }

// Kind checks

func IsKind(err error, k Kind) bool {
    var ae *Error
    if errors.As(err, &ae) {
        return ae.Kind == k
    }
    return false
}

func IsNotFound(err error) bool      { return IsKind(err, NotFound) }
func IsAlreadyExists(err error) bool { return IsKind(err, AlreadyExists) }
func IsInvalid(err error) bool       { return IsKind(err, Invalid) || IsKind(err, Required) }
