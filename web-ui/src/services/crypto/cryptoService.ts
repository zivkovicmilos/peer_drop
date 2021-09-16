import { RestService } from '../rest/restService';
import { IValidatePublicKeyResponse } from './cryptoService.types';

class CryptoService {
  public static async validatePublicKey(
    publicKey: string
  ): Promise<IValidatePublicKeyResponse> {
    try {
      return await RestService.post<IValidatePublicKeyResponse>({
        url: `crypto/validate-public-key`,
        data: {
          publicKey
        }
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }
}

export default CryptoService