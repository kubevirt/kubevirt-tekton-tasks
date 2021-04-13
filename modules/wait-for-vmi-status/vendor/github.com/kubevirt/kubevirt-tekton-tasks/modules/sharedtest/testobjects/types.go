package testobjects

import "sigs.k8s.io/yaml"

type CloudConfigSSHKeys struct {
	RSAPrivate string `json:"rsa_private"`
	RSAPublic  string `json:"rsa_public"`
}

type CloudConfig struct {
	Password string `json:"password"`
	Chpasswd struct {
		Expire bool `json:"expire"`
	} `json:"chpasswd"`
	SSHAuthorizedKeys []string           `json:"ssh_authorized_keys"`
	SSHKeys           CloudConfigSSHKeys `json:"ssh_keys"`
}

func (c CloudConfig) ToString() string {
	cloudConfigBytes, _ := yaml.Marshal(c)
	cloudConfigStr := string(cloudConfigBytes)

	return "#cloud-config\n" + cloudConfigStr
}
