package vault

import (
	"fmt"
	"testing"

	"github.com/BESTSELLER/harpocrates/config"
	"github.com/BESTSELLER/harpocrates/files"
	"github.com/BESTSELLER/harpocrates/util"
	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"gotest.tools/assert"
)

var testVault vaultTest

// rootVaultToken is the Vault token used for tests
var rootVaultToken = "unittesttoken"

type vaultTest struct {
	Cluster *vault.TestCluster
	Client  *api.Client
}

func TestMain(t *testing.T) {
	fmt.Println("hej")
	testVault = GetTestVaultServer(t)
}

// GetTestVaultServer creates the test server
func GetTestVaultServer(t *testing.T) vaultTest {
	t.Helper()

	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		DevToken: rootVaultToken,
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client

	// put secrets
	secretPath := fmt.Sprintf("secret/data/secret")
	secret := map[string]interface{}{"key1": "value1", "key2": "value2", "key3": "value3"}

	_, err := client.Logical().Write(secretPath, secret)
	if err != nil {
		panic(err)
	}

	// return server
	return vaultTest{
		Cluster: cluster,
		Client:  client,
	}

}

// TestExtractSecretsAsExpected tests if a simple secret is extracted correct
func TestExtractSecretsAsExpected(t *testing.T) {
	// arrange

	// define input
	data := files.Read("../test_data/single_secret.yaml")
	input := util.ReadInput(data)

	// mock prefix
	config.Config.Prefix = input.Prefix

	var vaultClient *API
	vaultClient = &API{
		Client: testVault.Client,
	}

	// act
	result := vaultClient.ExtractSecrets(input)

	// assert
	expected := fmt.Sprintf("%v", map[string]interface{}{input.Prefix + "key1": "value1", input.Prefix + "key2": "value2", input.Prefix + "key3": "value3"})
	actual := fmt.Sprintf("%v", result)

	assert.Equal(t, expected, actual)

}

// TestExtractSecretsWithPrefixAsExpected tests if a simple secret is extracted correct
func TestExtractSecretsWithPrefixAsExpected(t *testing.T) {
	// arrange

	// define input
	data := files.Read("../test_data/keys_with_prefix.yaml")
	input := util.ReadInput(data)

	// mock prefix
	config.Config.Prefix = input.Prefix

	var vaultClient *API
	vaultClient = &API{
		Client: testVault.Client,
	}

	// act
	result := vaultClient.ExtractSecrets(input)

	// assert
	expected := fmt.Sprintf("%v", map[string]interface{}{"PRE_key1": "value1", "FIX_key2": "value2", input.Prefix + "key3": "value3"})
	actual := fmt.Sprintf("%v", result)

	assert.Equal(t, expected, actual)

}
