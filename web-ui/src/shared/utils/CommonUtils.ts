import {
  EIdentitySortDirection,
  EIdentitySortParam
} from '../../components/molecules/IdentitySort/identitySort.types';
import { ISortParams } from '../../services/rest/restService.types';

export default class CommonUtils {
  static formatBytes(bytes: number, decimals = 2) {
    if (bytes === 0) return '0 Bytes';

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
  }

  static removeLineBreaks(originalString: string): string {
    return originalString.replace(/(\r)/gm, '');
  }

  static fileToBase64 = (file: File): Promise<string> => {
    return new Promise<string>((resolve, reject) => {
      const reader = new FileReader();

      reader.readAsDataURL(file);
      reader.onload = () => resolve(reader.result.toString());

      reader.onerror = (error) => reject(error);
    });
  };

  static formatSortParams = (
    sortParam: EIdentitySortParam,
    sortDirection: EIdentitySortDirection
  ): ISortParams => {
    let sortParams: ISortParams = {};

    switch (sortParam) {
      case EIdentitySortParam.NAME:
        sortParams.sortParam = 'name';
        break;
      case EIdentitySortParam.PUBLIC_KEY:
        sortParams.sortParam = 'publicKeyID';
        break;
      case EIdentitySortParam.CREATION_DATE:
        sortParams.sortParam = 'dateCreated';
        break;
      case EIdentitySortParam.NUMBER_OF_WORKSPACES:
        sortParams.sortParam = 'numWorkspaces';
        break;
    }

    if (sortDirection == EIdentitySortDirection.ASC) {
      sortParams.sortDirection = 'asc';
    } else {
      sortParams.sortDirection = 'desc';
    }

    return sortParams;
  };

  static formatMnemonic(mnemonic: string): string {
    return mnemonic.replace(/\s/g, '-');
  }

  static unformatMnemonic(mnemonic: string): string {
    return mnemonic.replace(/-/g, ' ');
  }
}
