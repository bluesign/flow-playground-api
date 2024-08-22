// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"github.com/Masterminds/semver"
	"github.com/google/uuid"
)

type Event struct {
	Type   string   `json:"type"`
	Values []string `json:"values"`
}

type Mutation struct {
}

type NewContractDeployment struct {
	ProjectID uuid.UUID `json:"projectId"`
	Script    string    `json:"script"`
	Address   Address   `json:"address"`
	Arguments []string  `json:"arguments,omitempty"`
}

type NewContractTemplate struct {
	ProjectID uuid.UUID `json:"projectId"`
	Title     string    `json:"title"`
	Script    string    `json:"script"`
}

type NewFile struct {
	ProjectID uuid.UUID `json:"projectId"`
	Title     string    `json:"title"`
	Script    string    `json:"script"`
}

type NewProject struct {
	ParentID             *uuid.UUID                       `json:"parentId,omitempty"`
	Title                string                           `json:"title"`
	Description          string                           `json:"description"`
	Readme               string                           `json:"readme"`
	Seed                 int                              `json:"seed"`
	NumberOfAccounts     int                              `json:"numberOfAccounts"`
	TransactionTemplates []*NewProjectTransactionTemplate `json:"transactionTemplates,omitempty"`
	ScriptTemplates      []*NewProjectScriptTemplate      `json:"scriptTemplates,omitempty"`
	ContractTemplates    []*NewProjectContractTemplate    `json:"contractTemplates,omitempty"`
}

type NewProjectContractTemplate struct {
	Title  string `json:"title"`
	Script string `json:"script"`
}

type NewProjectFile struct {
	Title  string `json:"title"`
	Script string `json:"script"`
}

type NewProjectScriptTemplate struct {
	Title  string `json:"title"`
	Script string `json:"script"`
}

type NewProjectTransactionTemplate struct {
	Title  string `json:"title"`
	Script string `json:"script"`
}

type NewScriptExecution struct {
	ProjectID uuid.UUID `json:"projectId"`
	Script    string    `json:"script"`
	Arguments []string  `json:"arguments,omitempty"`
}

type NewScriptTemplate struct {
	ProjectID uuid.UUID `json:"projectId"`
	Title     string    `json:"title"`
	Script    string    `json:"script"`
}

type NewTransactionExecution struct {
	ProjectID uuid.UUID `json:"projectId"`
	Script    string    `json:"script"`
	Signers   []Address `json:"signers,omitempty"`
	Arguments []string  `json:"arguments,omitempty"`
}

type NewTransactionTemplate struct {
	ProjectID uuid.UUID `json:"projectId"`
	Title     string    `json:"title"`
	Script    string    `json:"script"`
}

type PlaygroundInfo struct {
	APIVersion      semver.Version `json:"apiVersion"`
	CadenceVersion  semver.Version `json:"cadenceVersion"`
	EmulatorVersion semver.Version `json:"emulatorVersion"`
}

type ProgramError struct {
	Message       string           `json:"message"`
	StartPosition *ProgramPosition `json:"startPosition,omitempty"`
	EndPosition   *ProgramPosition `json:"endPosition,omitempty"`
}

type ProgramPosition struct {
	Offset int `json:"offset"`
	Line   int `json:"line"`
	Column int `json:"column"`
}

type ProjectList struct {
	Projects []*Project `json:"projects,omitempty"`
}

type Query struct {
}

type UpdateContractTemplate struct {
	ID        uuid.UUID `json:"id"`
	Title     *string   `json:"title,omitempty"`
	ProjectID uuid.UUID `json:"projectId"`
	Index     *int      `json:"index,omitempty"`
	Script    *string   `json:"script,omitempty"`
}

type UpdateFile struct {
	ID        uuid.UUID `json:"id"`
	Title     *string   `json:"title,omitempty"`
	ProjectID uuid.UUID `json:"projectId"`
	Index     *int      `json:"index,omitempty"`
	Script    *string   `json:"script,omitempty"`
}

type UpdateProject struct {
	ID          uuid.UUID `json:"id"`
	Title       *string   `json:"title,omitempty"`
	Description *string   `json:"description,omitempty"`
	Readme      *string   `json:"readme,omitempty"`
	Persist     *bool     `json:"persist,omitempty"`
}

type UpdateScriptTemplate struct {
	ID        uuid.UUID `json:"id"`
	Title     *string   `json:"title,omitempty"`
	ProjectID uuid.UUID `json:"projectId"`
	Index     *int      `json:"index,omitempty"`
	Script    *string   `json:"script,omitempty"`
}

type UpdateTransactionTemplate struct {
	ID        uuid.UUID `json:"id"`
	Title     *string   `json:"title,omitempty"`
	ProjectID uuid.UUID `json:"projectId"`
	Index     *int      `json:"index,omitempty"`
	Script    *string   `json:"script,omitempty"`
}
