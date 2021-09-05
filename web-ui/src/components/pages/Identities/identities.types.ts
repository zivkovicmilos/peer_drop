export interface IIdentitiesProps {}

export interface IIdentity {
  id: string;
  picture: string;
  name: string;

  publicKeyID: string;
  numWorkspaces: number;
  creationDate: string;
}