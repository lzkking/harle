package assets

import _ "embed"

//go:embed pem/rsa_public.pem
var RsaPublicPem []byte
