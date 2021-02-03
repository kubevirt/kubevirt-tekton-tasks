package parse

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/env"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"k8s.io/apimachinery/pkg/util/validation"
	"strings"
	"unicode"
)

func (c *CLIOptions) trimSpaces() {
	for _, strVariablePtr := range []*string{&c.PublicKeySecretName, &c.PublicKeySecretNamespace, &c.PrivateKeySecretName, &c.PrivateKeySecretNamespace} {
		*strVariablePtr = strings.TrimSpace(*strVariablePtr)
	}

	for i, v := range c.PrivateKeyConnectionOptions {
		c.PrivateKeyConnectionOptions[i] = strings.TrimLeftFunc(v, unicode.IsSpace)
	}
}

func (c *CLIOptions) validateNames() error {
	for optionName, optionValue := range map[string]string{
		publicKeySecretNameOptionName:       c.PublicKeySecretName,
		publicKeySecretNamespaceOptionName:  c.PublicKeySecretNamespace,
		privateKeySecretNameOptionName:      c.PrivateKeySecretName,
		privateKeySecretNamespaceOptionName: c.PrivateKeySecretNamespace,
	} {
		if optionValue != "" {
			if errors := validation.IsDNS1123Subdomain(optionValue); len(errors) > 0 {
				return zerrors.NewMissingRequiredError("invalid %v value: %v", optionName, strings.Join(errors, ", "))
			}
		}
	}
	return nil
}

func (c *CLIOptions) resolveDefaultNamespaces() error {
	var activeNamespace string

	for optionName, namespacePtr := range map[string]*string{
		publicKeySecretNamespaceOptionName:  &c.PublicKeySecretNamespace,
		privateKeySecretNamespaceOptionName: &c.PrivateKeySecretNamespace,
	} {
		if *namespacePtr == "" {
			if activeNamespace == "" {
				var err error
				if activeNamespace, err = env.GetActiveNamespace(); err != nil {
					return zerrors.NewMissingRequiredError("%v: %v option is empty", err.Error(), optionName)
				}
			}
			*namespacePtr = activeNamespace

		}
	}
	return nil
}
