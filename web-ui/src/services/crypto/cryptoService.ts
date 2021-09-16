import { RestService } from '../rest/restService';
import {
  IGenerateKeyPairRequest,
  IGenerateKeyPairResponse,
  IValidatePrivateKeyResponse,
  IValidatePublicKeyResponse
} from './cryptoService.types';

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

  public static async validatePrivateKey(
    privateKey: string
  ): Promise<IValidatePrivateKeyResponse> {
    try {
      return await RestService.post<IValidatePrivateKeyResponse>({
        url: `crypto/validate-private-key`,
        data: {
          privateKey
        }
      });
    } catch (err) {
      console.warn(err);
      throw err;
    }
  }

  public static async generateKeyPair(
    request: IGenerateKeyPairRequest
  ): Promise<IGenerateKeyPairResponse> {
    try {
      return await RestService.post<IGenerateKeyPairResponse>({
        url: `crypto/generate-key-pair`,
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

export default CryptoService;
