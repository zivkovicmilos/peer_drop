export interface IIdentitySortProps {
  activeSort: EIdentitySortParam;
  sortDirection: EIdentitySortDirection;

  setActiveSort: (activeSort: EIdentitySortParam) => void;
  setSortDirection: (activeSort: EIdentitySortDirection) => void;
}

export enum EIdentitySortParam {
  NAME = 'Name',
  PUBLIC_KEY = 'Public Key ID',
  NUMBER_OF_WORKSPACES = 'Num. of Workspaces',
  CREATION_DATE = 'Creation Date'
}

export enum EIdentitySortDirection {
  ASC,
  DESC
}
