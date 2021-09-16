export interface INewContactRequest {
  name: string;
  publicKey: string;
}

export interface IUpdateContactRequest extends INewContactRequest {
  contactId: string;
}

export interface INewContactResponse {
  id: string;
  name: string;
  publicKey: string;
  publicKeyID: string;
}

export interface IContactResponse {
  id: string;
  name: string;
  email: string;
  dateAdded: string;
  publicKey: string;
  publicKeyID: string;
}
