package secret

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/types"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/utils/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/generate-ssh-keys/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	machinerytypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/rand"
	"net/http"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type SecretFacade struct {
	clioptions *parse.CLIOptions
	kubeClient *kubernetes.Clientset
	keys       types.SshKeys
}

func NewSecretFacade(clioptions *parse.CLIOptions, keys types.SshKeys) (*SecretFacade, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	kubeClient := kubernetes.NewForConfigOrDie(config)

	return &SecretFacade{clioptions: clioptions, kubeClient: kubeClient, keys: keys}, nil
}

func (s SecretFacade) CheckPrivateKeySecretExistence() error {
	if secretName := s.clioptions.GetPrivateKeySecretName(); secretName != "" {
		log.GetLogger().Debug("checking private key existence arguments", zap.String("secretName", secretName))
		if _, err := s.kubeClient.CoreV1().Secrets(s.clioptions.GetPrivateKeySecretNamespace()).Get(secretName, v1.GetOptions{}); !zerrors.IsStatusError(err, http.StatusNotFound) {
			return zerrors.NewMissingRequiredError("%v secret already exists", secretName)
		}
	}
	return nil
}

func (s *SecretFacade) GetPublicKeySecret() (*corev1.Secret, error) {
	var secret *corev1.Secret
	if secretName := s.clioptions.GetPublicKeySecretName(); secretName != "" {
		var err error
		secret, err = s.kubeClient.CoreV1().Secrets(s.clioptions.GetPublicKeySecretNamespace()).Get(secretName, v1.GetOptions{})
		if err != nil {
			secret = nil
			if !zerrors.IsStatusError(err, http.StatusNotFound) {
				return nil, err
			}
		}
	}

	return secret, nil
}

func (s *SecretFacade) CreatePrivateKeySecret() (*corev1.Secret, error) {
	data := map[string]string{}

	for key, value := range s.clioptions.GetPrivateKeyConnectionOptions() {
		data[key] = value
	}

	data[constants.ConnectionOptions.Type] = constants.ConnectionSSHType
	data[constants.ConnectionOptions.PrivateKey] = s.keys.PrivateKey

	secret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{},
		StringData: data,
	}

	if secretName := s.clioptions.GetPrivateKeySecretName(); secretName != "" {
		secret.Name = secretName
	} else {
		secret.GenerateName = constants.PrivateKeyGenerateName
	}

	log.GetLogger().Debug("creating private key secret")
	return s.kubeClient.CoreV1().Secrets(s.clioptions.GetPrivateKeySecretNamespace()).Create(secret)

}

func (s *SecretFacade) AppendPublicKeySecret(secret *corev1.Secret) (*corev1.Secret, error) {
	var publicKeyId string

	for {
		publicKeyId = generatePublicKeyId()
		if secret.Data[publicKeyId] == nil {
			break
		}
	}

	publicKeyBase64 := base64.StdEncoding.EncodeToString([]byte(s.keys.PublicKey))

	patches := []SecretPatch{
		{
			Op:    "add",
			Path:  "/data/" + publicKeyId,
			Value: publicKeyBase64,
		},
	}

	patchBytes, err := json.Marshal(patches)

	if err != nil {
		return nil, err
	}

	log.GetLogger().Debug("appending public key secret")
	return s.kubeClient.CoreV1().Secrets(s.clioptions.GetPublicKeySecretNamespace()).Patch(secret.Name, machinerytypes.JSONPatchType, patchBytes)
}

func (s *SecretFacade) CreatePublicKeySecret() (*corev1.Secret, error) {
	secret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{},

		StringData: map[string]string{
			generatePublicKeyId(): s.keys.PublicKey,
		},
	}

	if secretName := s.clioptions.GetPublicKeySecretName(); secretName != "" {
		secret.Name = secretName
	} else {
		secret.GenerateName = constants.PublicKeyGenerateName
	}

	log.GetLogger().Debug("creating public key secret")
	return s.kubeClient.CoreV1().Secrets(s.clioptions.GetPublicKeySecretNamespace()).Create(secret)
}

func (s *SecretFacade) DeleteSecret(secret *corev1.Secret) error {
	if secret == nil {
		return nil
	}
	log.GetLogger().Debug("deleting secret", zap.String("namespace", secret.Namespace), zap.String("name", secret.Name))
	return s.kubeClient.CoreV1().Secrets(secret.Namespace).Delete(secret.Name, &v1.DeleteOptions{})
}

func generatePublicKeyId() string {
	return fmt.Sprintf("id-rsa-%v.pub", rand.String(5))
}
