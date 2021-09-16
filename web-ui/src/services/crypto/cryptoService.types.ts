export interface IValidatePublicKeyRequest {
  publicKey: string;
}

export interface IValidatePrivateKeyRequest {
  privateKey: string;
}

export interface IValidatePublicKeyResponse {
  publicKeyID: string;
}

export interface IValidatePrivateKeyResponse {
  publicKeyID: string;
  publicKey: string;
}

export interface IGenerateKeyPairResponse {
  privateKey: string;
  publicKeyID: string;
}

export interface IGenerateKeyPairRequest {
  keySize: number;
  name: string;
  email: string;
}
