import axios, {
  AxiosRequestConfig,
  AxiosInstance,
  AxiosResponse,
  AxiosPromise,
} from "axios";

export interface ResponseData {
  code: 200 | 400 | 500;
  message: string;
  data: any;
}
interface ApiResponse<T = any> extends AxiosResponse {
  token: string | undefined;
  code: number;
  message: string;
  data: T;
}

class HtttpRequest {
  // Merge configuration items
  private mergeConfig(...configs: AxiosRequestConfig[]): AxiosRequestConfig {
    return Object.assign({}, ...configs);
  }

  // Set the get request alias
  public get(
    url: string,
    params: any = {},
    config: AxiosRequestConfig = {}
  ): AxiosPromise<ApiResponse> {
    const newConfig = this.mergeConfig(config, {
      url,
      method: "GET",
      params,
    });
    return this.request(newConfig);
  }

  // Set the post request alias
  public post(
    url: string,
    data: any = {},
    config: AxiosRequestConfig = {}
  ): AxiosPromise<ApiResponse> {
    const newConfig = this.mergeConfig(config, {
      method: "POST",
      url,      
      data,
    });
    return this.request(newConfig);
  }

  public put(
    url: string,
    data: any = {},
    config: AxiosRequestConfig = {}
  ): AxiosPromise<ApiResponse> {
    const newConfig = this.mergeConfig(config, { url, data, method: "PUT" });
    return this.request(newConfig);
  }

  public delete(
    url: string,
    params: any = {},
    config: AxiosRequestConfig = {}
  ): AxiosPromise<ApiResponse> {
    const newConfig = this.mergeConfig(config, {
      url,
      params,
      method: "DELETE",
    });
    return this.request(newConfig);
  }

  public request(config: AxiosRequestConfig): AxiosPromise<ApiResponse> {
    const instance: AxiosInstance = axios.create();
    this.interceptor(instance);
    return instance(config);
  }

  // Add intercept
  private interceptor(instance: AxiosInstance) {
    // Intercept request
    instance.interceptors.request.use(
      (config: any) => {
        config.baseURL = process.env.NEXT_PUBLIC_BASE_URL?.startsWith("http")
          ? process.env.NEXT_PUBLIC_BASE_URL
          : location.origin + "/" + process.env.NEXT_PUBLIC_BASE_URL;
        return config;
      },
      (error: any) => {
        return Promise.reject(error);
      }
    );

    //Intercept response
    instance.interceptors.response.use(
      (response: AxiosResponse) => {
        const {
          data: { code, message, data },
          status,
        } = response;
        if (status === 200 && code === 200) {
        } else {
          console.log(message)
        }
        return response;
      },
      (error: any) => {
        console.log(error.message);
        return Promise.reject(error);
      }
    );
  }
}

export default HtttpRequest;
