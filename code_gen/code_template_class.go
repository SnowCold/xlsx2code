package code_gen

type CodeTemplateMember struct {
	MemberType        string
	Name              string
	FormatPrivateName string
	NeedVirtual       bool
	FormatGetPropName string
	Comment           string
}

type CodeTemplateClass struct {
	ClassName      string
	Members        []CodeTemplateMember
	KeyType        string
	KeyPropName    string
	NeedCopyMethod bool
	FileBaseName   string
}
