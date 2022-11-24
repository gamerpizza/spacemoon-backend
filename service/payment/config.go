package payment

type (
	PaymentCfg struct {
		PayPal PaypalCfg `yaml:"paypal"`
		Stripe StripeCfg `yaml:"stripe"`
	}

	PaypalCfg struct {
		ClientID     string `yaml:"clientId"`
		ClientSecret string `yaml:"clientSecret"`
		URL          string `yaml:"url"`
	}

	StripeCfg struct {
		PublicKey  string `yaml:"publicKey"`
		APIKey     string `yaml:"apiKey"`
		CancelUrl  string `yaml:"cancelUrl"`
		SuccessUrl string `yaml:"successUrl"`
	}
)
