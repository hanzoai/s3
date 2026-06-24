// Code generated from plugin.proto; DO NOT EDIT.

package pluginwire

import (
	zap "github.com/zap-proto/go"
)

// This file carries the job-type descriptor / config-form cluster:
// JobTypeDescriptor, ConfigForm, ConfigSection, ConfigField, ConfigOption,
// ValidationRule, AdminRuntimeDefaults, AdminRuntimeConfig.

// --- JobTypeDescriptor ---

const (
	jobTypeDescriptorJobTypeOff              = 0
	jobTypeDescriptorDisplayNameOff          = 8
	jobTypeDescriptorDescriptionOff          = 16
	jobTypeDescriptorIconOff                 = 24
	jobTypeDescriptorDescriptorVersionOff    = 32
	jobTypeDescriptorAdminConfigFormOff      = 36
	jobTypeDescriptorWorkerConfigFormOff     = 44
	jobTypeDescriptorAdminRuntimeDefaultsOff = 52
	jobTypeDescriptorWorkerDefaultValuesOff  = 60
	jobTypeDescriptorSize                    = 68
)

// JobTypeDescriptor is a zero-copy view into a ZAP-encoded JobTypeDescriptor.
type JobTypeDescriptor struct{ o zap.Object }

// WrapJobTypeDescriptor parses b and returns a typed view.
func WrapJobTypeDescriptor(b []byte) (JobTypeDescriptor, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return JobTypeDescriptor{}, err
	}
	return JobTypeDescriptor{o: m.Root()}, nil
}

func (t JobTypeDescriptor) JobType() string     { return t.o.Text(jobTypeDescriptorJobTypeOff) }
func (t JobTypeDescriptor) DisplayName() string { return t.o.Text(jobTypeDescriptorDisplayNameOff) }
func (t JobTypeDescriptor) Description() string { return t.o.Text(jobTypeDescriptorDescriptionOff) }
func (t JobTypeDescriptor) Icon() string        { return t.o.Text(jobTypeDescriptorIconOff) }
func (t JobTypeDescriptor) DescriptorVersion() uint32 {
	return t.o.Uint32(jobTypeDescriptorDescriptorVersionOff)
}
func (t JobTypeDescriptor) AdminConfigForm() ConfigForm {
	return ConfigForm{o: childObject(t.o.Bytes(jobTypeDescriptorAdminConfigFormOff))}
}
func (t JobTypeDescriptor) WorkerConfigForm() ConfigForm {
	return ConfigForm{o: childObject(t.o.Bytes(jobTypeDescriptorWorkerConfigFormOff))}
}
func (t JobTypeDescriptor) AdminRuntimeDefaults() AdminRuntimeDefaults {
	return AdminRuntimeDefaults{o: childObject(t.o.Bytes(jobTypeDescriptorAdminRuntimeDefaultsOff))}
}

// WorkerDefaultValues reads worker_default_values (proto field 9,
// map<string,ConfigValue>); each element is a ConfigValueMapEntry object.
func (t JobTypeDescriptor) WorkerDefaultValues() zap.List {
	return t.o.List(jobTypeDescriptorWorkerDefaultValuesOff)
}

// JobTypeDescriptorInput collects the field values for NewJobTypeDescriptor.
// The form/defaults members take a pre-built sub-buffer; WorkerDefaultValues is
// a list of pre-built ConfigValueMapEntry sub-buffers.
type JobTypeDescriptorInput struct {
	JobType              string
	DisplayName          string
	Description          string
	Icon                 string
	DescriptorVersion    uint32
	AdminConfigForm      []byte
	WorkerConfigForm     []byte
	AdminRuntimeDefaults []byte
	WorkerDefaultValues  [][]byte
}

// NewJobTypeDescriptor builds a ZAP-encoded JobTypeDescriptor message.
func NewJobTypeDescriptor(in JobTypeDescriptorInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(jobTypeDescriptorSize)
	ob.SetText(jobTypeDescriptorJobTypeOff, in.JobType)
	ob.SetText(jobTypeDescriptorDisplayNameOff, in.DisplayName)
	ob.SetText(jobTypeDescriptorDescriptionOff, in.Description)
	ob.SetText(jobTypeDescriptorIconOff, in.Icon)
	ob.SetUint32(jobTypeDescriptorDescriptorVersionOff, in.DescriptorVersion)
	setChild(ob, jobTypeDescriptorAdminConfigFormOff, in.AdminConfigForm)
	setChild(ob, jobTypeDescriptorWorkerConfigFormOff, in.WorkerConfigForm)
	setChild(ob, jobTypeDescriptorAdminRuntimeDefaultsOff, in.AdminRuntimeDefaults)
	lb := b.StartList(0)
	for _, elem := range in.WorkerDefaultValues {
		lb.AddObjectBytes(elem)
	}
	ob.SetList(jobTypeDescriptorWorkerDefaultValuesOff, lb.FinishOffset(), len(in.WorkerDefaultValues))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ConfigForm ---

const (
	configFormFormIDOff        = 0
	configFormTitleOff         = 8
	configFormDescriptionOff   = 16
	configFormSectionsOff      = 24
	configFormDefaultValuesOff = 32
	configFormSize             = 40
)

// ConfigForm is a zero-copy view into a ZAP-encoded ConfigForm message.
type ConfigForm struct{ o zap.Object }

// WrapConfigForm parses b and returns a typed view.
func WrapConfigForm(b []byte) (ConfigForm, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ConfigForm{}, err
	}
	return ConfigForm{o: m.Root()}, nil
}

func (t ConfigForm) FormID() string      { return t.o.Text(configFormFormIDOff) }
func (t ConfigForm) Title() string       { return t.o.Text(configFormTitleOff) }
func (t ConfigForm) Description() string { return t.o.Text(configFormDescriptionOff) }

// Sections reads sections (proto field 4, repeated ConfigSection); each element
// is a ConfigSection object.
func (t ConfigForm) Sections() zap.List { return t.o.List(configFormSectionsOff) }

// DefaultValues reads default_values (proto field 5, map<string,ConfigValue>);
// each element is a ConfigValueMapEntry object.
func (t ConfigForm) DefaultValues() zap.List { return t.o.List(configFormDefaultValuesOff) }

// ConfigFormInput collects the field values for NewConfigForm.
type ConfigFormInput struct {
	FormID        string
	Title         string
	Description   string
	Sections      [][]byte
	DefaultValues [][]byte
}

// NewConfigForm builds a ZAP-encoded ConfigForm message.
func NewConfigForm(in ConfigFormInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(configFormSize)
	ob.SetText(configFormFormIDOff, in.FormID)
	ob.SetText(configFormTitleOff, in.Title)
	ob.SetText(configFormDescriptionOff, in.Description)
	sb := b.StartList(0)
	for _, elem := range in.Sections {
		sb.AddObjectBytes(elem)
	}
	ob.SetList(configFormSectionsOff, sb.FinishOffset(), len(in.Sections))
	db := b.StartList(0)
	for _, elem := range in.DefaultValues {
		db.AddObjectBytes(elem)
	}
	ob.SetList(configFormDefaultValuesOff, db.FinishOffset(), len(in.DefaultValues))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ConfigSection ---

const (
	configSectionSectionIDOff   = 0
	configSectionTitleOff       = 8
	configSectionDescriptionOff = 16
	configSectionFieldsOff      = 24
	configSectionSize           = 32
)

// ConfigSection is a zero-copy view into a ZAP-encoded ConfigSection message.
type ConfigSection struct{ o zap.Object }

// WrapConfigSection parses b and returns a typed view.
func WrapConfigSection(b []byte) (ConfigSection, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ConfigSection{}, err
	}
	return ConfigSection{o: m.Root()}, nil
}

func (t ConfigSection) SectionID() string   { return t.o.Text(configSectionSectionIDOff) }
func (t ConfigSection) Title() string       { return t.o.Text(configSectionTitleOff) }
func (t ConfigSection) Description() string { return t.o.Text(configSectionDescriptionOff) }

// Fields reads fields (proto field 4, repeated ConfigField); each element is a
// ConfigField object.
func (t ConfigSection) Fields() zap.List { return t.o.List(configSectionFieldsOff) }

// ConfigSectionInput collects the field values for NewConfigSection.
type ConfigSectionInput struct {
	SectionID   string
	Title       string
	Description string
	Fields      [][]byte
}

// NewConfigSection builds a ZAP-encoded ConfigSection message.
func NewConfigSection(in ConfigSectionInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(configSectionSize)
	ob.SetText(configSectionSectionIDOff, in.SectionID)
	ob.SetText(configSectionTitleOff, in.Title)
	ob.SetText(configSectionDescriptionOff, in.Description)
	lb := b.StartList(0)
	for _, elem := range in.Fields {
		lb.AddObjectBytes(elem)
	}
	ob.SetList(configSectionFieldsOff, lb.FinishOffset(), len(in.Fields))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ConfigField ---

const (
	configFieldNameOff              = 0
	configFieldLabelOff             = 8
	configFieldDescriptionOff       = 16
	configFieldHelpTextOff          = 24
	configFieldPlaceholderOff       = 32
	configFieldFieldTypeOff         = 40
	configFieldWidgetOff            = 44
	configFieldRequiredOff          = 48
	configFieldReadOnlyOff          = 49
	configFieldSensitiveOff         = 50
	configFieldMinValueOff          = 52
	configFieldMaxValueOff          = 60
	configFieldOptionsOff           = 68
	configFieldValidationRulesOff   = 76
	configFieldVisibleWhenFieldOff  = 84
	configFieldVisibleWhenEqualsOff = 92
	configFieldSize                 = 100
)

// ConfigField is a zero-copy view into a ZAP-encoded ConfigField message.
type ConfigField struct{ o zap.Object }

// WrapConfigField parses b and returns a typed view.
func WrapConfigField(b []byte) (ConfigField, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ConfigField{}, err
	}
	return ConfigField{o: m.Root()}, nil
}

func (t ConfigField) Name() string        { return t.o.Text(configFieldNameOff) }
func (t ConfigField) Label() string       { return t.o.Text(configFieldLabelOff) }
func (t ConfigField) Description() string { return t.o.Text(configFieldDescriptionOff) }
func (t ConfigField) HelpText() string    { return t.o.Text(configFieldHelpTextOff) }
func (t ConfigField) Placeholder() string { return t.o.Text(configFieldPlaceholderOff) }

// FieldType reads field_type (proto field 6, enum ConfigFieldType).
func (t ConfigField) FieldType() uint32 { return t.o.Uint32(configFieldFieldTypeOff) }

// Widget reads widget (proto field 7, enum ConfigWidget).
func (t ConfigField) Widget() uint32  { return t.o.Uint32(configFieldWidgetOff) }
func (t ConfigField) Required() bool  { return t.o.Bool(configFieldRequiredOff) }
func (t ConfigField) ReadOnly() bool  { return t.o.Bool(configFieldReadOnlyOff) }
func (t ConfigField) Sensitive() bool { return t.o.Bool(configFieldSensitiveOff) }
func (t ConfigField) MinValue() ConfigValue {
	return ConfigValue{o: childObject(t.o.Bytes(configFieldMinValueOff))}
}
func (t ConfigField) MaxValue() ConfigValue {
	return ConfigValue{o: childObject(t.o.Bytes(configFieldMaxValueOff))}
}

// Options reads options (proto field 13, repeated ConfigOption); each element
// is a ConfigOption object.
func (t ConfigField) Options() zap.List { return t.o.List(configFieldOptionsOff) }

// ValidationRules reads validation_rules (proto field 14, repeated
// ValidationRule); each element is a ValidationRule object.
func (t ConfigField) ValidationRules() zap.List { return t.o.List(configFieldValidationRulesOff) }

func (t ConfigField) VisibleWhenField() string { return t.o.Text(configFieldVisibleWhenFieldOff) }
func (t ConfigField) VisibleWhenEquals() ConfigValue {
	return ConfigValue{o: childObject(t.o.Bytes(configFieldVisibleWhenEqualsOff))}
}

// ConfigFieldInput collects the field values for NewConfigField. ConfigValue
// members and list elements take pre-built sub-buffers.
type ConfigFieldInput struct {
	Name              string
	Label             string
	Description       string
	HelpText          string
	Placeholder       string
	FieldType         uint32
	Widget            uint32
	Required          bool
	ReadOnly          bool
	Sensitive         bool
	MinValue          []byte
	MaxValue          []byte
	Options           [][]byte
	ValidationRules   [][]byte
	VisibleWhenField  string
	VisibleWhenEquals []byte
}

// NewConfigField builds a ZAP-encoded ConfigField message.
func NewConfigField(in ConfigFieldInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(configFieldSize)
	ob.SetText(configFieldNameOff, in.Name)
	ob.SetText(configFieldLabelOff, in.Label)
	ob.SetText(configFieldDescriptionOff, in.Description)
	ob.SetText(configFieldHelpTextOff, in.HelpText)
	ob.SetText(configFieldPlaceholderOff, in.Placeholder)
	ob.SetUint32(configFieldFieldTypeOff, in.FieldType)
	ob.SetUint32(configFieldWidgetOff, in.Widget)
	ob.SetBool(configFieldRequiredOff, in.Required)
	ob.SetBool(configFieldReadOnlyOff, in.ReadOnly)
	ob.SetBool(configFieldSensitiveOff, in.Sensitive)
	setChild(ob, configFieldMinValueOff, in.MinValue)
	setChild(ob, configFieldMaxValueOff, in.MaxValue)
	opb := b.StartList(0)
	for _, elem := range in.Options {
		opb.AddObjectBytes(elem)
	}
	ob.SetList(configFieldOptionsOff, opb.FinishOffset(), len(in.Options))
	vrb := b.StartList(0)
	for _, elem := range in.ValidationRules {
		vrb.AddObjectBytes(elem)
	}
	ob.SetList(configFieldValidationRulesOff, vrb.FinishOffset(), len(in.ValidationRules))
	ob.SetText(configFieldVisibleWhenFieldOff, in.VisibleWhenField)
	setChild(ob, configFieldVisibleWhenEqualsOff, in.VisibleWhenEquals)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ConfigOption ---

const (
	configOptionValueOff       = 0
	configOptionLabelOff       = 8
	configOptionDescriptionOff = 16
	configOptionDisabledOff    = 24
	configOptionSize           = 25
)

// ConfigOption is a zero-copy view into a ZAP-encoded ConfigOption message.
type ConfigOption struct{ o zap.Object }

// WrapConfigOption parses b and returns a typed view.
func WrapConfigOption(b []byte) (ConfigOption, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ConfigOption{}, err
	}
	return ConfigOption{o: m.Root()}, nil
}

func (t ConfigOption) Value() string       { return t.o.Text(configOptionValueOff) }
func (t ConfigOption) Label() string       { return t.o.Text(configOptionLabelOff) }
func (t ConfigOption) Description() string { return t.o.Text(configOptionDescriptionOff) }
func (t ConfigOption) Disabled() bool      { return t.o.Bool(configOptionDisabledOff) }

// ConfigOptionInput collects the field values for NewConfigOption.
type ConfigOptionInput struct {
	Value       string
	Label       string
	Description string
	Disabled    bool
}

// NewConfigOption builds a ZAP-encoded ConfigOption message.
func NewConfigOption(in ConfigOptionInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(configOptionSize)
	ob.SetText(configOptionValueOff, in.Value)
	ob.SetText(configOptionLabelOff, in.Label)
	ob.SetText(configOptionDescriptionOff, in.Description)
	ob.SetBool(configOptionDisabledOff, in.Disabled)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ValidationRule ---

const (
	validationRuleTypeOff         = 0
	validationRuleExpressionOff   = 4
	validationRuleErrorMessageOff = 12
	validationRuleSize            = 20
)

// ValidationRule is a zero-copy view into a ZAP-encoded ValidationRule message.
type ValidationRule struct{ o zap.Object }

// WrapValidationRule parses b and returns a typed view.
func WrapValidationRule(b []byte) (ValidationRule, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ValidationRule{}, err
	}
	return ValidationRule{o: m.Root()}, nil
}

// Type reads type (proto field 1, enum ValidationRuleType).
func (t ValidationRule) Type() uint32         { return t.o.Uint32(validationRuleTypeOff) }
func (t ValidationRule) Expression() string   { return t.o.Text(validationRuleExpressionOff) }
func (t ValidationRule) ErrorMessage() string { return t.o.Text(validationRuleErrorMessageOff) }

// ValidationRuleInput collects the field values for NewValidationRule.
type ValidationRuleInput struct {
	Type         uint32
	Expression   string
	ErrorMessage string
}

// NewValidationRule builds a ZAP-encoded ValidationRule message.
func NewValidationRule(in ValidationRuleInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(validationRuleSize)
	ob.SetUint32(validationRuleTypeOff, in.Type)
	ob.SetText(validationRuleExpressionOff, in.Expression)
	ob.SetText(validationRuleErrorMessageOff, in.ErrorMessage)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- AdminRuntimeDefaults ---

const (
	adminRuntimeDefaultsEnabledOff                    = 0
	adminRuntimeDefaultsDetectionIntervalMinutesOff   = 4
	adminRuntimeDefaultsDetectionTimeoutSecondsOff    = 8
	adminRuntimeDefaultsMaxJobsPerDetectionOff        = 12
	adminRuntimeDefaultsGlobalExecutionConcurrencyOff = 16
	adminRuntimeDefaultsPerWorkerExecutionConcOff     = 20
	adminRuntimeDefaultsRetryLimitOff                 = 24
	adminRuntimeDefaultsRetryBackoffSecondsOff        = 28
	adminRuntimeDefaultsJobTypeMaxRuntimeSecondsOff   = 32
	adminRuntimeDefaultsExecutionTimeoutSecondsOff    = 36
	adminRuntimeDefaultsSize                          = 40
)

// AdminRuntimeDefaults is a zero-copy view into a ZAP-encoded
// AdminRuntimeDefaults message.
type AdminRuntimeDefaults struct{ o zap.Object }

// WrapAdminRuntimeDefaults parses b and returns a typed view.
func WrapAdminRuntimeDefaults(b []byte) (AdminRuntimeDefaults, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return AdminRuntimeDefaults{}, err
	}
	return AdminRuntimeDefaults{o: m.Root()}, nil
}

func (t AdminRuntimeDefaults) Enabled() bool { return t.o.Bool(adminRuntimeDefaultsEnabledOff) }
func (t AdminRuntimeDefaults) DetectionIntervalMinutes() int32 {
	return t.o.Int32(adminRuntimeDefaultsDetectionIntervalMinutesOff)
}
func (t AdminRuntimeDefaults) DetectionTimeoutSeconds() int32 {
	return t.o.Int32(adminRuntimeDefaultsDetectionTimeoutSecondsOff)
}
func (t AdminRuntimeDefaults) MaxJobsPerDetection() int32 {
	return t.o.Int32(adminRuntimeDefaultsMaxJobsPerDetectionOff)
}
func (t AdminRuntimeDefaults) GlobalExecutionConcurrency() int32 {
	return t.o.Int32(adminRuntimeDefaultsGlobalExecutionConcurrencyOff)
}
func (t AdminRuntimeDefaults) PerWorkerExecutionConcurrency() int32 {
	return t.o.Int32(adminRuntimeDefaultsPerWorkerExecutionConcOff)
}
func (t AdminRuntimeDefaults) RetryLimit() int32 {
	return t.o.Int32(adminRuntimeDefaultsRetryLimitOff)
}
func (t AdminRuntimeDefaults) RetryBackoffSeconds() int32 {
	return t.o.Int32(adminRuntimeDefaultsRetryBackoffSecondsOff)
}
func (t AdminRuntimeDefaults) JobTypeMaxRuntimeSeconds() int32 {
	return t.o.Int32(adminRuntimeDefaultsJobTypeMaxRuntimeSecondsOff)
}
func (t AdminRuntimeDefaults) ExecutionTimeoutSeconds() int32 {
	return t.o.Int32(adminRuntimeDefaultsExecutionTimeoutSecondsOff)
}

// AdminRuntimeDefaultsInput collects the field values for NewAdminRuntimeDefaults.
type AdminRuntimeDefaultsInput struct {
	Enabled                       bool
	DetectionIntervalMinutes      int32
	DetectionTimeoutSeconds       int32
	MaxJobsPerDetection           int32
	GlobalExecutionConcurrency    int32
	PerWorkerExecutionConcurrency int32
	RetryLimit                    int32
	RetryBackoffSeconds           int32
	JobTypeMaxRuntimeSeconds      int32
	ExecutionTimeoutSeconds       int32
}

// NewAdminRuntimeDefaults builds a ZAP-encoded AdminRuntimeDefaults message.
func NewAdminRuntimeDefaults(in AdminRuntimeDefaultsInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(adminRuntimeDefaultsSize)
	ob.SetBool(adminRuntimeDefaultsEnabledOff, in.Enabled)
	ob.SetInt32(adminRuntimeDefaultsDetectionIntervalMinutesOff, in.DetectionIntervalMinutes)
	ob.SetInt32(adminRuntimeDefaultsDetectionTimeoutSecondsOff, in.DetectionTimeoutSeconds)
	ob.SetInt32(adminRuntimeDefaultsMaxJobsPerDetectionOff, in.MaxJobsPerDetection)
	ob.SetInt32(adminRuntimeDefaultsGlobalExecutionConcurrencyOff, in.GlobalExecutionConcurrency)
	ob.SetInt32(adminRuntimeDefaultsPerWorkerExecutionConcOff, in.PerWorkerExecutionConcurrency)
	ob.SetInt32(adminRuntimeDefaultsRetryLimitOff, in.RetryLimit)
	ob.SetInt32(adminRuntimeDefaultsRetryBackoffSecondsOff, in.RetryBackoffSeconds)
	ob.SetInt32(adminRuntimeDefaultsJobTypeMaxRuntimeSecondsOff, in.JobTypeMaxRuntimeSeconds)
	ob.SetInt32(adminRuntimeDefaultsExecutionTimeoutSecondsOff, in.ExecutionTimeoutSeconds)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- AdminRuntimeConfig ---

const (
	adminRuntimeConfigEnabledOff                    = 0
	adminRuntimeConfigDetectionIntervalMinutesOff   = 4
	adminRuntimeConfigDetectionTimeoutSecondsOff    = 8
	adminRuntimeConfigMaxJobsPerDetectionOff        = 12
	adminRuntimeConfigGlobalExecutionConcurrencyOff = 16
	adminRuntimeConfigPerWorkerExecutionConcOff     = 20
	adminRuntimeConfigRetryLimitOff                 = 24
	adminRuntimeConfigRetryBackoffSecondsOff        = 28
	adminRuntimeConfigJobTypeMaxRuntimeSecondsOff   = 32
	adminRuntimeConfigExecutionTimeoutSecondsOff    = 36
	adminRuntimeConfigSize                          = 40
)

// AdminRuntimeConfig is a zero-copy view into a ZAP-encoded AdminRuntimeConfig
// message.
type AdminRuntimeConfig struct{ o zap.Object }

// WrapAdminRuntimeConfig parses b and returns a typed view.
func WrapAdminRuntimeConfig(b []byte) (AdminRuntimeConfig, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return AdminRuntimeConfig{}, err
	}
	return AdminRuntimeConfig{o: m.Root()}, nil
}

func (t AdminRuntimeConfig) Enabled() bool { return t.o.Bool(adminRuntimeConfigEnabledOff) }
func (t AdminRuntimeConfig) DetectionIntervalMinutes() int32 {
	return t.o.Int32(adminRuntimeConfigDetectionIntervalMinutesOff)
}
func (t AdminRuntimeConfig) DetectionTimeoutSeconds() int32 {
	return t.o.Int32(adminRuntimeConfigDetectionTimeoutSecondsOff)
}
func (t AdminRuntimeConfig) MaxJobsPerDetection() int32 {
	return t.o.Int32(adminRuntimeConfigMaxJobsPerDetectionOff)
}
func (t AdminRuntimeConfig) GlobalExecutionConcurrency() int32 {
	return t.o.Int32(adminRuntimeConfigGlobalExecutionConcurrencyOff)
}
func (t AdminRuntimeConfig) PerWorkerExecutionConcurrency() int32 {
	return t.o.Int32(adminRuntimeConfigPerWorkerExecutionConcOff)
}
func (t AdminRuntimeConfig) RetryLimit() int32 {
	return t.o.Int32(adminRuntimeConfigRetryLimitOff)
}
func (t AdminRuntimeConfig) RetryBackoffSeconds() int32 {
	return t.o.Int32(adminRuntimeConfigRetryBackoffSecondsOff)
}
func (t AdminRuntimeConfig) JobTypeMaxRuntimeSeconds() int32 {
	return t.o.Int32(adminRuntimeConfigJobTypeMaxRuntimeSecondsOff)
}
func (t AdminRuntimeConfig) ExecutionTimeoutSeconds() int32 {
	return t.o.Int32(adminRuntimeConfigExecutionTimeoutSecondsOff)
}

// AdminRuntimeConfigInput collects the field values for NewAdminRuntimeConfig.
type AdminRuntimeConfigInput struct {
	Enabled                       bool
	DetectionIntervalMinutes      int32
	DetectionTimeoutSeconds       int32
	MaxJobsPerDetection           int32
	GlobalExecutionConcurrency    int32
	PerWorkerExecutionConcurrency int32
	RetryLimit                    int32
	RetryBackoffSeconds           int32
	JobTypeMaxRuntimeSeconds      int32
	ExecutionTimeoutSeconds       int32
}

// NewAdminRuntimeConfig builds a ZAP-encoded AdminRuntimeConfig message.
func NewAdminRuntimeConfig(in AdminRuntimeConfigInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(adminRuntimeConfigSize)
	ob.SetBool(adminRuntimeConfigEnabledOff, in.Enabled)
	ob.SetInt32(adminRuntimeConfigDetectionIntervalMinutesOff, in.DetectionIntervalMinutes)
	ob.SetInt32(adminRuntimeConfigDetectionTimeoutSecondsOff, in.DetectionTimeoutSeconds)
	ob.SetInt32(adminRuntimeConfigMaxJobsPerDetectionOff, in.MaxJobsPerDetection)
	ob.SetInt32(adminRuntimeConfigGlobalExecutionConcurrencyOff, in.GlobalExecutionConcurrency)
	ob.SetInt32(adminRuntimeConfigPerWorkerExecutionConcOff, in.PerWorkerExecutionConcurrency)
	ob.SetInt32(adminRuntimeConfigRetryLimitOff, in.RetryLimit)
	ob.SetInt32(adminRuntimeConfigRetryBackoffSecondsOff, in.RetryBackoffSeconds)
	ob.SetInt32(adminRuntimeConfigJobTypeMaxRuntimeSecondsOff, in.JobTypeMaxRuntimeSeconds)
	ob.SetInt32(adminRuntimeConfigExecutionTimeoutSecondsOff, in.ExecutionTimeoutSeconds)
	ob.FinishAsRoot()
	return b.Finish()
}
