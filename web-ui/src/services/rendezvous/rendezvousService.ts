import { RestService } from '../rest/restService';
import { IListResponse } from '../rest/restService.types';
import {
  IRendezvousNodeResponse,
  IRendezvousWrapper
} from './rendezvousService.types';

class RendezvousService {
  public static async getRendezvousNodes(): Promise<
    IListResponse<IRendezvousNodeResponse>
  > {
    try {
      return await RestService.get<IListResponse<IRendezvousNodeResponse>>({
        url: 'rendezvous'
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }

  public static async addNewRendezvousNode(
    request: IRendezvousWrapper
  ): Promise<string> {
    try {
      return await RestService.post<string>({
        url: 'rendezvous',
        data: {
          ...request
        }
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }

  public static async removeRendezvousNode(
    request: IRendezvousWrapper
  ): Promise<string> {
    try {
      return await RestService.delete<string>({
        url: 'rendezvous',
        data: {
          ...request
        }
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }
}

export default RendezvousService;