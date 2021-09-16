import { RestService } from '../rest/restService';
import { IListResponse, IPagination } from '../rest/restService.types';
import {
  IContactResponse,
  INewContactRequest,
  INewContactResponse
} from './contactsService.types';

class ContactsService {
  public static async createContact(
    request: INewContactRequest
  ): Promise<INewContactResponse> {
    try {
      return await RestService.post<INewContactResponse>({
        url: `contacts`,
        data: {
          name: request.name,
          publicKey: request.publicKey
        }
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }

  public static async getContacts(
    pagination: IPagination
  ): Promise<IListResponse<IContactResponse>> {
    try {
      return await RestService.get<IListResponse<IContactResponse>>({
        url: `contacts?page=${pagination.page}&limit=${pagination.limit}`
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }
}

export default ContactsService;
