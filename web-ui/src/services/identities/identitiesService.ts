import { RestService } from '../rest/restService';
import {
  IListResponse,
  IPagination,
  ISortParams
} from '../rest/restService.types';
import {
  IIdentityPrivateKeyResponse,
  IIdentityPublicKeyResponse,
  IIdentityResponse,
  INewIdentityRequest,
  IUpdateIdentityRequest
} from './identitiesService.types';

class IdentitiesService {
  public static async createIdentity(
    request: INewIdentityRequest
  ): Promise<IIdentityResponse> {
    try {
      return await RestService.post<IIdentityResponse>({
        url: `identities`,
        data: {
          ...request
        }
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }

  public static async getIdentities(
    pagination: IPagination,
    sortParams: ISortParams
  ): Promise<IListResponse<IIdentityResponse>> {
    try {
      return await RestService.get<IListResponse<IIdentityResponse>>({
        url: `identities?page=${pagination.page}&limit=${pagination.limit}&sortParam=${sortParams.sortParam}&sortDirection=${sortParams.sortDirection}`
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }

  public static async getIdentity(
    identityId: string
  ): Promise<IIdentityResponse> {
    try {
      return await RestService.get<IIdentityResponse>({
        url: `identities/${identityId}`
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }

  public static async getPrimaryIdentity(): Promise<IIdentityResponse> {
    try {
      return await RestService.get<IIdentityResponse>({
        url: `me`
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }

  public static async getIdentityPublicKey(
    identityId: string
  ): Promise<IIdentityPublicKeyResponse> {
    try {
      return await RestService.get<IIdentityPublicKeyResponse>({
        url: `identities/${identityId}/public-key`
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }

  public static async getIdentityPrivateKey(
    identityId: string
  ): Promise<IIdentityPrivateKeyResponse> {
    try {
      return await RestService.get<IIdentityPrivateKeyResponse>({
        url: `identities/${identityId}/private-key`
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }

  public static async updateIdentity(
    request: IUpdateIdentityRequest
  ): Promise<string> {
    try {
      return await RestService.put<string>({
        url: `identities/${request.identityId}`,
        data: {
          name: request.name,
          picture: request.picture,
          privateKey: request.privateKey
        }
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }

  public static async setPrimaryIdentity(identityId: string): Promise<string> {
    try {
      return await RestService.put<string>({
        url: `identities/${identityId}/set-primary`
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }

  public static async deleteIdentity(identityId: string): Promise<string> {
    try {
      return await RestService.delete<string>({
        url: `identities/${identityId}`
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }
}

export default IdentitiesService;
