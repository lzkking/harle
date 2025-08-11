package crypto

var (
	RSAPublicKey  = "" //	全局公钥
	RSAPrivateKey = "" //	全局的私钥
)

func init() {
	//	读取配置文件获取RSA的公钥和私钥,没有的话进行创建并写到配置文件中
}
