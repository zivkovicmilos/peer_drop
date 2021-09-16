export interface INewIdentityRequest {
  name: string;
  picture: string;
  privateKey: string;
}

export interface IUpdateIdentityRequest extends INewIdentityRequest {
  identityId: string;
}

export interface IIdentityResponse {
  id: string;
  name: string;
  picture: string;
  dateCreated: string;
  publicKeyID: string;
  isPrimary: boolean;
  numWorkspaces: number;
}

export interface IIdentityPublicKeyResponse {
  publicKey: string;
}

export interface IIdentityPrivateKeyResponse {
  privateKey: string;
}
