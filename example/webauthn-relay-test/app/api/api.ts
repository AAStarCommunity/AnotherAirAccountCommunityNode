enum API {
    PASSKEY_REG = "/api/passkey/v1/reg",
    PASSKEY_REG_VERIFY = "/api/passkey/v1/reg/verify",
    PASSKEY_AUTH = "/api/passkey/v1/sign",
    PASSKEY_AUTH_VERIFY = "/api/passkey/v1/sign/verify",
    PASSKEY_PAYMENT = "/api/passkey/v1/payment/sign",
    PASSKEY_PAYMENT_VERIFY = "/api/passkey/v1/payment/sign/verify",
}

export default API;
