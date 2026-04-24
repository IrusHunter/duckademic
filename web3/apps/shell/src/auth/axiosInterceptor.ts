// shell/src/auth/axiosInterceptor.ts

import axios from 'axios';
import { tokenManager } from './tokenManager';

export function setupAxiosInterceptor() {

  // ── Request: додаємо токен ──────────────────────────────────────
  axios.interceptors.request.use((config) => {
    if (config.skipAuthInterceptor) return config;

    const token = tokenManager.get();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  });

  // ── Response: 401 → refresh → retry ────────────────────────────
  axios.interceptors.response.use(
    (response) => response,
    async (error) => {
      const original = error.config;

      const is401 = error.response?.status === 401;
      const alreadyRetried = original._retry;
      const isRefreshCall = original.skipAuthInterceptor;

      if (is401 && !alreadyRetried && !isRefreshCall) {
        original._retry = true;

        try {
          const newToken = await tokenManager.refresh();
          original.headers.Authorization = `Bearer ${newToken}`;
          return axios(original); // повторний запит — MFE нічого не знає
        } catch {
          // refresh провалився — logout вже задиспатчений у tokenManager
          return Promise.reject(error);
        }
      }

      return Promise.reject(error);
    }
  );
}