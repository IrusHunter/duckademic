// shell/src/auth/axios.d.ts

import 'axios';

declare module 'axios' {
  interface AxiosRequestConfig {
    skipAuthInterceptor?: boolean;
  }
}