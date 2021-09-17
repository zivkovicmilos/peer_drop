import { RestService } from '../rest/restService';

class CommonService {
  public static async shutdown(): Promise<void> {
    try {
      return await RestService.post<void>({
        url: `shutdown`
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }
}

export default CommonService;
