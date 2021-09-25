import { RestService } from '../rest/restService';
import { ISearchResults } from './searchService.types';

class SearchService {
  public static async getSearchResults(input: string): Promise<ISearchResults> {
    try {
      return await RestService.get<ISearchResults>({
        url: `search?input=${input}`
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }
}

export default SearchService;