import { AxiosRequestConfig } from 'axios';

export interface RequestParams {
  url: string;
  data?: any;
  config?: AxiosRequestConfig;
}

export interface IListResponse<T> {
  count: number;
  data: T[];
}

export interface IPagination {
  page?: number;
  limit?: number;
}
