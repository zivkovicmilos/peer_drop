import { RestService } from '../rest/restService';
import { IListResponse, IPagination } from '../rest/restService.types';
import {
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
    pagination: IPagination
    // TODO add sort param
  ): Promise<IListResponse<IIdentityResponse>> {
    try {
      return await RestService.get<IListResponse<IIdentityResponse>>({
        url: `identities?page=${pagination.page}&limit=${pagination.limit}`
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
}

export default IdentitiesService;
