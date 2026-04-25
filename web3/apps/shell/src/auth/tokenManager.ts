// shell/src/auth/tokenManager.ts
import axios from "axios";

let accessToken: string | null = null;
let refreshPromise: Promise<string> | null = null;

export const tokenManager = {
  get: () => accessToken,

  set: (token: string) => {
    accessToken = token;
  },

  clear: () => {
    accessToken = null;
  },

  refresh: (): Promise<string> => {
    // якщо вже рефрешимо — повертаємо той самий Promise
    // захист від race condition при паралельних 401
    if (refreshPromise) return refreshPromise;

    refreshPromise = axios
      .post<{ access_token: string }>('/api/auth/refresh', null, {
        withCredentials: true, // refresh_token в httpOnly cookie
        skipAuthInterceptor: true, // <-- щоб сам refresh не потрапив у перехоплювач
      })
      .then(({ data }) => {
        accessToken = data.access_token;
        console.log(accessToken);
        return data.access_token;
      })
      .catch((err) => {
        accessToken = null;
        window.dispatchEvent(new CustomEvent('auth:logout'));
        return Promise.reject(err);
      })
      .finally(() => {
        refreshPromise = null;
      });

    return refreshPromise;
  },
};